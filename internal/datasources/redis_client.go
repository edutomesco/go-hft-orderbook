package datasources

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var ErrKeyNotFound = errors.New("key not found")

type Cache interface {
	Set(ctx context.Context, key string, data []byte) error
	SAdd(ctx context.Context, key string, member string) error
	Get(ctx context.Context, key string) ([]byte, error)
	GetMembersByKey(ctx context.Context, key string) ([]string, error)
	Del(ctx context.Context, key string) error
}

type Redis struct {
	client *redis.Client
}

func NewRedis(ctx context.Context, host, port, password string, db int) (*Redis, error) {
	var tlsConfig *tls.Config
	options := &redis.Options{
		Addr:      fmt.Sprintf("%s:%s", host, port),
		Password:  password,
		DB:        db,
		TLSConfig: tlsConfig,
	}
	client := redis.NewClient(options)
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &Redis{
		client: client,
	}, nil
}

func (r Redis) Set(ctx context.Context, key string, data []byte) error {
	if err := r.client.Set(ctx, key, data, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (r Redis) SAdd(ctx context.Context, key string, member string) error {
	if err := r.client.SAdd(ctx, fmt.Sprintf("\"%s\"", key), member).Err(); err != nil {
		return err
	}
	return nil
}

func (r Redis) Get(ctx context.Context, key string) ([]byte, error) {
	res, err := r.client.Get(ctx, key).Bytes()
	if res == nil {
		return nil, ErrKeyNotFound
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r Redis) GetMembersByKey(ctx context.Context, key string) ([]string, error) {
	res, err := r.client.SMembers(ctx, fmt.Sprintf("\"%s\"", key)).Result()
	if res == nil {
		return nil, ErrKeyNotFound
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r Redis) Del(ctx context.Context, key string) error {
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}

func (r Redis) Client() *redis.Client {
	return r.client
}
