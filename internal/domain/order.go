
package domain

import "time"

type Order struct {
	ID           string
	CustomerName string
	TotalAmount  float64
	CreatedAt    time.Time
}
