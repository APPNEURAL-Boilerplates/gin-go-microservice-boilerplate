package items

import "time"

type Item struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateItemRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=120"`
	Description string  `json:"description" binding:"omitempty,max=500"`
	Price       float64 `json:"price" binding:"required,gt=0"`
}
