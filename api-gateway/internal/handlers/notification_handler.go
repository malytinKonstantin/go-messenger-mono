package handlers

import (
	"context"
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

	if err := mux.HandlePath("POST", "/v1/notification/send", handleSendNotification(client)); err != nil {
		return err
	}
	if err := mux.HandlePath("GET", "/v1/notification", handleGetNotifications(client)); err != nil {
		return err
	}
	if err := mux.HandlePath("PUT", "/v1/notification/mark-read", handleMarkNotificationAsRead(client)); err != nil {
		return err
	}
	if err := mux.HandlePath("PUT", "/v1/notification/update-preferences", handleUpdateNotificationPreferences(client)); err != nil {
		return err
	}
	if err := mux.HandlePath("GET", "/v1/notification/preferences", handleGetNotificationPreferences(client)); err != nil {
		return err
	}

	return nil
}

func handleSendNotification(client notification_service.NotificationServiceClient) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		var req notification_service.SendNotificationRequest
		if err := decodeJSONBody(w, r, &req); err != nil {
			return
		}

		if req.UserId == "" || req.Message == "" {
			http.Error(w, "Поля user_id и message обязательны", http.StatusBadRequest)
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
			http.Error(w, "Параметр user_id обязателен", http.StatusBadRequest)
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
			http.Error(w, "Поля notification_id и user_id обязательны", http.StatusBadRequest)
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
			http.Error(w, "Поля user_id и preferences обязательны", http.StatusBadRequest)
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
			http.Error(w, "Параметр user_id обязателен", http.StatusBadRequest)
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
