package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	auth_service "github.com/malytinKonstantin/go-messenger-mono/proto/auth-service"
	friendship_service "github.com/malytinKonstantin/go-messenger-mono/proto/friendship-service"
	messaging_service "github.com/malytinKonstantin/go-messenger-mono/proto/messaging-service"
	notification_service "github.com/malytinKonstantin/go-messenger-mono/proto/notification-service"
	user_service "github.com/malytinKonstantin/go-messenger-mono/proto/user-service"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	grpcMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Регистрируем сервисы
	if err := registerAuthService(ctx, grpcMux, "auth-service:50051", opts); err != nil {
		log.Fatalf("Failed to register auth service: %v", err)
	}
	if err := registerFriendshipService(ctx, grpcMux, "friendship-service:50052", opts); err != nil {
		log.Fatalf("Failed to register friendship service: %v", err)
	}
	if err := registerMessagingService(ctx, grpcMux, "messaging-service:50053", opts); err != nil {
		log.Fatalf("Failed to register messaging service: %v", err)
	}
	if err := registerNotificationService(ctx, grpcMux, "notification-service:50054", opts); err != nil {
		log.Fatalf("Failed to register notification service: %v", err)
	}
	if err := registerUserService(ctx, grpcMux, "user-service:50055", opts); err != nil {
		log.Fatalf("Failed to register user service: %v", err)
	}

	app := fiber.New()
	app.All("/*", adaptor.HTTPHandler(grpcMux))

	log.Println("Starting API Gateway on :3000")
	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func registerAuthService(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return err
	}
	client := auth_service.NewAuthServiceClient(conn)
	return mux.HandlePath("POST", "/v1/auth/authenticate", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req auth_service.AuthRequest
		body, _ := ioutil.ReadAll(r.Body)
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp, err := client.Authenticate(r.Context(), &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
}

func registerFriendshipService(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return err
	}
	client := friendship_service.NewFriendshipServiceClient(conn)

	err = mux.HandlePath("POST", "/v1/friendship/send-request", handleGrpcRequest(client.SendFriendRequest))
	if err != nil {
		return err
	}

	err = mux.HandlePath("POST", "/v1/friendship/accept-request", handleGrpcRequest(client.AcceptFriendRequest))
	if err != nil {
		return err
	}

	return nil
}

func registerMessagingService(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return err
	}
	client := messaging_service.NewMessagingServiceClient(conn)

	err = mux.HandlePath("POST", "/v1/messaging/send", handleGrpcRequest(client.SendMessage))
	if err != nil {
		return err
	}

	err = mux.HandlePath("GET", "/v1/messaging/get", handleGrpcRequest(client.GetMessages))
	if err != nil {
		return err
	}

	return nil
}

func registerNotificationService(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return err
	}
	client := notification_service.NewNotificationServiceClient(conn)

	err = mux.HandlePath("POST", "/v1/notification/send", handleGrpcRequest(client.SendNotification))
	if err != nil {
		return err
	}

	err = mux.HandlePath("GET", "/v1/notification/get", handleGrpcRequest(client.GetNotifications))
	if err != nil {
		return err
	}

	return nil
}

func registerUserService(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return err
	}
	client := user_service.NewUserServiceClient(conn)

	err = mux.HandlePath("POST", "/v1/user/create", handleGrpcRequest(client.CreateUser))
	if err != nil {
		return err
	}

	err = mux.HandlePath("GET", "/v1/user/get", handleGrpcRequest(client.GetUser))
	if err != nil {
		return err
	}

	return nil
}

func handleGrpcRequest(grpcFunc interface{}) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		ctx := r.Context()

		// Создаем пустой запрос нужного типа
		reqValue := reflect.New(reflect.TypeOf(grpcFunc).In(1).Elem()).Interface()

		// Парсим тело запроса
		body, _ := ioutil.ReadAll(r.Body)
		if err := json.Unmarshal(body, reqValue); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Вызываем gRPC функцию
		results := reflect.ValueOf(grpcFunc).Call([]reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(reqValue),
		})

		// Проверяем ошибку
		if err := results[1].Interface(); err != nil {
			http.Error(w, err.(error).Error(), http.StatusInternalServerError)
			return
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results[0].Interface())
	}
}
