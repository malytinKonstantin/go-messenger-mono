package handlers

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	auth_service "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/auth_service/v1"
)

func RegisterAuthService(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return err
	}

	client := auth_service.NewAuthServiceClient(conn)

	if err := mux.HandlePath("POST", "/v1/auth/register", handleRegister(client)); err != nil {
		return err
	}
	if err := mux.HandlePath("POST", "/v1/auth/authenticate", handleAuthenticate(client)); err != nil {
		return err
	}
	if err := mux.HandlePath("POST", "/v1/auth/oauth", handleOAuthAuthenticate(client)); err != nil {
		return err
	}
	if err := mux.HandlePath("POST", "/v1/auth/verify-email", handleVerifyEmail(client)); err != nil {
		return err
	}
	if err := mux.HandlePath("POST", "/v1/auth/reset-password", handleResetPassword(client)); err != nil {
		return err
	}
	if err := mux.HandlePath("POST", "/v1/auth/change-password", handleChangePassword(client)); err != nil {
		return err
	}

	return nil
}

func handleRegister(client auth_service.AuthServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req auth_service.RegisterRequest
		if err := decodeJSONBody(w, r, &req); err != nil {
			return
		}
		resp, err := client.Register(r.Context(), &req)
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		writeJSONResponse(w, http.StatusOK, resp)
	}
}

func handleAuthenticate(client auth_service.AuthServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req auth_service.AuthenticateRequest
		if err := decodeJSONBody(w, r, &req); err != nil {
			return
		}
		resp, err := client.Authenticate(r.Context(), &req)
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		writeJSONResponse(w, http.StatusOK, resp)
	}
}

func handleOAuthAuthenticate(client auth_service.AuthServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req auth_service.OAuthAuthenticateRequest
		if err := decodeJSONBody(w, r, &req); err != nil {
			return
		}
		resp, err := client.OAuthAuthenticate(r.Context(), &req)
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		writeJSONResponse(w, http.StatusOK, resp)
	}
}

func handleVerifyEmail(client auth_service.AuthServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req auth_service.VerifyEmailRequest
		if err := decodeJSONBody(w, r, &req); err != nil {
			return
		}
		resp, err := client.VerifyEmail(r.Context(), &req)
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		writeJSONResponse(w, http.StatusOK, resp)
	}
}

func handleResetPassword(client auth_service.AuthServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req auth_service.ResetPasswordRequest
		if err := decodeJSONBody(w, r, &req); err != nil {
			return
		}
		resp, err := client.ResetPassword(r.Context(), &req)
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		writeJSONResponse(w, http.StatusOK, resp)
	}
}

func handleChangePassword(client auth_service.AuthServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req auth_service.ChangePasswordRequest
		if err := decodeJSONBody(w, r, &req); err != nil {
			return
		}
		resp, err := client.ChangePassword(r.Context(), &req)
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		writeJSONResponse(w, http.StatusOK, resp)
	}
}
