package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(ctx context.Context, host, port string) (*Redis, error) {
	var tlsConfig *tls.Config
	options := &redis.Options{
		Addr:      fmt.Sprintf("%s:%s", host, port),
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

func (r Redis) Client() *redis.Client {
	return r.client
}
