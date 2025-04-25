package redis

import (
	"fmt"
	"context"
	"errors"
	"sync"

	"github.com/Toast-2-Siblings/micro-board-comment/config"
	"github.com/Toast-2-Siblings/micro-board-comment/utils/convert"
	"github.com/go-redis/redis/v8"
)

type AuthRedis interface {
	GetClient() *redis.Client
	Close()
	SetAuth(ctx context.Context, key string, value string) error
}

type authRedis struct {
	client *redis.Client
}

var (
	auth_instance *authRedis
	auth_once sync.Once
)

func NewAuthRedis(ctx context.Context) error {
	cfg := config.GetConfig()
	db := convert.InterfaceToInt(cfg.Redis.RedisAuthDB)

	new_redis := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.RedisHost + ":" + cfg.Redis.RedisPort,
		Password: cfg.Redis.RedisPass,
		DB: db,
	})

	_, err := new_redis.Ping(ctx).Result()
	if err != nil {
		return err
	}

	auth_instance = &authRedis{
		client: new_redis,
	}

	return nil
}

func GetAuthRedis(ctx context.Context) (AuthRedis, error) {
	var err error

	auth_once.Do(func() {
		err = NewAuthRedis(ctx)
	})

	if err != nil {
		return nil, err
	}

	if auth_instance == nil {
		return nil, errors.New("auth instance is nil")
	}

	return auth_instance, nil
}

func (a *authRedis) GetClient() *redis.Client {
	return a.client
}

func (a *authRedis) Close() {
	if a.client != nil {
		a.client.Close()
	}

	auth_instance = nil
}

func (a *authRedis) SetAuth(ctx context.Context, key string, value string) error {
	auth_key := fmt.Sprintf("auth_%s", key)
	if err := a.client.Set(ctx, auth_key, value, 0).Err(); err != nil {
		return err
	}

	return nil
}
