package items

import (
	"context"
	"time"

	"github.com/example/gin-go-microservice/internal/common"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) List(ctx context.Context) ([]Item, error) {
	return s.repository.List(ctx)
}

func (s *Service) Get(ctx context.Context, id string) (Item, error) {
	if id == "" {
		return Item{}, common.BadRequest("item id is required")
	}

	return s.repository.Get(ctx, id)
}

func (s *Service) Create(ctx context.Context, request CreateItemRequest) (Item, error) {
	now := time.Now().UTC()
	item := Item{
		ID:          common.NewID("item"),
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return s.repository.Create(ctx, item)
}
