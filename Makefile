# Переменные

DOCKER_HOST := unix:///var/run/docker.sock
DOCKER_REGISTRY := constmalytin
SERVICES := api-gateway auth-service
K8S_NAMESPACE := go-messenger
VERSION ?= $(shell git rev-parse --short HEAD)

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOVET=$(GOCMD) vet
GOLINT=golangci-lint

# Пути
SRC_DIRS := ./...

# Цели
.PHONY: all build clean test coverage deps fmt lint vet check push deploy update-images

all: build

clean:
	$(GOCLEAN)

test:
	$(GOTEST) -v $(SRC_DIRS)

coverage:
	$(GOTEST) -v -coverprofile=coverage.out $(SRC_DIRS)
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

deps:
	$(GOGET) -v -t -d $(SRC_DIRS)
	$(GOMOD) tidy

fmt:
	$(GOFMT) -w .

lint:
	$(GOLINT) run

vet:
	$(GOVET) $(SRC_DIRS)

check: fmt lint vet test

# Цель сборки
build:
	@for service in $(SERVICES); do \
		echo "Building $$service with version $(VERSION)..."; \
		docker build -t $(DOCKER_REGISTRY)/$$service:$(VERSION) -f ./$$service/Dockerfile .; \
	done

# Цель пуша образов
push:
	@for service in $(SERVICES); do \
		echo "Pushing $$service with version $(VERSION)..."; \
		docker push $(DOCKER_REGISTRY)/$$service:$(VERSION); \
	done

# Сборка отдельного образа
build-service:
	@if [ -z "$(SERVICE)" ]; then \
		echo "Укажите SERVICE=<имя_сервиса>"; \
		exit 1; \
	fi
	@echo "Building $(SERVICE) with version $(VERSION)..."
	docker build -t $(DOCKER_REGISTRY)/$(SERVICE):$(VERSION) -f ./$(SERVICE)/Dockerfile .

# Пуш отдельного образа
push-service:
	@if [ -z "$(SERVICE)" ]; then \
		echo "Укажите SERVICE=<имя_сервиса>"; \
		exit 1; \
	fi
	@echo "Pushing $(SERVICE) with version $(VERSION)..."
	docker push $(DOCKER_REGISTRY)/$(SERVICE):$(VERSION)

# Обновление образов в Kubernetes
update-images:
	@for service in $(SERVICES); do \
		kubectl set image deployment/$$service $$service=$(DOCKER_REGISTRY)/$$service:$(VERSION) -n $(K8S_NAMESPACE); \
	done

rebuild-auth:
	docker build -t $(DOCKER_REGISTRY)/auth-service:$(VERSION) -f ./auth-service/Dockerfile .
	docker push $(DOCKER_REGISTRY)/auth-service:$(VERSION)
	kubectl set image deployment/auth-service auth-service=$(DOCKER_REGISTRY)/auth-service:$(VERSION) -n $(K8S_NAMESPACE)
	kubectl rollout restart deployment auth-service -n $(K8S_NAMESPACE)

deploy:
	kubectl apply -f k8s/namespace.yaml
	kubectl apply -f k8s/common/configmap.yaml
	kubectl apply -f k8s/auth-service/
	kubectl apply -f k8s/api-gateway/
	kubectl apply -f k8s/auth-service/
	kubectl apply -f k8s/ingress.yaml
	kubectl apply -f k8s/registry/deployment.yaml
	kubectl get pods -n $(K8S_NAMESPACE)
	kubectl get services -n $(K8S_NAMESPACE)
	kubectl get deployments -n $(K8S_NAMESPACE)
	kubectl get ingress -n $(K8S_NAMESPACE)
	kubectl get configmaps,secrets -n $(K8S_NAMESPACE)

log:
	kubectl get pods -n $(K8S_NAMESPACE)
	kubectl get services -n $(K8S_NAMESPACE)
	kubectl get deployments -n $(K8S_NAMESPACE)
	kubectl get ingress -n $(K8S_NAMESPACE)
	kubectl get configmaps,secrets -n $(K8S_NAMESPACE)