package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	auth_service "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/auth_service/v1"
)

// Регистрация сервиса аутентификации с использованием JWT-валидации для конкретных маршрутов
func RegisterAuthService(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return registerService(ctx, mux, endpoint, opts, func(conn *grpc.ClientConn) error {
		client := auth_service.NewAuthServiceClient(conn)

		handlers := []struct {
			method  string
			pattern string
			handler runtime.HandlerFunc
		}{
			{"POST", "/v1/auth/register", handleRegister(client)},
			{"POST", "/v1/auth/authenticate", handleAuthenticate(client)},
			{"POST", "/v1/auth/oauth", handleOAuthAuthenticate(client)},
			{"POST", "/v1/auth/verify-email", handleVerifyEmail(client)},
			{"POST", "/v1/auth/reset-password", withJWTValidation(handleResetPassword(client))},
			{"POST", "/v1/auth/change-password", withJWTValidation(handleChangePassword(client))},
		}

		for _, h := range handlers {
			if err := mux.HandlePath(h.method, h.pattern, h.handler); err != nil {
				return err
			}
		}

		return nil
	})
}

func handleRegister(client auth_service.AuthServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req auth_service.RegisterRequest
		if err := decodeJSONBody(w, r, &req); err != nil {
			return
		}

		ctx, cancel := withTimeout(r.Context(), 5*time.Second)
		defer cancel()

		respInterface, err := cb.Execute(func() (interface{}, error) {
			return client.Register(ctx, &req)
		})

		if err != nil {
			handleGrpcError(w, err)
			return
		}

		resp := respInterface.(*auth_service.RegisterResponse)
		writeJSONResponse(w, http.StatusOK, resp)
	}
}

func handleAuthenticate(client auth_service.AuthServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req auth_service.AuthenticateRequest
		if err := decodeJSONBody(w, r, &req); err != nil {
			return
		}

		ctx, cancel := withTimeout(r.Context(), 5*time.Second)
		defer cancel()

		respInterface, err := cb.Execute(func() (interface{}, error) {
			return client.Authenticate(ctx, &req)
		})

		if err != nil {
			handleGrpcError(w, err)
			return
		}

		resp := respInterface.(*auth_service.AuthenticateResponse)
		writeJSONResponse(w, http.StatusOK, resp)
	}
}

func handleOAuthAuthenticate(client auth_service.AuthServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req auth_service.OAuthAuthenticateRequest
		if err := decodeJSONBody(w, r, &req); err != nil {
			return
		}

		ctx, cancel := withTimeout(r.Context(), 5*time.Second)
		defer cancel()

		respInterface, err := cb.Execute(func() (interface{}, error) {
			return client.OAuthAuthenticate(ctx, &req)
		})

		if err != nil {
			handleGrpcError(w, err)
			return
		}

		resp := respInterface.(*auth_service.OAuthAuthenticateResponse)
		writeJSONResponse(w, http.StatusOK, resp)
	}
}

func handleVerifyEmail(client auth_service.AuthServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req auth_service.VerifyEmailRequest
		if err := decodeJSONBody(w, r, &req); err != nil {
			return
		}

		ctx, cancel := withTimeout(r.Context(), 5*time.Second)
		defer cancel()

		respInterface, err := cb.Execute(func() (interface{}, error) {
			return client.VerifyEmail(ctx, &req)
		})

		if err != nil {
			handleGrpcError(w, err)
			return
		}

		resp := respInterface.(*auth_service.VerifyEmailResponse)
		writeJSONResponse(w, http.StatusOK, resp)
	}
}

func handleResetPassword(client auth_service.AuthServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req auth_service.ResetPasswordRequest
		if err := decodeJSONBody(w, r, &req); err != nil {
			return
		}

		ctx, cancel := withTimeout(r.Context(), 5*time.Second)
		defer cancel()

		respInterface, err := cb.Execute(func() (interface{}, error) {
			return client.ResetPassword(ctx, &req)
		})

		if err != nil {
			handleGrpcError(w, err)
			return
		}

		resp := respInterface.(*auth_service.ResetPasswordResponse)
		writeJSONResponse(w, http.StatusOK, resp)
	}
}

func handleChangePassword(client auth_service.AuthServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req auth_service.ChangePasswordRequest
		if err := decodeJSONBody(w, r, &req); err != nil {
			return
		}

		ctx, cancel := withTimeout(r.Context(), 5*time.Second)
		defer cancel()

		respInterface, err := cb.Execute(func() (interface{}, error) {
			return client.ChangePassword(ctx, &req)
		})

		if err != nil {
			handleGrpcError(w, err)
			return
		}

		resp := respInterface.(*auth_service.ChangePasswordResponse)
		writeJSONResponse(w, http.StatusOK, resp)
	}
}
