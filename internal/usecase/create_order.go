
package usecase

import (
	"time"

	"github.com/example/clean-arch-orders/internal/domain"
	"github.com/example/clean-arch-orders/internal/repository"
)

type CreateOrder struct{ repo repository.OrderRepository }

func NewCreateOrder(r repository.OrderRepository) *CreateOrder { return &CreateOrder{repo: r} }

func (uc *CreateOrder) Execute(id, customer string, total float64, now time.Time) error {
	return uc.repo.Create(&domain.Order{
		ID:           id,
		CustomerName: customer,
		TotalAmount:  total,
		CreatedAt:    now,
	})
}
