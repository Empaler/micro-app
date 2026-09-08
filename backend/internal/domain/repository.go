package domain

import (
	"context"
)

type Repository[T any] interface {
	Create(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id int64) (*T, error)
	GetAll(ctx context.Context) ([]T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id int64) error
}
