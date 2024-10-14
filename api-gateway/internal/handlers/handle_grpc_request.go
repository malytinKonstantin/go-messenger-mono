package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func handleGrpcRequest(grpcFunc interface{}) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		ctx := r.Context()

		// Создаем пустой запрос нужного типа
		reqValue := reflect.New(reflect.TypeOf(grpcFunc).In(1).Elem()).Interface()

		// Парсим тело запроса
		body, _ := ioutil.ReadAll(r.Body)
		if err := json.Unmarshal(body, reqValue); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Вызываем gRPC функцию
		results := reflect.ValueOf(grpcFunc).Call([]reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(reqValue),
		})

		// Проверяем ошибку
		if err := results[1].Interface(); err != nil {
			http.Error(w, err.(error).Error(), http.StatusInternalServerError)
			return
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results[0].Interface())
	}
}
