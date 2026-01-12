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
	// Handler orquestra entrada HTTP e delega regra de negócio para o service.
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
	// Validação sintática/JSON ocorre no handler para responder 400 sem invocar o domínio.
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
			// ServiceError já traz o status apropriado definido pelas regras de negócio.
			c.JSON(se.StatusCode, gin.H{
				"error": gin.H{
					"code":    "PACKING_ERROR",
					"message": se.Message,
				},
			})
			return
		}

		// Fallback 500 para falhas inesperadas não mapeadas pelo service.
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "erro interno inesperado",
			},
		})
		return
	}

	// Service já garante preservação da ordem dos produtos; handler apenas serializa o DTO final.
	c.JSON(http.StatusOK, resp)
}
