# Переменные
DOCKER_REGISTRY := localhost:5000
SERVICES := api-gateway auth-service
DOCKER_COMPOSE = docker-compose
K8S_NAMESPACE := go-messenger
SRC_DIRS := ./...
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
.PHONY: all build clean test coverage deps fmt lint vet

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

build:
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		$(DOCKER_COMPOSE) -f ./$$service/docker-compose.yml build; \
	done

up:
	@for service in $(SERVICES); do \
		echo "Starting $$service..."; \
		$(DOCKER_COMPOSE) -f ./$$service/docker-compose.yml up -d; \
	done

down:
	@for service in $(SERVICES); do \
		echo "Stopping $$service..."; \
		$(DOCKER_COMPOSE) -f ./$$service/docker-compose.yml down; \
	done

push:
	@for service in $(SERVICES); do \
		echo "Pushing $$service..."; \
		docker push $(DOCKER_REGISTRY)/$$service:latest; \
	done

deploy:
	kubectl apply -f k8s/namespace.yaml
	kubectl apply -f k8s/common/configmap.yaml
	kubectl apply -f k8s/auth-service/secrets.yaml
	kubectl apply -f k8s/api-gateway/
	kubectl apply -f k8s/auth-service/
	kubectl apply -f k8s/ingress.yaml

# Цель для обновления образов в деплойментах
update-images:
	@for service in $(SERVICES); do \
		kubectl set image deployment/$$service $$service=$(DOCKER_REGISTRY)/$$service:latest -n $(K8S_NAMESPACE); \
	done