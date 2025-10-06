
# Desafio — Listagem de Orders (Clean Architecture)

Projeto em Go com **REST**, **gRPC** e **GraphQL** para **listar orders**.

As migrações rodam automaticamente e criam a tabela `orders`.

## Endpoints

### REST
- `POST /order` — cria uma order
- `GET /order` — lista as orders

Exemplos no arquivo `api.http`.

### GraphQL
- `POST /` com `{ "query": "{ listOrders { id customer_name total_amount created_at } }" }`
- Ou `GET /?query=...`

### gRPC
Serviço: `OrderService` com método `ListOrders`.

Proto em `api/proto/order.proto`.

Exemplo com `evans` (CLI gRPC):

```bash
evans --host localhost --port 50051 -r repl
> package orderpb
> service OrderService
> call ListOrders
{}
```
