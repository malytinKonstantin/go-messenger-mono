package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	user_service "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/user_service/v1"
)

func RegisterUserService(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return err
	}

	client := user_service.NewUserServiceClient(conn)

	if err := mux.HandlePath("GET", "/v1/users/{user_id}", handleGetUser(client)); err != nil {
		return err
	}

	if err := mux.HandlePath("POST", "/v1/users", handleCreateUserProfile(client)); err != nil {
		return err
	}

	if err := mux.HandlePath("PUT", "/v1/users/{user_id}", handleUpdateUserProfile(client)); err != nil {
		return err
	}

	if err := mux.HandlePath("GET", "/v1/users/search", handleSearchUsers(client)); err != nil {
		return err
	}

	return nil
}

func handleGetUser(client user_service.UserServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		userID, ok := pathParams["user_id"]
		if !ok {
			http.Error(w, "user_id not provided", http.StatusBadRequest)
			return
		}
		req := &user_service.GetUserRequest{UserId: userID}
		resp, err := client.GetUser(r.Context(), req)
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func handleCreateUserProfile(client user_service.UserServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req user_service.CreateUserProfileRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp, err := client.CreateUserProfile(r.Context(), &req)
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func handleUpdateUserProfile(client user_service.UserServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		userID, ok := pathParams["user_id"]
		if !ok {
			http.Error(w, "user_id not provided", http.StatusBadRequest)
			return
		}
		var req user_service.UpdateUserProfileRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		req.UserId = userID
		resp, err := client.UpdateUserProfile(r.Context(), &req)
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func handleSearchUsers(client user_service.UserServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		query := r.URL.Query().Get("query")
		if query == "" {
			http.Error(w, "query parameter is required", http.StatusBadRequest)
			return
		}

		limit := int32(parseIntParam(r, "limit", 10))
		offset := int32(parseIntParam(r, "offset", 0))

		req := &user_service.SearchUsersRequest{
			Query:  query,
			Limit:  limit,
			Offset: offset,
		}
		resp, err := client.SearchUsers(r.Context(), req)
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
