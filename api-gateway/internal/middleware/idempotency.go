package middleware

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/malytinKonstantin/go-messenger-mono/shared/cache"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

type IdempotencyInterceptor struct {
	redisClient *redis.Client
}

func NewIdempotencyInterceptor() *IdempotencyInterceptor {
	return &IdempotencyInterceptor{
		redisClient: cache.GetRedisClient(),
	}
}

func (i *IdempotencyInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}

		idempotencyKeys := md.Get("idempotency-key")
		if len(idempotencyKeys) == 0 {
			return handler(ctx, req)
		}

		idempotencyKey := idempotencyKeys[0]
		hash := sha256.Sum256([]byte(idempotencyKey))
		key := fmt.Sprintf("idempotency:%x", hash)

		// Попытка получить сохраненный результат из Redis
		val, err := i.redisClient.Get(ctx, key).Bytes()
		if err == redis.Nil {
			// Ключ не найден, выполняем запрос и сохраняем результат
			resp, err := handler(ctx, req)
			if err != nil {
				return nil, err
			}

			// Сериализуем ответ
			respData, err := proto.Marshal(resp.(proto.Message))
			if err != nil {
				return nil, err
			}

			// Сохраняем ответ в Redis с временем жизни
			err = i.redisClient.Set(ctx, key, respData, time.Hour).Err()
			if err != nil {
				return nil, err
			}

			return resp, nil
		} else if err != nil {
			// Ошибка при обращении к Redis
			return nil, err
		} else {
			// Ключ найден, возвращаем сохраненный результат
			resp := proto.Clone(req.(proto.Message))
			err = proto.Unmarshal(val, resp.(proto.Message))
			if err != nil {
				return nil, err
			}

			return resp, nil
		}
	}
}
