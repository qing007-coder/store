package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"store/pkg/config"
	"time"
)

type Client struct {
	Client *redis.Client
}

func NewClient(conf *config.GlobalConfig) *Client {
	rdb := new(Client)

	rdb.init(conf)
	return rdb
}

func (r *Client) init(conf *config.GlobalConfig) {
	r.Client = redis.NewClient(r.GetConfig(conf))
}

func (r *Client) GetConfig(conf *config.GlobalConfig) *redis.Options {
	return &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", conf.Redis.Address, conf.Redis.Port),
		Password: conf.Redis.Password,
		DB:       conf.Redis.DB,
	}
}

func (r *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

func (r *Client) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

// SetNX key的值如果存在，则不做任何操作
func (r *Client) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return r.Client.SetNX(ctx, key, value, expiration).Result()
}

func (r *Client) Del(ctx context.Context, keys ...string) error {
	return r.Client.Del(ctx, keys...).Err()
}
