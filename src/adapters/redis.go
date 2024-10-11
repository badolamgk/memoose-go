package adapters

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCacheProvider struct {
	CacheProvider[string]
	client      *redis.Client
	connection  *redis.Conn
	connect     bool
	storesAsObj bool
	name        string
}

func (rcp *RedisCacheProvider) Name() string {
	return rcp.name
}

func (rcp *RedisCacheProvider) NewRedisCacheProvider(name string, options redis.Options, connect bool) {
	rcp.name = name
	rcp.client = redis.NewClient(&options)
	rcp.connect = connect
	if connect {
		rcp.connection = rcp.client.Conn()
	} else {
		rcp.connection = nil
	}
}

func (rcp *RedisCacheProvider) GetClientConnection() (*redis.Conn, error) {
	if rcp.connect {
		return rcp.connection, nil
	}
	return nil, errors.New("shouldnt call as client initialised with connect=false")
}

func (rcp *RedisCacheProvider) Get(ctx context.Context, key string) (string, error) {
	return rcp.client.Get(ctx, key).Result()
}

func (rcp *RedisCacheProvider) Set(ctx context.Context, key string, value string, ttl int64) (string, error) {
	var durationStr string
	if ttl > 0 {
		durationStr = fmt.Sprintf("%ds", ttl)
	} else {
		durationStr = "0s"
	}
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return "", err
	}
	return rcp.client.Set(ctx, key, value, duration).Result()
}

func (rcp *RedisCacheProvider) Del(ctx context.Context, keys ...string) (int64, error) {
	return rcp.client.Del(ctx, keys...).Result()
}

func (rcp *RedisCacheProvider) Pipeline() redis.Pipeliner {
	return rcp.client.Pipeline()
}

func (rcp *RedisCacheProvider) Expire(ctx context.Context, key string, ttl int64) (bool, error) {
	var durationStr string
	if ttl > 0 {
		durationStr = fmt.Sprintf("%ds", ttl)
	} else {
		durationStr = "0s"
	}
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return false, err
	}
	return rcp.client.Expire(ctx, key, duration).Result()
}

func (rcp *RedisCacheProvider) MGet(ctx context.Context, keys ...string) ([]any, error) {
	return rcp.client.MGet(ctx, keys...).Result()
}

func (rcp *RedisCacheProvider) MSet(ctx context.Context, keyValues ...interface{}) (string, error) {
	return rcp.client.MSet(ctx, keyValues...).Result()
}

func (rcp *RedisCacheProvider) Exists(ctx context.Context, keys ...string) (int64, error) {
	return rcp.client.Exists(ctx, keys...).Result()
}

func (rcp *RedisCacheProvider) LPush(ctx context.Context, key string, values ...any) (int64, error) {
	return rcp.client.LPush(ctx, key, values...).Result()
}

func (rcp *RedisCacheProvider) LLen(ctx context.Context, key string) (int64, error) {
	return rcp.client.LLen(ctx, key).Result()
}

func (rcp *RedisCacheProvider) SCard(ctx context.Context, key string) (int64, error) {
	return rcp.client.SCard(ctx, key).Result()
}

func (rcp *RedisCacheProvider) SIsMember(ctx context.Context, key string, member any) (bool, error) {
	return rcp.client.SIsMember(ctx, key, member).Result()
}

func (rcp *RedisCacheProvider) SAdd(ctx context.Context, key string, members ...any) (int64, error) {
	return rcp.client.SAdd(ctx, key, members...).Result()
}

func (rcp *RedisCacheProvider) SRem(ctx context.Context, key string, members ...any) (int64, error) {
	return rcp.client.SRem(ctx, key, members...).Result()
}

func (rcp *RedisCacheProvider) SMembers(ctx context.Context, key string) ([]string, error) {
	return rcp.client.SMembers(ctx, key).Result()
}

func (rcp *RedisCacheProvider) RPush(ctx context.Context, key string, values ...any) (int64, error) {
	return rcp.client.RPush(ctx, key, values...).Result()
}

func (rcp *RedisCacheProvider) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return rcp.client.LRange(ctx, key, start, stop).Result()
}

func (rcp *RedisCacheProvider) ping(ctx context.Context) (string, error) {
	return rcp.client.Ping(ctx).Result()
}

func (rcp *RedisCacheProvider) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return rcp.client.Subscribe(ctx, channels...)
}

func (rcp *RedisCacheProvider) Publish(ctx context.Context, channel string, message any) (int64, error) {
	return rcp.client.Publish(ctx, channel, message).Result()
}

func (rcp *RedisCacheProvider) FlushDB(ctx context.Context) (string, error) {
	return rcp.client.FlushDB(ctx).Result()
}

func (rcp *RedisCacheProvider) LPop(ctx context.Context, key string) (string, error) {
	return rcp.client.LPop(ctx, key).Result()
}

func (rcp *RedisCacheProvider) decrby(ctx context.Context, key string, count int64) (int64, error) {
	return rcp.client.DecrBy(ctx, key, count).Result()
}
