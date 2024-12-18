package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	messaging_service "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/messaging_service/v1"
)

func RegisterMessagingService(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return registerService(ctx, mux, endpoint, opts, func(conn *grpc.ClientConn) error {
		client := messaging_service.NewMessagingServiceClient(conn)

		handlers := []struct {
			method  string
			pattern string
			handler runtime.HandlerFunc
		}{
			{"POST", "/v1/messaging/send-message", withJWTValidation(handleSendMessage(client))},
			{"GET", "/v1/messaging/messages", withJWTValidation(handleGetMessages(client))},
			{"POST", "/v1/messaging/update-message-status", withJWTValidation(handleUpdateMessageStatus(client))},
		}

		for _, h := range handlers {
			if err := mux.HandlePath(h.method, h.pattern, h.handler); err != nil {
				return err
			}
		}

		return nil
	})
}

func handleSendMessage(client messaging_service.MessagingServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req messaging_service.SendMessageRequest
		if err := decodeJSONBody(w, r, &req); err != nil {
			return
		}

		ctx, cancel := withTimeout(r.Context(), 5*time.Second)
		defer cancel()

		respInterface, err := cb.Execute(func() (interface{}, error) {
			return client.SendMessage(ctx, &req)
		})
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		resp := respInterface.(*messaging_service.SendMessageResponse)
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

		ctx, cancel := withTimeout(r.Context(), 5*time.Second)
		defer cancel()

		respInterface, err := cb.Execute(func() (interface{}, error) {
			return client.GetMessages(ctx, &req)
		})
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		resp := respInterface.(*messaging_service.GetMessagesResponse)
		writeJSONResponse(w, http.StatusOK, resp)
	}
}

func handleUpdateMessageStatus(client messaging_service.MessagingServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req messaging_service.UpdateMessageStatusRequest
		if err := decodeJSONBody(w, r, &req); err != nil {
			return
		}

		ctx, cancel := withTimeout(r.Context(), 5*time.Second)
		defer cancel()

		respInterface, err := cb.Execute(func() (interface{}, error) {
			return client.UpdateMessageStatus(ctx, &req)
		})
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		resp := respInterface.(*messaging_service.UpdateMessageStatusResponse)
		writeJSONResponse(w, http.StatusOK, resp)
	}
}
