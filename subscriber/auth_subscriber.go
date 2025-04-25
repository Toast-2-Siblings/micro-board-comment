package subscriber

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Toast-2-Siblings/micro-board-comment/redis"
)

// 장고 서버로부터 발생된 auth_user_created 이벤트를 수신하기 위한 구조체
type AuthMessage struct {
	ID string `json:"user_id"`
	Name string `json:"user_name"`
}

// 장고 서버로부터 auth_user_created 이벤트를 구독하는 함수
func SubscribeAuthUserCreated(ctx context.Context) {
	auth_redis, err := redis.GetAuthRedis(ctx)
	if err != nil {
		SubscribeAuthUserCreated(ctx)
		log.Println("Error getting auth redis client:", err)
		return
	}

	client := auth_redis.GetClient()
	pubsub := client.Subscribe(ctx, "user_created")
	ch := pubsub.Channel()

	go func() {
		for msg := range ch {
			var authMessage AuthMessage
			if err := json.Unmarshal([]byte(msg.Payload), &authMessage); err != nil {
				log.Println("[Auth-Redis] Error unmarshalling message:", err)
				continue
			}

			if err :=  auth_redis.SetAuth(ctx, authMessage.ID, authMessage.Name); err != nil {
				log.Println("[Auth-Redis] Error setting auth in redis:", err)
				continue
			}
			log.Printf("[Auth-Redis] User created: %s, Name: %s\n", authMessage.ID, authMessage.Name)
		}
	}()
}
