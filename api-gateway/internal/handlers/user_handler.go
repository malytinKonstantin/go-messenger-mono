package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	user_service "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/user_service/v1"
)

func RegisterUserService(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return registerService(ctx, mux, endpoint, opts, func(conn *grpc.ClientConn) error {
		client := user_service.NewUserServiceClient(conn)

		handlers := []struct {
			method  string
			pattern string
			handler runtime.HandlerFunc
		}{
			{"GET", "/v1/users/{user_id}", withJWTValidation(handleGetUser(client))},
			{"POST", "/v1/users", withJWTValidation(handleCreateUserProfile(client))},
			{"PUT", "/v1/users/{user_id}", withJWTValidation(handleUpdateUserProfile(client))},
			{"GET", "/v1/users/search", withJWTValidation(handleSearchUsers(client))},
		}

		for _, h := range handlers {
			if err := mux.HandlePath(h.method, h.pattern, h.handler); err != nil {
				return err
			}
		}

		return nil
	})
}

func handleGetUser(client user_service.UserServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		userID, ok := pathParams["user_id"]
		if !ok {
			http.Error(w, "user_id is not specified", http.StatusBadRequest)
			return
		}
		req := &user_service.GetUserRequest{UserId: userID}

		ctx, cancel := withTimeout(r.Context(), 5*time.Second)
		defer cancel()

		respInterface, err := cb.Execute(func() (interface{}, error) {
			return client.GetUser(ctx, req)
		})
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		resp := respInterface.(*user_service.GetUserResponse)
		writeJSONResponse(w, http.StatusOK, resp)
	}
}

func handleCreateUserProfile(client user_service.UserServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req user_service.CreateUserProfileRequest
		if err := decodeJSONBody(w, r, &req); err != nil {
			return
		}

		ctx, cancel := withTimeout(r.Context(), 5*time.Second)
		defer cancel()

		respInterface, err := cb.Execute(func() (interface{}, error) {
			return client.CreateUserProfile(ctx, &req)
		})
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		resp := respInterface.(*user_service.CreateUserProfileResponse)
		writeJSONResponse(w, http.StatusCreated, resp)
	}
}

func handleUpdateUserProfile(client user_service.UserServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		userID, ok := pathParams["user_id"]
		if !ok {
			http.Error(w, "user_id is not specified", http.StatusBadRequest)
			return
		}
		var req user_service.UpdateUserProfileRequest
		if err := decodeJSONBody(w, r, &req); err != nil {
			return
		}
		req.UserId = userID

		ctx, cancel := withTimeout(r.Context(), 5*time.Second)
		defer cancel()

		respInterface, err := cb.Execute(func() (interface{}, error) {
			return client.UpdateUserProfile(ctx, &req)
		})
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		resp := respInterface.(*user_service.UpdateUserProfileResponse)
		writeJSONResponse(w, http.StatusOK, resp)
	}
}

func handleSearchUsers(client user_service.UserServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		query := parseStringParam(r, "query", "")
		if query == "" {
			http.Error(w, "Query parameter is required", http.StatusBadRequest)
			return
		}
		limit := int32(parseIntParam(r, "limit", 10))
		offset := int32(parseIntParam(r, "offset", 0))

		req := &user_service.SearchUsersRequest{
			Query:  query,
			Limit:  limit,
			Offset: offset,
		}

		ctx, cancel := withTimeout(r.Context(), 5*time.Second)
		defer cancel()

		respInterface, err := cb.Execute(func() (interface{}, error) {
			return client.SearchUsers(ctx, req)
		})
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		resp := respInterface.(*user_service.SearchUsersResponse)
		writeJSONResponse(w, http.StatusOK, resp)
	}
}
