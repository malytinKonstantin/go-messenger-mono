package handlers

import (
	"net/http"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"
)

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

func handleGrpcError(w http.ResponseWriter, err error) {
	grpcStatus, ok := status.FromError(err)
	if ok {
		httpCode := runtime.HTTPStatusFromCode(grpcStatus.Code())
		http.Error(w, grpcStatus.Message(), httpCode)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
