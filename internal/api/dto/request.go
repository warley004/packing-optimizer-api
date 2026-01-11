package dto

type PackingRequest struct {
	Pedidos []PedidoRequest `json:"pedidos" binding:"required,min=1"`
}

type PedidoRequest struct {
	PedidoID  int64            `json:"pedido_id" binding:"required"`
	Produtos  []ProdutoRequest `json:"produtos" binding:"required,min=1"`
}

type ProdutoRequest struct {
	ProdutoID  string        `json:"produto_id" binding:"required,min=1"`
	Dimensoes  DimensoesDTO  `json:"dimensoes" binding:"required"`
}

type DimensoesDTO struct {
	Altura      int `json:"altura" binding:"required,gt=0"`
	Largura     int `json:"largura" binding:"required,gt=0"`
	Comprimento int `json:"comprimento" binding:"required,gt=0"`
}
