include vendor.proto.mk

# Используем bin в текущей директории для установки плагинов protoc
LOCAL_BIN := $(CURDIR)/bin
export PATH := $(LOCAL_BIN):$(PATH)

# Опциональный параметр DIR для указания папки
DIR ?= .

.PHONY: .bin-deps proto-lint proto-breaking proto-generate proto-format proto-build protolint all

# Установка всех необходимых инструментов для работы с proto-файлами
.bin-deps: export GOBIN := $(LOCAL_BIN)
.bin-deps:
	@echo "Installing tools..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install github.com/yoheimuta/protolint/cmd/protolint@latest
	go install github.com/envoyproxy/protoc-gen-validate@v1.1.0
	

# Проверка синтаксиса и стиля proto-файлов с помощью buf
lint:
	$(LOCAL_BIN)/buf lint

# Проверка обратной совместимости proto-файлов с помощью buf
breaking:
	$(LOCAL_BIN)/buf breaking --against '.git#branch=main'

# Генерация кода из proto-файлов с помощью buf
generate:
	$(LOCAL_BIN)/buf generate --template buf.gen.yaml

# Форматирование proto-файлов с помощью buf
format:
	$(LOCAL_BIN)/buf format -w

# Сборка proto-файлов с помощью buf
build:
	$(LOCAL_BIN)/buf build

# Дополнительная проверка линтером protolint
protolint:
	find . -name "*.proto" | xargs $(LOCAL_BIN)/protolint lint

# Выполнение всех команд последовательно
all: .bin-deps lint breaking generate format build protolint
