APP_NAME := clean-arch-orders
DOCKER_COMPOSE := docker compose
GO := go

.PHONY: build run stop restart logs proto clean

## build: compila o binário da aplicação
build:
	@echo "🛠️  Compilando aplicação..."
	$(GO) build -o bin/$(APP_NAME) ./cmd/app
	@echo "✅ Build concluído -> bin/$(APP_NAME)"

## run: sobe toda a stack com Docker Compose
run:
	@echo "Subindo containers..."
	$(DOCKER_COMPOSE) up --build -d
	@echo "Aplicação pronta para uso, conforme as urls abaixo!"
	@echo "REST:     http://localhost:8080"
	@echo "GraphQL:  http://localhost:8081"
	@echo "gRPC:     localhost:50051"

## stop: para e remove os containers
stop:
	@echo "Parando containers..."
	$(DOCKER_COMPOSE) down
	@echo "Containers parados."

## restart: reinicia a aplicação rapidamente
restart: stop run

## logs: mostra logs em tempo real da aplicação
logs:
	@$(DOCKER_COMPOSE) logs -f app

## clean: remove binários e dependências temporárias
clean:
	@echo "Limpando o build..."
	rm -rf bin/
	@echo "Done!"

## proto: gera o código gRPC a partir do arquivo .proto
proto:
	@echo "Gerando código gRPC..."
	protoc --go_out=. --go-grpc_out=. api/proto/order.proto
	@echo "gRPC gerado!"