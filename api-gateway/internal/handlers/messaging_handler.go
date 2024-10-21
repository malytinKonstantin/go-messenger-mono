package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	messaging_service "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/messaging_service/v1"
)

func RegisterMessagingService(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return err
	}

	client := messaging_service.NewMessagingServiceClient(conn)

	err = mux.HandlePath("POST", "/v1/messaging/send-message", handleSendMessage(client))
	if err != nil {
		return err
	}

	err = mux.HandlePath("GET", "/v1/messaging/messages", handleGetMessages(client))
	if err != nil {
		return err
	}

	err = mux.HandlePath("POST", "/v1/messaging/update-message-status", handleUpdateMessageStatus(client))
	if err != nil {
		return err
	}

	return nil
}

func handleSendMessage(client messaging_service.MessagingServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req messaging_service.SendMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request format: "+err.Error(), http.StatusBadRequest)
			return
		}
		resp, err := client.SendMessage(r.Context(), &req)
		if err != nil {
			http.Error(w, "Error sending message: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func handleGetMessages(client messaging_service.MessagingServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req messaging_service.GetMessagesRequest

		queryParams := r.URL.Query()
		req.UserId = queryParams.Get("user_id")
		req.ConversationUserId = queryParams.Get("conversation_user_id")

		if req.UserId == "" || req.ConversationUserId == "" {
			http.Error(w, "Parameters 'user_id' and 'conversation_user_id' are required", http.StatusBadRequest)
			return
		}

		if limitStr := queryParams.Get("limit"); limitStr != "" {
			limit, err := strconv.Atoi(limitStr)
			if err != nil || limit <= 0 {
				http.Error(w, "Invalid 'limit' parameter", http.StatusBadRequest)
				return
			}
			req.Limit = int32(limit)
		}

		if offsetStr := queryParams.Get("offset"); offsetStr != "" {
			offset, err := strconv.Atoi(offsetStr)
			if err != nil || offset < 0 {
				http.Error(w, "Invalid 'offset' parameter", http.StatusBadRequest)
				return
			}
			req.Offset = int32(offset)
		}

		resp, err := client.GetMessages(r.Context(), &req)
		if err != nil {
			http.Error(w, "Error getting messages: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func handleUpdateMessageStatus(client messaging_service.MessagingServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req messaging_service.UpdateMessageStatusRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request format: "+err.Error(), http.StatusBadRequest)
			return
		}
		resp, err := client.UpdateMessageStatus(r.Context(), &req)
		if err != nil {
			http.Error(w, "Error updating message status: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
