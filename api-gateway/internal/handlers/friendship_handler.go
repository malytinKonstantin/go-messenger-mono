package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	friendship_service "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/friendship_service/v1"
)

func RegisterFriendshipService(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := friendship_service.NewFriendshipServiceClient(conn)

	err = mux.HandlePath("POST", "/v1/friendship/send-request", handleSendFriendRequest(client))
	if err != nil {
		return err
	}

	err = mux.HandlePath("POST", "/v1/friendship/accept-request", handleAcceptFriendRequest(client))
	if err != nil {
		return err
	}

	err = mux.HandlePath("POST", "/v1/friendship/reject-request", handleRejectFriendRequest(client))
	if err != nil {
		return err
	}

	err = mux.HandlePath("POST", "/v1/friendship/remove-friend", handleRemoveFriend(client))
	if err != nil {
		return err
	}

	err = mux.HandlePath("GET", "/v1/friendship/friends-list", handleGetFriendsList(client))
	if err != nil {
		return err
	}

	err = mux.HandlePath("GET", "/v1/friendship/pending-requests", handleGetPendingRequests(client))
	if err != nil {
		return err
	}

	return nil
}

func handleSendFriendRequest(client friendship_service.FriendshipServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req friendship_service.SendFriendRequestRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp, err := client.SendFriendRequest(r.Context(), &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func handleAcceptFriendRequest(client friendship_service.FriendshipServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req friendship_service.AcceptFriendRequestRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp, err := client.AcceptFriendRequest(r.Context(), &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func handleRejectFriendRequest(client friendship_service.FriendshipServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req friendship_service.RejectFriendRequestRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp, err := client.RejectFriendRequest(r.Context(), &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func handleRemoveFriend(client friendship_service.FriendshipServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req friendship_service.RemoveFriendRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp, err := client.RemoveFriend(r.Context(), &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func handleGetFriendsList(client friendship_service.FriendshipServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		userID := r.URL.Query().Get("user_id")
		req := &friendship_service.GetFriendsListRequest{UserId: userID}
		resp, err := client.GetFriendsList(r.Context(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func handleGetPendingRequests(client friendship_service.FriendshipServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		userID := r.URL.Query().Get("user_id")
		req := &friendship_service.GetPendingRequestsRequest{UserId: userID}
		resp, err := client.GetPendingRequests(r.Context(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
