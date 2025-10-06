
package repository

import "github.com/example/clean-arch-orders/internal/domain"

type OrderRepository interface {
	Create(o *domain.Order) error
	List() ([]domain.Order, error)
}
