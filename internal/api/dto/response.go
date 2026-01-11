package dto

type PackingResponse struct {
	Pedidos []PedidoResponse `json:"pedidos"`
}

type PedidoResponse struct {
	PedidoID int64          `json:"pedido_id"`
	Caixas   []CaixaResponse `json:"caixas"`
}

type CaixaResponse struct {
	CaixaID   string   `json:"caixa_id"`
	Produtos  []string `json:"produtos"`
}
