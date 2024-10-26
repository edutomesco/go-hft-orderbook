package internal

import "context"

type Cache interface {
	Enqueue(ctx context.Context, price string, o *Order) error
	Dequeue(ctx context.Context, price string) (*Order, error)
	Size(ctx context.Context, price string) (int, error)
	Delete(ctx context.Context, price string, o *Order) error
	DeleteAll(ctx context.Context, price string) error
	IsEmpty(ctx context.Context, price string) (bool, error)
}
