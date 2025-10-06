
package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/example/clean-arch-orders/internal/domain"
	"github.com/example/clean-arch-orders/internal/repository"
	"github.com/jackc/pgx/v5"
)

var _ repository.OrderRepository = (*OrderPostgres)(nil)

type OrderPostgres struct {
	Conn *pgx.Conn
}

func New(conn *pgx.Conn) *OrderPostgres { return &OrderPostgres{Conn: conn} }

func (r *OrderPostgres) Create(o *domain.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.Conn.Exec(ctx, `INSERT INTO orders (id, customer_name, total_amount, created_at)
		VALUES ($1,$2,$3,$4)`, o.ID, o.CustomerName, o.TotalAmount, o.CreatedAt.UTC())
	return err
}

func (r *OrderPostgres) List() ([]domain.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.Conn.Query(ctx, `SELECT id, customer_name, total_amount, created_at
		FROM orders ORDER BY created_at DESC`)
	if err != nil { return nil, err }
	defer rows.Close()

	var out []domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(&o.ID, &o.CustomerName, &o.TotalAmount, &o.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, rows.Err()
}
