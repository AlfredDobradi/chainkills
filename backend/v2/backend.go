package v2

import (
	"context"
	"time"

	"git.sr.ht/~barveyhirdman/chainkills/backend/v2/memory"
	"git.sr.ht/~barveyhirdman/chainkills/config"
	"github.com/redis/go-redis/v9"
)

type Engine interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	SMembers(ctx context.Context, key string) *redis.StringSliceCmd
	SAdd(ctx context.Context, key string, members ...interface{}) *redis.IntCmd
}

func Get(kind string) Engine {
	switch kind {
	case "memory":
		return memory.New()
	case "redis":
		return redis.NewClient(&redis.Options{
			Addr: config.Get().Redict.Address,
			DB:   config.Get().Redict.Database,
		})
	}

	return nil
}
