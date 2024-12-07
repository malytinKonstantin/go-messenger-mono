package handlers

import (
	"context"
	"net/http"
	"time"

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
			{"POST", "/v1/notification/send", withJWTValidation(handleSendNotification(client))},
			{"GET", "/v1/notification", withJWTValidation(handleGetNotifications(client))},
			{"PUT", "/v1/notification/mark-read", withJWTValidation(handleMarkNotificationAsRead(client))},
			{"PUT", "/v1/notification/update-preferences", withJWTValidation(handleUpdateNotificationPreferences(client))},
			{"GET", "/v1/notification/preferences", withJWTValidation(handleGetNotificationPreferences(client))},
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

		ctx, cancel := withTimeout(r.Context(), 5*time.Second)
		defer cancel()

		respInterface, err := cb.Execute(func() (interface{}, error) {
			return client.SendNotification(ctx, &req)
		})
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		resp := respInterface.(*notification_service.SendNotificationResponse)
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

		ctx, cancel := withTimeout(r.Context(), 5*time.Second)
		defer cancel()

		respInterface, err := cb.Execute(func() (interface{}, error) {
			return client.GetNotifications(ctx, &req)
		})
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		resp := respInterface.(*notification_service.GetNotificationsResponse)
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

		ctx, cancel := withTimeout(r.Context(), 5*time.Second)
		defer cancel()

		respInterface, err := cb.Execute(func() (interface{}, error) {
			return client.MarkNotificationAsRead(ctx, &req)
		})
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		resp := respInterface.(*notification_service.MarkNotificationAsReadResponse)
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

		ctx, cancel := withTimeout(r.Context(), 5*time.Second)
		defer cancel()

		respInterface, err := cb.Execute(func() (interface{}, error) {
			return client.UpdateNotificationPreferences(ctx, &req)
		})
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		resp := respInterface.(*notification_service.UpdateNotificationPreferencesResponse)
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

		ctx, cancel := withTimeout(r.Context(), 5*time.Second)
		defer cancel()

		respInterface, err := cb.Execute(func() (interface{}, error) {
			return client.GetNotificationPreferences(ctx, &req)
		})
		if err != nil {
			handleGrpcError(w, err)
			return
		}
		resp := respInterface.(*notification_service.GetNotificationPreferencesResponse)
		writeJSONResponse(w, http.StatusOK, resp)
	}
}
