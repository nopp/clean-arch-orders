package usecase

import "github.com/nopp/clean-arch-orders/internal/repository"

type ListOrders struct {
	repo repository.OrderRepository
}

func NewListOrders(repo repository.OrderRepository) *ListOrders {
	return &ListOrders{repo: repo}
}

func (uc *ListOrders) Execute() (any, error) {
	return uc.repo.List()
}
