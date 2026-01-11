package packing

import (
	"fmt"
	"sort"
)

type Dimensions struct {
	Height int
	Width  int
	Length int
}

func (d Dimensions) Volume() int {
	return d.Height * d.Width * d.Length
}

// Rotations returns all unique 3D rotations of the dimensions.
func (d Dimensions) Rotations() []Dimensions {
	perms := []Dimensions{
		{d.Height, d.Width, d.Length},
		{d.Height, d.Length, d.Width},
		{d.Width, d.Height, d.Length},
		{d.Width, d.Length, d.Height},
		{d.Length, d.Height, d.Width},
		{d.Length, d.Width, d.Height},
	}

	unique := make([]Dimensions, 0)
	seen := make(map[Dimensions]bool)

	for _, p := range perms {
		if !seen[p] {
			seen[p] = true
			unique = append(unique, p)
		}
	}

	return unique
}

func (d Dimensions) FitsIn(space Dimensions) bool {
	return d.Height <= space.Height && d.Width <= space.Width && d.Length <= space.Length
}

type Item struct {
	ProductID string
	Dim       Dimensions
	Volume    int
}

type PackedBox struct {
	BoxType   BoxType
	Products  []string
	freeSpaces []Dimensions
}

func newPackedBox(bt BoxType) PackedBox {
	return PackedBox{
		BoxType:    bt,
		Products:   []string{},
		freeSpaces: []Dimensions{{Height: bt.Height, Width: bt.Width, Length: bt.Length}},
	}
}

func (b *PackedBox) boxVolume() int {
	return b.BoxType.Height * b.BoxType.Width * b.BoxType.Length
}

type placement struct {
	boxIndex     int
	spaceIndex   int
	rot          Dimensions
	wasteVolume  int
	boxVolume    int
}

// TryPlace tries to place an item into this box using free-space splitting.
// Returns true if placed.
func (b *PackedBox) TryPlace(item Item) bool {
	best := placement{wasteVolume: int(^uint(0) >> 1)} // max int
	found := false

	for si, space := range b.freeSpaces {
		for _, rot := range item.Dim.Rotations() {
			if !rot.FitsIn(space) {
				continue
			}
			waste := space.Volume() - rot.Volume()
			if waste < best.wasteVolume {
				best = placement{
					spaceIndex:  si,
					rot:         rot,
					wasteVolume: waste,
					boxVolume:   b.boxVolume(),
				}
				found = true
			}
		}
	}

	if !found {
		return false
	}

	// Place into chosen space and split it deterministically.
	space := b.freeSpaces[best.spaceIndex]
	rot := best.rot

	// Remove used space
	b.freeSpaces = append(b.freeSpaces[:best.spaceIndex], b.freeSpaces[best.spaceIndex+1:]...)

	// Split strategy (guillotine-like, simple and deterministic):
	// Place item at origin of the space; generate up to 3 residual spaces.
	// 1) Right slice (remaining width)
	if space.Width-rot.Width > 0 {
		b.freeSpaces = append(b.freeSpaces, Dimensions{
			Height: space.Height,
			Width:  space.Width - rot.Width,
			Length: space.Length,
		})
	}
	// 2) Front slice (remaining length) within the item's width region
	if space.Length-rot.Length > 0 {
		b.freeSpaces = append(b.freeSpaces, Dimensions{
			Height: space.Height,
			Width:  rot.Width,
			Length: space.Length - rot.Length,
		})
	}
	// 3) Top slice (remaining height) within the item's footprint
	if space.Height-rot.Height > 0 {
		b.freeSpaces = append(b.freeSpaces, Dimensions{
			Height: space.Height - rot.Height,
			Width:  rot.Width,
			Length: rot.Length,
		})
	}

	// Optional: keep spaces sorted by volume descending to try larger spaces first.
	sort.Slice(b.freeSpaces, func(i, j int) bool {
		return b.freeSpaces[i].Volume() > b.freeSpaces[j].Volume()
	})

	b.Products = append(b.Products, item.ProductID)
	return true
}

type OrderPackingResult struct {
	Boxes []PackedBox
}

// PackOrder packs items into available boxes minimizing number of boxes (heuristic).
// Returns an error if an item can't fit in any available box (even with rotation).
func PackOrder(items []Item, boxTypes []BoxType) (OrderPackingResult, error) {
	if len(items) == 0 {
		return OrderPackingResult{Boxes: []PackedBox{}}, nil
	}

	// Sort items by volume desc (FFD)
	for i := range items {
		items[i].Volume = items[i].Dim.Volume()
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Volume == items[j].Volume {
			return items[i].ProductID < items[j].ProductID
		}
		return items[i].Volume > items[j].Volume
	})

	// Sort box types by volume asc (pick smallest that fits when opening new one)
	sort.Slice(boxTypes, func(i, j int) bool {
		vi := boxTypes[i].Height * boxTypes[i].Width * boxTypes[i].Length
		vj := boxTypes[j].Height * boxTypes[j].Width * boxTypes[j].Length
		if vi == vj {
			return boxTypes[i].ID < boxTypes[j].ID
		}
		return vi < vj
	})

	var opened []PackedBox

	for _, it := range items {
		placed := false

		// Try existing boxes first
		for bi := range opened {
			if opened[bi].TryPlace(it) {
				placed = true
				break
			}
		}

		if placed {
			continue
		}

		// Open a new box: smallest box that can fit this item in some rotation
		var chosen *BoxType
		for i := range boxTypes {
			bt := boxTypes[i]
			space := Dimensions{Height: bt.Height, Width: bt.Width, Length: bt.Length}
			canFit := false
			for _, rot := range it.Dim.Rotations() {
				if rot.FitsIn(space) {
					canFit = true
					break
				}
			}
			if canFit {
				chosen = &bt
				break
			}
		}
		if chosen == nil {
			return OrderPackingResult{}, fmt.Errorf("produto '%s' não cabe em nenhuma caixa disponível (mesmo com rotação)", it.ProductID)
		}

		nb := newPackedBox(*chosen)
		if !nb.TryPlace(it) {
			// This should not happen given the fit check, but keep safe.
			return OrderPackingResult{}, fmt.Errorf("falha inesperada ao alocar produto '%s' na caixa '%s'", it.ProductID, chosen.ID)
		}
		opened = append(opened, nb)
	}

	return OrderPackingResult{Boxes: opened}, nil
}
