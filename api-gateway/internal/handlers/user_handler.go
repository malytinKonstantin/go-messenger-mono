package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	user_service "github.com/malytinKonstantin/go-messenger-mono/proto/user-service"
)

func RegisterUserService(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := user_service.NewUserServiceClient(conn)

	err = mux.HandlePath("GET", "/v1/user/get", handleGetUser(client))
	if err != nil {
		return err
	}

	err = mux.HandlePath("POST", "/v1/user/create", handleCreateUserProfile(client))
	if err != nil {
		return err
	}

	err = mux.HandlePath("PUT", "/v1/user/update", handleUpdateUserProfile(client))
	if err != nil {
		return err
	}

	err = mux.HandlePath("GET", "/v1/user/search", handleSearchUsers(client))
	if err != nil {
		return err
	}

	return nil
}

func handleGetUser(client user_service.UserServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		userID := r.URL.Query().Get("user_id")
		req := &user_service.GetUserRequest{UserId: userID}
		resp, err := client.GetUser(r.Context(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func handleUpdateUserProfile(client user_service.UserServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req user_service.UpdateUserProfileRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp, err := client.UpdateUserProfile(r.Context(), &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func handleSearchUsers(client user_service.UserServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		query := r.URL.Query().Get("query")
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		req := &user_service.SearchUsersRequest{
			Query:  query,
			Limit:  int32(limit),
			Offset: int32(offset),
		}
		resp, err := client.SearchUsers(r.Context(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
