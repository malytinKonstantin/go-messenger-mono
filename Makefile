# Переменные

DOCKER_HOST := unix:///var/run/docker.sock
DOCKER_REGISTRY := constmalytin
SERVICES := api-gateway auth-service
K8S_NAMESPACE := go-messenger
VERSION ?= $(shell git rev-parse --short HEAD)
DEPLOYMENT_VERSION ?= blue

GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod
GOFMT = gofmt
GOLINT = golangci-lint

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
	$(GOCMD) vet $(SRC_DIRS)

check: fmt lint vet test

# Сборка всех сервисов
build:
	@for service in $(SERVICES); do \
		echo "Building $$service with version $(VERSION)..."; \
		docker build -t $(DOCKER_REGISTRY)/$$service:$(VERSION) -f ./$$service/Dockerfile .; \
	done

# Пуш всех образов
push:
	@for service in $(SERVICES); do \
		echo "Pushing $$service with version $(VERSION)..."; \
		docker push $(DOCKER_REGISTRY)/$$service:$(VERSION); \
	done

# Сборка и пуш всех сервисов
build-and-push: build push

# Сборка отдельного сервиса
build-service:
	@if [ -z "$(SERVICE)" ]; then \
		echo "Укажите SERVICE=<имя_сервиса>"; \
		exit 1; \
	fi
	echo "Building $(SERVICE) with version $(VERSION)..."
	docker build -t $(DOCKER_REGISTRY)/$(SERVICE):$(VERSION) -f ./$(SERVICE)/Dockerfile .

# Пуш отдельного образа
push-service:
	@if [ -z "$(SERVICE)" ]; then \
		echo "Укажите SERVICE=<имя_сервиса>"; \
		exit 1; \
	fi
	echo "Pushing $(SERVICE) with version $(VERSION)..."
	docker push $(DOCKER_REGISTRY)/$(SERVICE):$(VERSION)

# Сборка и пуш отдельного сервиса
build-and-push-service: build-service push-service

# Обновление образов в Kubernetes для Blue-Green деплоя
update-images:
	@for service in $(SERVICES); do \
		kubectl set image deployment/$$service-$(DEPLOYMENT_VERSION) $$service=$(DOCKER_REGISTRY)/$$service:$(VERSION) -n $(K8S_NAMESPACE); \
	done

# Рестарт деплойментов (опционально)
restart-deployments:
	@for service in $(SERVICES); do \
		kubectl rollout restart deployment/$$service-$(DEPLOYMENT_VERSION) -n $(K8S_NAMESPACE); \
	done

# Цели для Blue-Green деплоя с правильными Docker сборками

deploy-blue: DEPLOYMENT_VERSION=blue
deploy-blue: build-and-push deploy-blue-k8s

deploy-green: DEPLOYMENT_VERSION=green
deploy-green: build-and-push deploy-green-k8s

deploy-blue-k8s:
	@for service in $(SERVICES); do \
		echo "Deploying $$service to blue environment..."; \
		export VERSION=$(VERSION); \
		cat k8s/$$service/deployment-blue.yaml | envsubst '$${VERSION}' | kubectl apply -f -; \
		kubectl apply -f k8s/$$service/service.yaml; \
	done

deploy-green-k8s:
	@for service in $(SERVICES); do \
		echo "Deploying $$service to green environment..."; \
		export VERSION=$(VERSION); \
		cat k8s/$$service/deployment-green.yaml | envsubst '$${VERSION}' | kubectl apply -f -; \
		kubectl apply -f k8s/$$service/service.yaml; \
	done

switch-to-blue:
	@for service in $(SERVICES); do \
		echo "Switching $$service to blue version..."; \
		kubectl patch service $$service -n $(K8S_NAMESPACE) -p '{"spec":{"selector":{"app":"'"$$service"'","version":"blue"}}}'; \
	done

switch-to-green:
	@for service in $(SERVICES); do \
		echo "Switching $$service to green version..."; \
		kubectl patch service $$service -n $(K8S_NAMESPACE) -p '{"spec":{"selector":{"app":"'"$$service"'","version":"green"}}}'; \
	done

# Деплой отдельных сервисов с учётом версии деплоя

deploy-auth-service:
	$(MAKE) build-service SERVICE=auth-service
	$(MAKE) push-service SERVICE=auth-service
	envsubst < k8s/auth-service/deployment-$(DEPLOYMENT_VERSION).yaml | kubectl apply -f -
	kubectl apply -f k8s/auth-service/service.yaml

deploy-api-gateway:
	$(MAKE) build-service SERVICE=api-gateway
	$(MAKE) push-service SERVICE=api-gateway
	envsubst < k8s/api-gateway/deployment-$(DEPLOYMENT_VERSION).yaml | kubectl apply -f -
	kubectl apply -f k8s/api-gateway/service.yaml

# Общий деплой ресурсов Kubernetes
deploy:
	kubectl apply -f k8s/namespace.yaml
	kubectl apply -f k8s/common/configmap.yaml
	$(MAKE) deploy-blue
	kubectl apply -f k8s/ingress.yaml
	kubectl apply -f k8s/registry/deployment.yaml
	$(MAKE) log

log:
	kubectl get pods -n $(K8S_NAMESPACE)
	kubectl get services -n $(K8S_NAMESPACE)
	kubectl get deployments -n $(K8S_NAMESPACE)
	kubectl get ingress -n $(K8S_NAMESPACE)
	kubectl get configmaps,secrets -n $(K8S_NAMESPACE)