package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	graphqlsrv "github.com/nopp/clean-arch-orders/internal/adapter/graphql"
	grpcsrv "github.com/nopp/clean-arch-orders/internal/adapter/grpc"
	"github.com/nopp/clean-arch-orders/internal/adapter/grpc/pb"
	rest "github.com/nopp/clean-arch-orders/internal/adapter/http/rest"
	"github.com/nopp/clean-arch-orders/internal/infra/db"
	"github.com/nopp/clean-arch-orders/internal/repository/postgres"
	"github.com/nopp/clean-arch-orders/internal/usecase"
	"google.golang.org/grpc"
)

func mustEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func main() {
	// DB
	connStr := db.ConnString(
		mustEnv("DB_HOST", "postgres"),
		mustEnv("DB_PORT", "5432"),
		mustEnv("DB_USER", "app"),
		mustEnv("DB_PASS", "app"),
		mustEnv("DB_NAME", "ordersdb"),
		mustEnv("DB_SSLMODE", "disable"),
	)
	conn, err := db.Connect(connStr)
	if err != nil {
		log.Fatalf("erro ao conectar no banco: %v", err)
	}
	defer conn.Close(context.Background())

	if err := db.RunMigrations(conn); err != nil {
		log.Fatalf("migrações falharam: %v", err)
	}

	repo := postgres.New(conn)
	listUC := usecase.NewListOrders(repo)
	createUC := usecase.NewCreateOrder(repo)

	// REST
	restSrv := rest.NewServer(listUC, createUC)
	restPort := mustEnv("REST_PORT", "8080")
	restHTTP := &http.Server{Addr: ":" + restPort, Handler: restSrv.Router()}

	// GraphQL
	graphqlSrv := graphqlsrv.NewServer(listUC)
	graphqlPort := mustEnv("GRAPHQL_PORT", "8081")
	graphqlHTTP := &http.Server{Addr: ":" + graphqlPort, Handler: graphqlSrv.Handler()}

	// gRPC
	grpcPort := mustEnv("GRPC_PORT", "50051")
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, grpcsrv.NewServer(listUC))

	// Start servers
	go func() { log.Printf("REST ouvindo em :%s", restPort); log.Println(restHTTP.ListenAndServe()) }()
	go func() { log.Printf("GraphQL ouvindo em :%s", graphqlPort); log.Println(graphqlHTTP.ListenAndServe()) }()
	go func() { log.Printf("gRPC ouvindo em :%s", grpcPort); log.Println(grpcServer.Serve(lis)) }()

	// Graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = restHTTP.Shutdown(ctx)
	_ = graphqlHTTP.Shutdown(ctx)
	grpcServer.GracefulStop()
	_ = conn.Close(ctx)
	log.Println("bye")
}

// expose _ to prevent unused import of pgx/v5 in go.mod checksum when only indrectly used
var _ pgx.Row
