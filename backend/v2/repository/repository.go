package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	v2 "git.sr.ht/~barveyhirdman/chainkills/backend/v2"
	"git.sr.ht/~barveyhirdman/chainkills/config"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

var packageName string = "git.sr.ht/~barveyhirdman/chainkills/backend/redict"

const (
	spanAddKillmail         = "AddKillmail"
	spanKillmailExists      = "KillmailExists"
	spanGetIgnoredSystemIDs = "GetIgnoredSystemIDs"
	spanGetIgnoredRegionIDs = "GetIgnoredRegionIDs"
	spanIgnoreSystemID      = "IgnoreSystemID"
	spanIgnoreRegionID      = "IgnoreRegionID"

	keyIgnoredSystemIDs = "ignored_system_ids"
	keyIgnoredRegionIDs = "ignored_region_ids"
)

type Backend struct {
	store v2.Engine
}

func New(engine v2.Engine) (*Backend, error) {
	return &Backend{
		store: engine,
	}, nil
}

func (r *Backend) AddKillmail(ctx context.Context, id string) error {
	_, span := otel.Tracer(packageName).Start(ctx, spanAddKillmail)
	defer span.End()

	span.SetAttributes(attribute.String("id", id))

	key := fmt.Sprintf("%s:%s", config.Get().Redict.Prefix, id)
	if err := r.store.Set(context.Background(), key, "", time.Duration(config.Get().Redict.TTL)*time.Minute).Err(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	span.SetStatus(codes.Ok, "ok")

	return nil
}

func (r *Backend) KillmailExists(ctx context.Context, id string) (bool, error) {
	_, span := otel.Tracer(packageName).Start(ctx, spanKillmailExists)
	defer span.End()

	span.SetAttributes(attribute.String("id", id))

	key := fmt.Sprintf("%s:%s", config.Get().Redict.Prefix, id)
	_, err := r.store.Get(context.Background(), key).Result()

	switch err {
	case nil:
		span.SetAttributes(attribute.String("cache", "hit"))
		slog.Debug("cache hit", "id", id)
		return true, nil
	case redis.Nil:
		span.SetAttributes(attribute.String("cache", "miss"))
		slog.Debug("cache miss", "id", id)
		return false, nil
	}

	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())

	return false, err
}

func (r *Backend) GetIgnoredSystemIDs(ctx context.Context) ([]string, error) {
	_, span := otel.Tracer(packageName).Start(ctx, spanGetIgnoredSystemIDs)
	defer span.End()

	key := fmt.Sprintf("%s:%s", config.Get().Redict.Prefix, keyIgnoredSystemIDs)
	ids, err := r.store.SMembers(context.Background(), key).Result()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "ok")
	return ids, nil
}

func (r *Backend) GetIgnoredRegionIDs(ctx context.Context) ([]string, error) {
	_, span := otel.Tracer(packageName).Start(ctx, spanGetIgnoredRegionIDs)
	defer span.End()

	key := fmt.Sprintf("%s:%s", config.Get().Redict.Prefix, keyIgnoredRegionIDs)
	ids, err := r.store.SMembers(context.Background(), key).Result()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "ok")
	return ids, nil
}

func (r *Backend) IgnoreSystemID(ctx context.Context, id int64) error {
	sctx, span := otel.Tracer(packageName).Start(ctx, spanIgnoreSystemID)
	defer span.End()

	span.SetAttributes(attribute.Int64("id", id))

	key := fmt.Sprintf("%s:%s", config.Get().Redict.Prefix, keyIgnoredSystemIDs)
	if _, err := r.store.SAdd(sctx, key, id).Result(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	span.SetStatus(codes.Ok, "ok")
	return nil
}

func (r *Backend) IgnoreRegionID(ctx context.Context, id int64) error {
	sctx, span := otel.Tracer(packageName).Start(ctx, spanIgnoreRegionID)
	defer span.End()

	span.SetAttributes(attribute.Int64("id", id))

	key := fmt.Sprintf("%s:%s", config.Get().Redict.Prefix, keyIgnoredRegionIDs)
	if _, err := r.store.SAdd(sctx, key, id).Result(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	span.SetStatus(codes.Ok, "ok")
	return nil
}
