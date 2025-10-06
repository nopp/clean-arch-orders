
package grpcsrv

import (
	"context"
	"time"

	"github.com/example/clean-arch-orders/internal/adapter/grpc/pb"
	"github.com/example/clean-arch-orders/internal/domain"
	"github.com/example/clean-arch-orders/internal/usecase"
)

type Server struct {
	pb.UnimplementedOrderServiceServer
	ListUC *usecase.ListOrders
}

func NewServer(list *usecase.ListOrders) *Server { return &Server{ListUC: list} }

func (s *Server) ListOrders(ctx context.Context, in *pb.Empty) (*pb.ListOrdersResponse, error) {
	res, err := s.ListUC.Execute()
	if err != nil { return nil, err }
	orders := res.([]domain.Order)
	out := &pb.ListOrdersResponse{}
	for _, o := range orders {
		out.Orders = append(out.Orders, &pb.Order{
			Id: o.ID, CustomerName: o.CustomerName, TotalAmount: o.TotalAmount, CreatedAt: o.CreatedAt.Format(time.RFC3339),
		})
	}
	return out, nil
}
