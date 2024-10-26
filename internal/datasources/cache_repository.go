package datasources

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"go-hft-orderbook/internal"
)

var ErrKeyNotFound = errors.New("key not found")

type CacheRepository struct {
	redisClient *redis.Client
}

func NewCacheRepository(redisClient *redis.Client) *CacheRepository {
	return &CacheRepository{
		redisClient: redisClient,
	}
}

func (r CacheRepository) Enqueue(ctx context.Context, price string, o *internal.Order) error {
	bytes, err := json.Marshal(o)
	if err != nil {
		return err
	}

	if err = r.redisClient.RPush(ctx, price, bytes).Err(); err != nil {
		return err
	}
	return nil
}

func (r CacheRepository) Dequeue(ctx context.Context, price string) (*internal.Order, error) {
	res, err := r.redisClient.LPop(ctx, price).Bytes()
	if res == nil {
		return nil, ErrKeyNotFound
	}
	if err != nil {
		return nil, err
	}

	var order internal.Order
	if err = json.Unmarshal(res, &order); err != nil {
		return nil, err
	}
	return &order, nil
}

func (r CacheRepository) Size(ctx context.Context, price string) (int, error) {
	res, err := r.redisClient.LLen(ctx, price).Result()
	if err != nil {
		return 0, err
	}
	return int(res), nil
}

func (r CacheRepository) Delete(ctx context.Context, price string, o *internal.Order) error {
	bytes, err := json.Marshal(o)
	if err != nil {
		return err
	}

	err = r.redisClient.LRem(ctx, price, 0, bytes).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r CacheRepository) DeleteAll(ctx context.Context, price string) error {
	err := r.redisClient.Del(ctx, price).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r CacheRepository) IsEmpty(ctx context.Context, price string) (bool, error) {
	size, err := r.Size(ctx, price)
	return size == 0, err
}
