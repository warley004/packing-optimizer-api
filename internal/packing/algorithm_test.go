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

// Modo enunciado: sem rotação por padrão.
func TestPackOrder_SingleItemChoosesSmallestBoxThatFits_NoRotation(t *testing.T) {
	items := []Item{
		{ProductID: "PS5", Dim: Dimensions{Height: 40, Width: 10, Length: 25}, Index: 0},
	}

	res, err := PackOrder(items, AvailableBoxes(), false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Boxes) != 1 {
		t.Fatalf("expected 1 box, got %d", len(res.Boxes))
	}

	// Sem rotação, PS5 (altura 40) não cabe na Caixa 1 (altura 30), então deve ir para Caixa 2.
	if res.Boxes[0].BoxType.ID != "Caixa 2" {
		t.Fatalf("expected Caixa 2, got %s", res.Boxes[0].BoxType.ID)
	}
}

// Diferencial: com rotação, o item pode caber em caixas menores.
func TestPackOrder_UsesRotationToFit_WhenEnabled(t *testing.T) {
	items := []Item{
		{ProductID: "Rot", Dim: Dimensions{Height: 80, Width: 30, Length: 40}, Index: 0},
	}

	res, err := PackOrder(items, AvailableBoxes(), true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Boxes) != 1 {
		t.Fatalf("expected 1 box, got %d", len(res.Boxes))
	}

	// Com rotação habilitada, 80x30x40 cabe na Caixa 1 (30x40x80).
	if res.Boxes[0].BoxType.ID != "Caixa 1" {
		t.Fatalf("expected Caixa 1, got %s", res.Boxes[0].BoxType.ID)
	}
}

// Este teste garante o comportamento do exemplo do enunciado (Pedido 1: Caixa 2 com PS5 e Volante).
func TestPackOrder_MatchesSpecExample_Pedido1_NoRotation(t *testing.T) {
	items := []Item{
		{ProductID: "PS5", Dim: Dimensions{Height: 40, Width: 10, Length: 25}, Index: 0},
		{ProductID: "Volante", Dim: Dimensions{Height: 40, Width: 30, Length: 30}, Index: 1},
	}

	res, err := PackOrder(items, AvailableBoxes(), false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Boxes) != 1 {
		t.Fatalf("expected 1 box, got %d", len(res.Boxes))
	}

	if res.Boxes[0].BoxType.ID != "Caixa 2" {
		t.Fatalf("expected Caixa 2, got %s", res.Boxes[0].BoxType.ID)
	}

	// A ordem dos produtos na caixa será normalizada no service (pela ordem de input).
	// Aqui a gente só valida a caixa escolhida e o fato de conter 2 itens.
	if len(res.Boxes[0].Products) != 2 {
		t.Fatalf("expected 2 products, got %d", len(res.Boxes[0].Products))
	}
}

func TestPackOrder_TwoItemsSameBoxWhenPossible_NoRotation(t *testing.T) {
	items := []Item{
		{ProductID: "A", Dim: Dimensions{Height: 10, Width: 10, Length: 10}, Index: 0},
		{ProductID: "B", Dim: Dimensions{Height: 10, Width: 10, Length: 10}, Index: 1},
	}

	res, err := PackOrder(items, AvailableBoxes(), false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Boxes) != 1 {
		t.Fatalf("expected 1 box, got %d", len(res.Boxes))
	}
}

func TestPackOrder_ItemUnpackableReturnsError_NoRotation(t *testing.T) {
	items := []Item{
		{ProductID: "GIGANTE", Dim: Dimensions{Height: 999, Width: 999, Length: 999}, Index: 0},
	}

	_, err := PackOrder(items, AvailableBoxes(), false)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
