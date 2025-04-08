package backend

import (
	"context"
	"time"

	"git.sr.ht/~barveyhirdman/chainkills/backend/memory"
	"git.sr.ht/~barveyhirdman/chainkills/config"
	"github.com/redis/go-redis/v9"
)

type Engine interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	SMembers(ctx context.Context, key string) *redis.StringSliceCmd
	SAdd(ctx context.Context, key string, members ...interface{}) *redis.IntCmd
}

func Get() Engine {
	switch config.Get().Backend.Kind {
	default:
		fallthrough
	case "memory":
		return memory.New()
	case "redis":
		return redis.NewClient(&redis.Options{
			Addr: config.Get().Backend.Address,
			DB:   config.Get().Backend.Database,
		})
	}
}
