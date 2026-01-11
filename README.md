# Packing Optimizer API

API em Go para otimizar o empacotamento de produtos em caixas disponíveis, buscando **minimizar o número de caixas** por pedido.
Inclui **rotação 3D (6 orientações)** e heurística determinística baseada em *First Fit Decreasing* com *free-spaces*.

## Stack
- Go 1.23+
- Gin (HTTP)
- Swagger (swaggo)
- Docker / Docker Compose

## Como rodar

### Local (Go)
```bash
go run ./cmd/api

A API sobe em:

Health: GET http://localhost:8080/healthz

Swagger: GET http://localhost:8080/swagger/index.html

Packing: POST http://localhost:8080/v1/packing

Docker (recomendado)
docker compose up --build

Endpoint principal
POST /v1/packing

Request

{
  "pedidos": [
    {
      "pedido_id": 1,
      "produtos": [
        {
          "produto_id": "PS5",
          "dimensoes": { "altura": 40, "largura": 10, "comprimento": 25 }
        }
      ]
    }
    Health: GET http://localhost:8080/healthz

    Swagger: GET http://localhost:8080/swagger/index.html

    Packing: POST http://localhost:8080/v1/packing

    ### Docker (recomendado)
    ```bash
    docker compose up --build
    ```

    ### Endpoint principal
    POST /v1/packing

    ### Request
    ```json
    {
      "pedidos": [
        {
          "pedido_id": 1,
          "produtos": [
            {
              "produto_id": "PS5",
              "dimensoes": { "altura": 40, "largura": 10, "comprimento": 25 }
            }
          ]
        }
      ]
    }
    ```

    ### Response (exemplo)
    ```json
    {
      "pedidos": [
        {
          "pedido_id": 1,
          "caixas": [
            { "caixa_id": "Caixa 1", "produtos": ["PS5"] }
          ]
        }
      ]
    }
    ```

    ## Decisões de projeto

    ### Rotação 3D

    O enunciado não restringe orientação dos produtos, então o algoritmo permite rotação 3D, testando até 6 permutações únicas de (altura, largura, comprimento).
    Isso aumenta o aproveitamento e pode reduzir o número total de caixas.

    ### Heurística de empacotamento

    O problema se aproxima de 3D bin packing (NP-difícil). Para manter desempenho e previsibilidade, foi adotada uma heurística:

    ordena produtos por volume decrescente (First Fit Decreasing);

    tenta encaixar em caixas já abertas, avaliando todas as rotações;

    mantém uma lista de espaços livres (free-spaces) por caixa e aplica um split determinístico ao inserir itens;

    se não couber, abre a menor caixa disponível que comporte o produto.

    ### Erros personalizados

    400 para erros de validação de JSON/estrutura;

    422 quando um produto não cabe em nenhuma caixa (mesmo com rotação), com mensagem contextualizada por pedido;

    500 para falhas inesperadas.

    ## Testes
    ```bash
    go test ./...
    ```

    ## Notas

    Peso, fragilidade, empilhamento e outras restrições não foram consideradas por não estarem especificadas.