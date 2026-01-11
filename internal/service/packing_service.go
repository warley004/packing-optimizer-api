package service
import "fmt"

import (
	"net/http"

	"github.com/warley004/packing-optimizer-api/internal/api/dto"
	"github.com/warley004/packing-optimizer-api/internal/packing"
)

type PackingService struct {
	boxes []packing.BoxType
}

func NewPackingService() *PackingService {
	return &PackingService{
		boxes: packing.AvailableBoxes(),
	}
}

type ServiceError struct {
	StatusCode int
	Message    string
}

func (e *ServiceError) Error() string {
	return e.Message
}

func (s *PackingService) Pack(req dto.PackingRequest) (dto.PackingResponse, error) {
	resp := dto.PackingResponse{
		Pedidos: make([]dto.PedidoResponse, 0, len(req.Pedidos)),
	}

	for _, pedido := range req.Pedidos {
		items := make([]packing.Item, 0, len(pedido.Produtos))
		for _, p := range pedido.Produtos {
			items = append(items, packing.Item{
				ProductID: p.ProdutoID,
				Dim: packing.Dimensions{
					Height: p.Dimensoes.Altura,
					Width:  p.Dimensoes.Largura,
					Length: p.Dimensoes.Comprimento,
				},
			})
		}

		result, err := packing.PackOrder(items, s.boxes)
		if err != nil {
			return dto.PackingResponse{}, &ServiceError{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Pedido " + formatID(pedido.PedidoID) + ": " + err.Error(),
			}
		}

		pr := dto.PedidoResponse{
			PedidoID: pedido.PedidoID,
			Caixas:   make([]dto.CaixaResponse, 0, len(result.Boxes)),
		}

		for _, b := range result.Boxes {
			pr.Caixas = append(pr.Caixas, dto.CaixaResponse{
				CaixaID:  b.BoxType.ID,
				Produtos: b.Products,
			})
		}

		resp.Pedidos = append(resp.Pedidos, pr)
	}

	return resp, nil
}

func formatID(id int64) string {
	return fmt.Sprintf("%d", id)
}
