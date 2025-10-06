
FROM golang:latest as builder

RUN apt-get update && apt-get install -y protobuf-compiler && rm -rf /var/lib/apt/lists/*
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \\
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN protoc --go_out=. --go-grpc_out=. api/proto/order.proto
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/app ./cmd/app

# Imagem final
FROM debian:stable-slim

WORKDIR /
COPY --from=builder /out/app /app
ENV REST_PORT=8080 GRAPHQL_PORT=8081 GRPC_PORT=50051 DB_HOST=postgres DB_PORT=5432 DB_USER=app DB_PASS=app DB_NAME=ordersdb DB_SSLMODE=disable
EXPOSE 8080 8081 50051
ENTRYPOINT ["/app"]