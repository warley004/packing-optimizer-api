package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/warley004/packing-optimizer-api/internal/api/dto"
)

type PackingHandler struct{}

func NewPackingHandler() *PackingHandler {
	return &PackingHandler{}
}

func (h *PackingHandler) Pack(c *gin.Context) {
	var req dto.PackingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	// Placeholder determinístico (até implementarmos o algoritmo):
	// coloca todos os produtos do pedido na "Caixa 3".
	resp := dto.PackingResponse{Pedidos: make([]dto.PedidoResponse, 0, len(req.Pedidos))}
	for _, p := range req.Pedidos {
		prodIDs := make([]string, 0, len(p.Produtos))
		for _, pr := range p.Produtos {
			prodIDs = append(prodIDs, pr.ProdutoID)
		}

		resp.Pedidos = append(resp.Pedidos, dto.PedidoResponse{
			PedidoID: p.PedidoID,
			Caixas: []dto.CaixaResponse{
				{
					CaixaID:  "Caixa 3",
					Produtos: prodIDs,
				},
			},
		})
	}

	c.JSON(http.StatusOK, resp)
}
