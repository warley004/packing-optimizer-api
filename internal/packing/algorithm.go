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

// Rotations retorna todas as rotações 3D únicas das dimensões.
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

// rotationsFor devolve todas as rotações quando allowRotation=true ou apenas a orientação original.
// allowRotation define se mantemos orientação rígida (false) ou testamos todas as permutações 3D (true).
func rotationsFor(d Dimensions, allowRotation bool) []Dimensions {
	if allowRotation {
		return d.Rotations()
	}
	return []Dimensions{d}
}

type Item struct {
	ProductID string
	Dim       Dimensions
	Volume    int
	Index     int // posição do produto no pedido (ordem original do input)
}

type packedProduct struct {
	ID    string
	Index int
}

type PackedBox struct {
	BoxType    BoxType
	Products   []packedProduct
	freeSpaces []Dimensions
}

func newPackedBox(bt BoxType) PackedBox {
	return PackedBox{
		BoxType:    bt,
		Products:   []packedProduct{},
		freeSpaces: []Dimensions{{Height: bt.Height, Width: bt.Width, Length: bt.Length}},
	}
}

func (b *PackedBox) boxVolume() int {
	return b.BoxType.Height * b.BoxType.Width * b.BoxType.Length
}

type placement struct {
	spaceIndex  int
	rot         Dimensions
	wasteVolume int
}

// TryPlace tenta colocar o item aplicando free-space splitting.
// Retorna true se o item couber.
func (b *PackedBox) TryPlace(item Item, allowRotation bool) bool {
	best := placement{wasteVolume: int(^uint(0) >> 1)} // max int
	found := false

	for si, space := range b.freeSpaces {
		for _, rot := range rotationsFor(item.Dim, allowRotation) {
			if !rot.FitsIn(space) {
				continue
			}

			waste := space.Volume() - rot.Volume()
			if waste < best.wasteVolume {
				best = placement{
					spaceIndex:  si,
					rot:         rot,
					wasteVolume: waste,
				}
				found = true
			}
		}
	}

	if !found {
		return false
	}

	// Aloca no espaço escolhido e faz o split determinístico.
	space := b.freeSpaces[best.spaceIndex]
	rot := best.rot

	// Remove used space
	b.freeSpaces = append(b.freeSpaces[:best.spaceIndex], b.freeSpaces[best.spaceIndex+1:]...)

	// Estratégia de split (tipo guilhotina, simples e determinística) mantém layout previsível para encaixes futuros.
	// Coloca o item na origem do espaço; gera até 3 sobras.
	// 1) Slice lateral (largura restante)
	if space.Width-rot.Width > 0 {
		b.freeSpaces = append(b.freeSpaces, Dimensions{
			Height: space.Height,
			Width:  space.Width - rot.Width,
			Length: space.Length,
		})
	}
	// 2) Slice frontal (comprimento restante) dentro da largura do item
	if space.Length-rot.Length > 0 {
		b.freeSpaces = append(b.freeSpaces, Dimensions{
			Height: space.Height,
			Width:  rot.Width,
			Length: space.Length - rot.Length,
		})
	}
	// 3) Slice superior (altura restante) dentro da base do item
	if space.Height-rot.Height > 0 {
		b.freeSpaces = append(b.freeSpaces, Dimensions{
			Height: space.Height - rot.Height,
			Width:  rot.Width,
			Length: rot.Length,
		})
	}

	// Mantém os espaços ordenados por volume decrescente para tentar áreas maiores primeiro.
	sort.Slice(b.freeSpaces, func(i, j int) bool {
		return b.freeSpaces[i].Volume() > b.freeSpaces[j].Volume()
	})

	b.Products = append(b.Products, packedProduct{ID: item.ProductID, Index: item.Index})
	return true
}

type OrderPackingResult struct {
	Boxes []PackedBox
}

// PackOrder empacota itens com uma heurística determinística para o problema NP-difícil de bin packing 3D; busca minimizar caixas abertas, mas não garante ótimo global.
// Retorna erro se algum item não couber em nenhuma caixa disponível.
func PackOrder(items []Item, boxTypes []BoxType, allowRotation bool) (OrderPackingResult, error) {
	if len(items) == 0 {
		return OrderPackingResult{Boxes: []PackedBox{}}, nil
	}

	// Problema NP-difícil tratado via heurística determinística para reduzir caixas abertas.

	// Ordena itens por volume decrescente (FFD) para que maiores ocupem primeiro, reduzindo fragmentação.
	for i := range items {
		items[i].Volume = items[i].Dim.Volume()
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Volume == items[j].Volume {
			return items[i].ProductID < items[j].ProductID
		}
		return items[i].Volume > items[j].Volume
	})

	// Ordena caixas por volume crescente para testar primeiro o menor recipiente viável.
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
			if opened[bi].TryPlace(it, allowRotation) {
				placed = true
				break
			}
		}

		if placed {
			continue
		}

		// Abrir nova caixa: prefere encaixar sem rotação; recorre à rotação se necessário e permitido.
		noRotationIdx := -1
		rotationIdx := -1

		for i := range boxTypes {
			bt := boxTypes[i]
			space := Dimensions{Height: bt.Height, Width: bt.Width, Length: bt.Length}

			if it.Dim.FitsIn(space) {
				noRotationIdx = i
				break
			}

			if allowRotation && rotationIdx == -1 {
				for _, rot := range it.Dim.Rotations() {
					if rot.FitsIn(space) {
						rotationIdx = i
						break
					}
				}
			}
		}

		chosenIdx := -1
		if noRotationIdx != -1 {
			chosenIdx = noRotationIdx
		} else if rotationIdx != -1 {
			chosenIdx = rotationIdx
		}

		if chosenIdx == -1 {
			if allowRotation {
				return OrderPackingResult{}, fmt.Errorf("produto '%s' não cabe em nenhuma caixa disponível (mesmo com rotação)", it.ProductID)
			}
			return OrderPackingResult{}, fmt.Errorf("produto '%s' não cabe em nenhuma caixa disponível", it.ProductID)
		}

		chosen := boxTypes[chosenIdx]

		nb := newPackedBox(chosen)
		if !nb.TryPlace(it, allowRotation) {
			// Não deve acontecer após a checagem de ajuste, mas mantemos validação defensiva.
			return OrderPackingResult{}, fmt.Errorf("falha inesperada ao alocar produto '%s' na caixa '%s'", it.ProductID, chosen.ID)
		}
		opened = append(opened, nb)
	}

	return OrderPackingResult{Boxes: opened}, nil
}
