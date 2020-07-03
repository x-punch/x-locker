package redlock

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// RedisClient represents client for Redis server.
type RedisClient interface {
	Do(ctx context.Context, args ...interface{}) *redis.Cmd

	Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd
	EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd
	ScriptExists(ctx context.Context, hashes ...string) *redis.BoolSliceCmd
	ScriptLoad(ctx context.Context, script string) *redis.StringCmd
}
