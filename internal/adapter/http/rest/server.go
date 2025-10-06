
package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/example/clean-arch-orders/internal/usecase"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	ListUC   *usecase.ListOrders
	CreateUC *usecase.CreateOrder
}

func NewServer(list *usecase.ListOrders, create *usecase.CreateOrder) *Server {
	return &Server{ListUC: list, CreateUC: create}
}

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/order", s.listOrders)
	r.Post("/order", s.createOrder)
	return r
}

func (s *Server) listOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := s.ListUC.Execute()
	if err != nil { http.Error(w, err.Error(), 500); return }
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

type createReq struct {
	ID           string  `json:"id"`
	CustomerName string  `json:"customer_name"`
	TotalAmount  float64 `json:"total_amount"`
}

func (s *Server) createOrder(w http.ResponseWriter, r *http.Request) {
	var req createReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", 400)
		return
	}
	if req.ID == "" || req.CustomerName == "" { http.Error(w, "campos obrigatórios", 400); return }
	if err := s.CreateUC.Execute(req.ID, req.CustomerName, req.TotalAmount, time.Now()); err != nil {
		http.Error(w, err.Error(), 500); return
	}
	w.WriteHeader(http.StatusCreated)
}
