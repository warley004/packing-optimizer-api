package service

import (
	"fmt"
	"net/http"
	"runtime"
	"sort"
	"sync"

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
	total := len(req.Pedidos)
	resp := dto.PackingResponse{
		Pedidos: make([]dto.PedidoResponse, total),
	}

	if total == 0 {
		return resp, nil
	}

	type job struct {
		index  int
		pedido dto.PedidoRequest
	}

	type result struct {
		index  int
		pedido dto.PedidoResponse
		err    error
	}

	jobCh := make(chan job)
	resultCh := make(chan result, total)

	workerCount := runtime.NumCPU()
	if workerCount > total {
		workerCount = total
	}
	if workerCount < 1 {
		workerCount = 1
	}

	var wg sync.WaitGroup
	wg.Add(workerCount)

	for i := 0; i < workerCount; i++ {
		go func() {
			defer wg.Done()
			for j := range jobCh {
				pedidoResp, err := s.packSingleOrder(j.pedido)
				resultCh <- result{index: j.index, pedido: pedidoResp, err: err}
			}
		}()
	}

	go func() {
		for idx, pedido := range req.Pedidos {
			jobCh <- job{index: idx, pedido: pedido}
		}
		close(jobCh)
	}()

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	errors := make([]error, total)
	for res := range resultCh {
		if res.err != nil {
			errors[res.index] = res.err
			continue
		}
		resp.Pedidos[res.index] = res.pedido
	}

	for idx := 0; idx < total; idx++ {
		if errors[idx] != nil {
			return dto.PackingResponse{}, errors[idx]
		}
	}

	return resp, nil
}

func formatID(id int64) string {
	return fmt.Sprintf("%d", id)
}

func (s *PackingService) packSingleOrder(pedido dto.PedidoRequest) (dto.PedidoResponse, error) {
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

	result, err := packing.PackOrder(items, s.boxes, true)
	if err != nil {
		return dto.PedidoResponse{}, &ServiceError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Pedido " + formatID(pedido.PedidoID) + ": " + err.Error(),
		}
	}

	pr := dto.PedidoResponse{
		PedidoID: pedido.PedidoID,
		Caixas:   make([]dto.CaixaResponse, 0, len(result.Boxes)),
	}

	for _, b := range result.Boxes {
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

	return pr, nil
}
