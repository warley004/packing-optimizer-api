package packing

import "testing"

func TestDimensionsRotations(t *testing.T) {
	d := Dimensions{Height: 10, Width: 20, Length: 30}
	rots := d.Rotations()

	if len(rots) != 6 {
		t.Fatalf("expected 6 unique rotations, got %d", len(rots))
	}
}

func TestDimensionsRotationsWithDuplicates(t *testing.T) {
	d := Dimensions{Height: 10, Width: 10, Length: 20}
	rots := d.Rotations()

	if len(rots) != 3 {
		t.Fatalf("expected 3 unique rotations, got %d", len(rots))
	}
}

func TestPackOrder_SingleItemChoosesSmallestBoxThatFits(t *testing.T) {
	items := []Item{
		{ProductID: "PS5", Dim: Dimensions{Height: 40, Width: 10, Length: 25}},
	}

	res, err := PackOrder(items, AvailableBoxes())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Boxes) != 1 {
		t.Fatalf("expected 1 box, got %d", len(res.Boxes))
	}

	if res.Boxes[0].BoxType.ID != "Caixa 1" {
	t.Fatalf("expected Caixa 1, got %s", res.Boxes[0].BoxType.ID)
	}
}

func TestPackOrder_UsesRotationToFit(t *testing.T) {
	items := []Item{
		{ProductID: "Rot", Dim: Dimensions{Height: 80, Width: 30, Length: 40}},
	}

	res, err := PackOrder(items, AvailableBoxes())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Boxes) != 1 {
		t.Fatalf("expected 1 box, got %d", len(res.Boxes))
	}

	if res.Boxes[0].BoxType.ID != "Caixa 1" {
		t.Fatalf("expected Caixa 1, got %s", res.Boxes[0].BoxType.ID)
	}
}

func TestPackOrder_TwoItemsSameBoxWhenPossible(t *testing.T) {
	items := []Item{
		{ProductID: "A", Dim: Dimensions{Height: 10, Width: 10, Length: 10}},
		{ProductID: "B", Dim: Dimensions{Height: 10, Width: 10, Length: 10}},
	}

	res, err := PackOrder(items, AvailableBoxes())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Boxes) != 1 {
		t.Fatalf("expected 1 box, got %d", len(res.Boxes))
	}
}

func TestPackOrder_ItemUnpackableReturnsError(t *testing.T) {
	items := []Item{
		{ProductID: "GIGANTE", Dim: Dimensions{Height: 999, Width: 999, Length: 999}},
	}

	_, err := PackOrder(items, AvailableBoxes())
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
