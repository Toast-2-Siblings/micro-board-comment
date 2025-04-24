package redis

import "context"

func InitialRedis(ctx context.Context) error {
	if err := NewAuthRedis(ctx); err != nil {
		return err
	}

	return nil
}
