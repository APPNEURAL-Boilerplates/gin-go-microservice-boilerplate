package items

import (
	"context"
	"sort"
	"sync"

	"github.com/example/gin-microservice-boilerplate/internal/common"
)

type Repository interface {
	List(ctx context.Context) ([]Item, error)
	Get(ctx context.Context, id string) (Item, error)
	Create(ctx context.Context, item Item) (Item, error)
}

type MemoryRepository struct {
	mu    sync.RWMutex
	items map[string]Item
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		items: make(map[string]Item),
	}
}

func (r *MemoryRepository) List(ctx context.Context) ([]Item, error) {
	_ = ctx

	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]Item, 0, len(r.items))
	for _, item := range r.items {
		result = append(result, item)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})

	return result, nil
}

func (r *MemoryRepository) Get(ctx context.Context, id string) (Item, error) {
	_ = ctx

	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.items[id]
	if !exists {
		return Item{}, common.NotFound("item not found")
	}

	return item, nil
}

func (r *MemoryRepository) Create(ctx context.Context, item Item) (Item, error) {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	r.items[item.ID] = item
	return item, nil
}
