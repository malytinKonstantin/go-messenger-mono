package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	notification_service "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/notification_service/v1"
)

func RegisterNotificationService(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return err
	}

	client := notification_service.NewNotificationServiceClient(conn)

	err = mux.HandlePath("POST", "/v1/notification/send", handleSendNotification(client))
	if err != nil {
		return err
	}

	err = mux.HandlePath("GET", "/v1/notification", handleGetNotifications(client))
	if err != nil {
		return err
	}

	err = mux.HandlePath("PUT", "/v1/notification/mark-read", handleMarkNotificationAsRead(client))
	if err != nil {
		return err
	}

	err = mux.HandlePath("PUT", "/v1/notification/update-preferences", handleUpdateNotificationPreferences(client))
	if err != nil {
		return err
	}

	err = mux.HandlePath("GET", "/v1/notification/preferences", handleGetNotificationPreferences(client))
	if err != nil {
		return err
	}

	return nil
}

func handleSendNotification(client notification_service.NotificationServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req notification_service.SendNotificationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
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
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func handleGetNotifications(client notification_service.NotificationServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req notification_service.GetNotificationsRequest
		req.UserId = r.URL.Query().Get("user_id")

		if req.UserId == "" {
			http.Error(w, "Parameter user_id is required", http.StatusBadRequest)
			return
		}

		limit := parseIntParam(r, "limit", 10)
		if limit < 0 {
			http.Error(w, "Parameter limit cannot be negative", http.StatusBadRequest)
			return
		}
		req.Limit = int32(limit)

		offset := parseIntParam(r, "offset", 0)
		if offset < 0 {
			http.Error(w, "Parameter offset cannot be negative", http.StatusBadRequest)
			return
		}
		req.Offset = int32(offset)

		resp, err := client.GetNotifications(r.Context(), &req)
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func handleMarkNotificationAsRead(client notification_service.NotificationServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req notification_service.MarkNotificationAsReadRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
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
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func handleUpdateNotificationPreferences(client notification_service.NotificationServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req notification_service.UpdateNotificationPreferencesRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
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
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func handleGetNotificationPreferences(client notification_service.NotificationServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req notification_service.GetNotificationPreferencesRequest
		req.UserId = r.URL.Query().Get("user_id")

		if req.UserId == "" {
			http.Error(w, "Parameter user_id is required", http.StatusBadRequest)
			return
		}

		resp, err := client.GetNotificationPreferences(r.Context(), &req)
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
