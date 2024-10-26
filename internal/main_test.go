package internal_test

import (
	"context"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"go-hft-orderbook/internal"
	"go-hft-orderbook/internal/datasources"
	"go-hft-orderbook/internal/datasources/redis"
	"log"
	"os"
	"testing"
)

var cache internal.Cache
var mr *miniredis.Miniredis

func setup() {
	mr, err := miniredis.Run()
	if err != nil {
		panic(fmt.Sprintf("could not start miniredis: %s", err.Error()))
	}

	log.Println("miniredis running on:", mr.Port())

	redisClient, err := redis.NewRedis(context.Background(), mr.Host(), mr.Port())
	if err != nil {
		panic("error setting redis client")
	}
	cache = datasources.NewCacheRepository(redisClient.Client())
}

func teardown() {
	if mr != nil {
		mr.Close()
	}
}

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}
