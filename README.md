
# Desafio — Listagem de Orders (Clean Architecture)

Projeto em Go com **REST**, **gRPC** e **GraphQL** para **listar orders**. Inclui criação de order para facilitar os testes. Migrações são aplicadas automaticamente na subida do app.

## Stack
- Go 1.22
- Postgres 16
- gRPC (protoc)
- GraphQL (graphql-go)
- Router HTTP (chi)
- Driver PG (pgx)

## Como subir
Requisitos: Docker e Docker Compose.

```bash
docker compose up --build
```

Isso vai subir:
- **Postgres** em `localhost:5432`
- **REST** em `http://localhost:8080`
- **GraphQL** em `http://localhost:8081`
- **gRPC** em `localhost:50051`

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

## Arquitetura (resumo)
- `internal/domain` — entidades de domínio
- `internal/repository` — portas (interfaces) e adaptador Postgres
- `internal/usecase` — casos de uso (ListOrders, CreateOrder)
- `internal/adapter/http/rest` — servidor REST
- `internal/adapter/grpc` — servidor gRPC
- `internal/adapter/graphql` — servidor GraphQL
- `internal/infra/db` — conexão e migrações embutidas
- `cmd/app` — composição dos serviços

## Variáveis de ambiente
Definidas no `docker-compose.yaml`. Padrões:
- `REST_PORT=8080`
- `GRAPHQL_PORT=8081`
- `GRPC_PORT=50051`

## Notas
- O build usa **multi-stage** e gera os stubs gRPC com `protoc` durante o build.
- Migrações ficam em `internal/infra/db/migrations` e são aplicadas automaticamente.
