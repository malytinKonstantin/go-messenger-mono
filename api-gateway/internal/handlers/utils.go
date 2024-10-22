package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"
)

// parseIntParam извлекает целочисленный параметр из URL-запроса.
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

// parseStringParam извлекает строковый параметр из URL-запроса.
func parseStringParam(r *http.Request, name string, defaultValue string) string {
	value := r.URL.Query().Get(name)
	if value == "" {
		return defaultValue
	}
	return value
}

// handleGrpcError обрабатывает ошибки gRPC и возвращает соответствующий HTTP-статус.
func handleGrpcError(w http.ResponseWriter, err error) {
	grpcStatus, ok := status.FromError(err)
	if ok {
		httpCode := runtime.HTTPStatusFromCode(grpcStatus.Code())
		http.Error(w, grpcStatus.Message(), httpCode)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// decodeJSONBody декодирует JSON из тела запроса.
func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(dst); err != nil {
		http.Error(w, "Неверный формат JSON: "+err.Error(), http.StatusBadRequest)
		return err
	}
	return nil
}

// writeJSONResponse записывает ответ в формате JSON.
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
