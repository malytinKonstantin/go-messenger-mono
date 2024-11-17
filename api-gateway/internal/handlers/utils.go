package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/malytinKonstantin/go-messenger-mono/shared/circuitbreaker"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var cb *circuitbreaker.CircuitBreaker

func init() {
	cb = circuitbreaker.NewCircuitBreaker("APIGatewayCircuitBreaker")
}

func registerService(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption, registerFunc func(clientConn *grpc.ClientConn) error) error {
	dialCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(dialCtx, endpoint, append(opts, grpc.WithBlock())...)
	if err != nil {
		return fmt.Errorf("failed to establish connection: %w", err)
	}

	if err := registerFunc(conn); err != nil {
		return fmt.Errorf("error registering service: %w", err)
	}

	return nil
}

// parseIntParam извлекает целочисленный параметр из URL-запроса.
func parseIntParam(r *http.Request, name string, defaultValue int) int {
	value := r.URL.Query().Get(name)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

// parseStringParam извлекает строковый параметр из URL-запроса.
func parseStringParam(r *http.Request, name string, defaultValue string) string {
	value := r.URL.Query().Get(name)
	if value == "" {
		return defaultValue
	}
	return value
}

// handleGrpcError обрабатывает ошибки gRPC и возвращает соответствующий HTTP-статус.
func handleGrpcError(w http.ResponseWriter, err error) {
	grpcStatus, ok := status.FromError(err)
	if ok {
		httpCode := runtime.HTTPStatusFromCode(grpcStatus.Code())
		http.Error(w, grpcStatus.Message(), httpCode)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// decodeJSONBody декодирует JSON из тела запроса.
func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(dst); err != nil {
		http.Error(w, "Invalid JSON format: "+err.Error(), http.StatusBadRequest)
		return err
	}
	return nil
}

// writeJSONResponse записывает ответ в формате JSON.
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// Функция-обёртка для валидации JWT
func withJWTValidation(handler runtime.HandlerFunc) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			http.Error(w, "invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		secretKey := []byte(viper.GetString("JWT_SECRET"))

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid token signing method")
			}
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)
		md := metadata.Pairs("user", claims["user"].(string))
		ctx := metadata.NewIncomingContext(r.Context(), md)

		handler(w, r.WithContext(ctx), pathParams)
	}
}

// extractAndForwardAuthHeader извлекает заголовок Authorization из HTTP-запроса и добавляет его в контекст.
func extractAndForwardAuthHeader(ctx context.Context, r *http.Request) context.Context {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		md := metadata.Pairs("authorization", authHeader)
		ctx = metadata.NewOutgoingContext(ctx, md)
	}
	return ctx
}

func withTimeout(parentCtx context.Context, duration time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parentCtx, duration)
}

func extractAndForwardHeaders(ctx context.Context, r *http.Request) context.Context {
	md := metadata.MD{}

	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		md.Append("authorization", authHeader)
	}

	if idempotencyKey := r.Header.Get("Idempotency-Key"); idempotencyKey != "" {
		md.Append("idempotency-key", idempotencyKey)
	}

	return metadata.NewOutgoingContext(ctx, md)
}
