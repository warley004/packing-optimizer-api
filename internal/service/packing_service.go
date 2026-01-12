package service

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/warley004/packing-optimizer-api/internal/api/dto"
	"github.com/warley004/packing-optimizer-api/internal/packing"
)

type PackingService struct {
	boxes []packing.BoxType
}

// Service consolida regras de domínio de empacotamento; handlers apenas transformam HTTP <-> DTO e delegam aqui.
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
		for idx, p := range pedido.Produtos {
			items = append(items, packing.Item{
				ProductID: p.ProdutoID,
				Dim: packing.Dimensions{
					Height: p.Dimensoes.Altura,
					Width:  p.Dimensoes.Largura,
					Length: p.Dimensoes.Comprimento,
				},
				Index: idx, // preserva ordem do input
			})
		}

		// Permite rotação 3D para o algoritmo aproveitar melhor cada caixa.
		result, err := packing.PackOrder(items, s.boxes, true)
		if err != nil {
			// Handler converte ServiceError em resposta HTTP (status e payload).
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
			// Ordena produtos para manter a ordem recebida no payload original.
			sort.Slice(b.Products, func(i, j int) bool {
				return b.Products[i].Index < b.Products[j].Index
			})

			ids := make([]string, 0, len(b.Products))
			for _, pp := range b.Products {
				ids = append(ids, pp.ID)
			}

			pr.Caixas = append(pr.Caixas, dto.CaixaResponse{
				CaixaID:  b.BoxType.ID,
				Produtos: ids,
			})
		}

		resp.Pedidos = append(resp.Pedidos, pr)
	}

	return resp, nil
}

func formatID(id int64) string {
	return fmt.Sprintf("%d", id)
}
