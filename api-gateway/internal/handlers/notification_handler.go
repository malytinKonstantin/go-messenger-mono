package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

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

	err = mux.HandlePath("GET", "/v1/notification/get", handleGetNotifications(client))
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
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp, err := client.SendNotification(r.Context(), &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
		req.Limit = int32(parseIntParam(r, "limit", 10))
		req.Offset = int32(parseIntParam(r, "offset", 0))
		resp, err := client.GetNotifications(r.Context(), &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp, err := client.MarkNotificationAsRead(r.Context(), &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp, err := client.UpdateNotificationPreferences(r.Context(), &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
		resp, err := client.GetNotificationPreferences(r.Context(), &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func parseIntParam(r *http.Request, name string, defaultValue int) int {
	value := r.URL.Query().Get(name)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}
