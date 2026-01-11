package packing

type BoxType struct {
	ID          string
	Height      int
	Width       int
	Length      int
}

func AvailableBoxes() []BoxType {
	return []BoxType{
		{ID: "Caixa 1", Height: 30, Width: 40, Length: 80},
		{ID: "Caixa 2", Height: 50, Width: 50, Length: 40},
		{ID: "Caixa 3", Height: 50, Width: 80, Length: 60},
	}
}
