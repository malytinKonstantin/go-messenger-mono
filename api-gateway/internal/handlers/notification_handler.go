package handlers

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	notification_service "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/notification_service/v1"
)

func RegisterNotificationService(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return registerService(ctx, mux, endpoint, opts, func(conn *grpc.ClientConn) error {
		client := notification_service.NewNotificationServiceClient(conn)

		handlers := []struct {
			method  string
			pattern string
			handler runtime.HandlerFunc
		}{
			{"POST", "/v1/notification/send", handleSendNotification(client)},
			{"GET", "/v1/notification", handleGetNotifications(client)},
			{"PUT", "/v1/notification/mark-read", handleMarkNotificationAsRead(client)},
			{"PUT", "/v1/notification/update-preferences", handleUpdateNotificationPreferences(client)},
			{"GET", "/v1/notification/preferences", handleGetNotificationPreferences(client)},
		}

		for _, h := range handlers {
			if err := mux.HandlePath(h.method, h.pattern, h.handler); err != nil {
				return err
			}
		}

		return nil
	})
}

func handleSendNotification(client notification_service.NotificationServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req notification_service.SendNotificationRequest
		if err := decodeJSONBody(w, r, &req); err != nil {
			return
		}

		if req.UserId == "" || req.Message == "" {
			http.Error(w, "Fields user_id and message are required", http.StatusBadRequest)
			return
		}

		resp, err := client.SendNotification(r.Context(), &req)
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		writeJSONResponse(w, http.StatusCreated, resp)
	}
}

func handleGetNotifications(client notification_service.NotificationServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req notification_service.GetNotificationsRequest
		req.UserId = parseStringParam(r, "user_id", "")

		if req.UserId == "" {
			http.Error(w, "Parameter user_id is required", http.StatusBadRequest)
			return
		}

		req.Limit = int32(parseIntParam(r, "limit", 10))
		req.Offset = int32(parseIntParam(r, "offset", 0))

		resp, err := client.GetNotifications(r.Context(), &req)
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		writeJSONResponse(w, http.StatusOK, resp)
	}
}

func handleMarkNotificationAsRead(client notification_service.NotificationServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req notification_service.MarkNotificationAsReadRequest
		if err := decodeJSONBody(w, r, &req); err != nil {
			return
		}

		if req.NotificationId == "" || req.UserId == "" {
			http.Error(w, "Fields notification_id and user_id are required", http.StatusBadRequest)
			return
		}

		resp, err := client.MarkNotificationAsRead(r.Context(), &req)
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		writeJSONResponse(w, http.StatusOK, resp)
	}
}

func handleUpdateNotificationPreferences(client notification_service.NotificationServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req notification_service.UpdateNotificationPreferencesRequest
		if err := decodeJSONBody(w, r, &req); err != nil {
			return
		}

		if req.UserId == "" || req.Preferences == nil {
			http.Error(w, "Fields user_id and preferences are required", http.StatusBadRequest)
			return
		}

		resp, err := client.UpdateNotificationPreferences(r.Context(), &req)
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		writeJSONResponse(w, http.StatusOK, resp)
	}
}

func handleGetNotificationPreferences(client notification_service.NotificationServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req notification_service.GetNotificationPreferencesRequest
		req.UserId = parseStringParam(r, "user_id", "")

		if req.UserId == "" {
			http.Error(w, "Parameter user_id is required", http.StatusBadRequest)
			return
		}

		resp, err := client.GetNotificationPreferences(r.Context(), &req)
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		writeJSONResponse(w, http.StatusOK, resp)
	}
}
