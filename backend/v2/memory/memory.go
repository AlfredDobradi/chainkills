package memory

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Store struct {
	keyValue map[any]any
}

func New() *Store {
	return &Store{
		keyValue: make(map[any]any),
	}
}

func (s *Store) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	s.keyValue[key] = value

	return redis.NewStatusCmd(ctx, "OK")
}

func (s *Store) Get(ctx context.Context, key string) *redis.StringCmd {
	if value, ok := s.keyValue[key]; ok {
		return redis.NewStringResult(value.(string), nil)
	}
	return redis.NewStringResult("", redis.Nil)
}

func (s *Store) SMembers(ctx context.Context, key string) *redis.StringSliceCmd {
	result := make([]string, 0)
	ss, ok := s.keyValue[key]
	if !ok {
		return redis.NewStringSliceResult(result, redis.Nil)
	}

	value, ok := ss.(map[any]struct{})
	if !ok {
		return redis.NewStringSliceResult(result, redis.Nil)
	}

	members := make([]string, 0, len(value))
	for k := range value {
		members = append(members, fmt.Sprintf("%v", k))
	}
	return redis.NewStringSliceResult(members, nil)
}

func (s *Store) SAdd(ctx context.Context, key string, members ...any) *redis.IntCmd {
	ss, ok := s.keyValue[key]
	if !ok {
		ss = make(map[any]struct{})
		s.keyValue[key] = ss
	}

	value, ok := ss.(map[any]struct{})
	if !ok {
		return redis.NewIntResult(0, redis.Nil)
	}

	var delta int64 = 0
	for _, m := range members {
		if _, ok := value[m]; !ok {
			delta++
		}
		value[m] = struct{}{}
	}

	return redis.NewIntResult(delta, nil)
}
