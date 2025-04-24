package subscriber

import (
	"fmt"
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
				log.Println("Error unmarshalling message:", err)
				continue
			}

			// # TODO : auth_redis 에서 Method로 저장하게 수정
			key := fmt.Sprintf("auth_%s", authMessage.ID)
			err := client.Set(ctx, key, authMessage.Name, 0).Err()
			if err != nil {
				log.Println("Redis 저장 오류:", err)
			} else {
				log.Printf("[Auth-Redis] 저장됨: %s → %s\n", key, authMessage.Name)
			}}
	}()
}
