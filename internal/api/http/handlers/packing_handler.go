package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/warley004/packing-optimizer-api/internal/api/dto"
	"github.com/warley004/packing-optimizer-api/internal/service"
)

type PackingHandler struct {
	service *service.PackingService
}

func NewPackingHandler() *PackingHandler {
	return &PackingHandler{
		service: service.NewPackingService(),
	}
}

// Pack godoc
// @Summary      Empacotar pedidos
// @Description  Processa uma lista de pedidos e retorna a alocação de produtos em caixas disponíveis (minimizando o número de caixas).
// @Tags         packing
// @Accept       json
// @Produce      json
// @Param        request  body      dto.PackingRequest  true  "Lista de pedidos com produtos e dimensões"
// @Success      200      {object}  dto.PackingResponse
// @Failure      400      {object}  map[string]any  "Erro de validação do JSON/estrutura"
// @Failure      422      {object}  map[string]any  "Erro de empacotamento (produto não cabe)"
// @Failure      500      {object}  map[string]any  "Erro interno"
// @Router       /v1/packing [post]
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

	resp, err := h.service.Pack(req)
	if err != nil {
		if se, ok := err.(*service.ServiceError); ok {
			c.JSON(se.StatusCode, gin.H{
				"error": gin.H{
					"code":    "PACKING_ERROR",
					"message": se.Message,
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "erro interno inesperado",
			},
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
