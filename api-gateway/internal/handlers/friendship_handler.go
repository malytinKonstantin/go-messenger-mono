package handlers

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	friendship_service "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/friendship_service/v1"
)

func RegisterFriendshipService(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return registerService(ctx, mux, endpoint, opts, func(conn *grpc.ClientConn) error {
		client := friendship_service.NewFriendshipServiceClient(conn)

		handlers := []struct {
			method  string
			pattern string
			handler runtime.HandlerFunc
		}{
			{"POST", "/v1/friendship/send-request", handleSendFriendRequest(client)},
			{"POST", "/v1/friendship/accept-request", handleAcceptFriendRequest(client)},
			{"POST", "/v1/friendship/reject-request", handleRejectFriendRequest(client)},
			{"POST", "/v1/friendship/remove-friend", handleRemoveFriend(client)},
			{"GET", "/v1/friendship/friends-list", handleGetFriendsList(client)},
			{"GET", "/v1/friendship/pending-requests", handleGetPendingRequests(client)},
		}

		for _, h := range handlers {
			if err := mux.HandlePath(h.method, h.pattern, h.handler); err != nil {
				return err
			}
		}

		return nil
	})
}

func handleSendFriendRequest(client friendship_service.FriendshipServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req friendship_service.SendFriendRequestRequest
		if err := decodeJSONBody(w, r, &req); err != nil {
			return
		}
		resp, err := client.SendFriendRequest(r.Context(), &req)
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		writeJSONResponse(w, http.StatusCreated, resp)
	}
}

func handleAcceptFriendRequest(client friendship_service.FriendshipServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req friendship_service.AcceptFriendRequestRequest
		if err := decodeJSONBody(w, r, &req); err != nil {
			return
		}
		resp, err := client.AcceptFriendRequest(r.Context(), &req)
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		writeJSONResponse(w, http.StatusOK, resp)
	}
}

func handleRejectFriendRequest(client friendship_service.FriendshipServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req friendship_service.RejectFriendRequestRequest
		if err := decodeJSONBody(w, r, &req); err != nil {
			return
		}
		resp, err := client.RejectFriendRequest(r.Context(), &req)
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		writeJSONResponse(w, http.StatusOK, resp)
	}
}

func handleRemoveFriend(client friendship_service.FriendshipServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req friendship_service.RemoveFriendRequest
		if err := decodeJSONBody(w, r, &req); err != nil {
			return
		}
		resp, err := client.RemoveFriend(r.Context(), &req)
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		writeJSONResponse(w, http.StatusOK, resp)
	}
}

func handleGetFriendsList(client friendship_service.FriendshipServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		userID := parseStringParam(r, "user_id", "")
		if userID == "" {
			http.Error(w, "The user_id parameter is required", http.StatusBadRequest)
			return
		}
		req := &friendship_service.GetFriendsListRequest{UserId: userID}
		resp, err := client.GetFriendsList(r.Context(), req)
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		writeJSONResponse(w, http.StatusOK, resp)
	}
}

func handleGetPendingRequests(client friendship_service.FriendshipServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		userID := parseStringParam(r, "user_id", "")
		if userID == "" {
			http.Error(w, "The user_id parameter is required", http.StatusBadRequest)
			return
		}
		req := &friendship_service.GetPendingRequestsRequest{UserId: userID}
		resp, err := client.GetPendingRequests(r.Context(), req)
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		writeJSONResponse(w, http.StatusOK, resp)
	}
}
