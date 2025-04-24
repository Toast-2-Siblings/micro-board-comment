package redis

import (
	"context"
	"errors"
	"sync"

	"github.com/Toast-2-Siblings/micro-board-comment/config"
	"github.com/Toast-2-Siblings/micro-board-comment/utils/convert"
	"github.com/go-redis/redis/v8"
)

type AuthRedis interface {
	GetClient() *redis.Client
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


