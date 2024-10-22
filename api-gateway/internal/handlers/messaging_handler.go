package handlers

import (
	"context"
	"net/http"

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

	if err := mux.HandlePath("POST", "/v1/messaging/send-message", handleSendMessage(client)); err != nil {
		return err
	}
	if err := mux.HandlePath("GET", "/v1/messaging/messages", handleGetMessages(client)); err != nil {
		return err
	}
	if err := mux.HandlePath("POST", "/v1/messaging/update-message-status", handleUpdateMessageStatus(client)); err != nil {
		return err
	}

	return nil
}

func handleSendMessage(client messaging_service.MessagingServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req messaging_service.SendMessageRequest
		if err := decodeJSONBody(w, r, &req); err != nil {
			return
		}
		resp, err := client.SendMessage(r.Context(), &req)
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		writeJSONResponse(w, http.StatusCreated, resp)
	}
}

func handleGetMessages(client messaging_service.MessagingServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req messaging_service.GetMessagesRequest
		req.UserId = parseStringParam(r, "user_id", "")
		req.ConversationUserId = parseStringParam(r, "conversation_user_id", "")

		if req.UserId == "" || req.ConversationUserId == "" {
			http.Error(w, "Parameters 'user_id' and 'conversation_user_id' are required", http.StatusBadRequest)
			return
		}

		req.Limit = int32(parseIntParam(r, "limit", 0))
		req.Offset = int32(parseIntParam(r, "offset", 0))

		resp, err := client.GetMessages(r.Context(), &req)
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		writeJSONResponse(w, http.StatusOK, resp)
	}
}

func handleUpdateMessageStatus(client messaging_service.MessagingServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req messaging_service.UpdateMessageStatusRequest
		if err := decodeJSONBody(w, r, &req); err != nil {
			return
		}
		resp, err := client.UpdateMessageStatus(r.Context(), &req)
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		writeJSONResponse(w, http.StatusOK, resp)
	}
}
