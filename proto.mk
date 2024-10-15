# Используем bin в текущей директории для установки плагинов protoc
LOCAL_BIN := $(CURDIR)/bin
export PATH := $(LOCAL_BIN):$(PATH)

.PHONY: install-tools proto-lint proto-breaking proto-generate proto-format proto-build

# Установка всех необходимых инструментов для работы с proto-файлами
install-tools:
	@echo "Installing tools..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install github.com/yoheimuta/protolint/cmd/protolint@latest

# Проверка синтаксиса и стиля proto-файлов с помощью buf
proto-lint:
	@echo "Linting proto files..."
	buf lint

# Проверка обратной совместимости proto-файлов с помощью buf
proto-breaking:
	@echo "Checking for breaking changes..."
	buf breaking --against '.git#branch=main'

# Генерация кода из proto-файлов с помощью buf
proto-generate:
	@echo "Generating code from proto files..."
	buf generate --template buf.gen.yaml

# Форматирование proto-файлов с помощью buf
proto-format:
	@echo "Formatting proto files..."
	buf format -w

# Сборка proto-файлов с помощью buf
proto-build:
	@echo "Building proto files..."
	buf build

# Дополнительная проверка линтером protolint
protolint:
	@echo "Running protolint..."
	protolint lint ./proto

# Выполнение всех команд последовательно
all: install-tools proto-lint proto-breaking proto-generate proto-format proto-build protolint
	@echo "All tasks completed successfully!"