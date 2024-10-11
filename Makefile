# Переменные

DOCKER_HOST := unix:///var/run/docker.sock
DOCKER_REGISTRY := constmalytin
SERVICES := api-gateway auth-service friendship-service
K8S_NAMESPACE := go-messenger
VERSION ?= $(shell git rev-parse --short HEAD)

GOCMD = go
GOBUILD = $(GOCMD) build

# Цели
.PHONY: all build clean test coverage deps fmt lint vet check push deploy update-images deploy-all \
        deploy-blue deploy-green switch-to-blue switch-to-green log \
        check-events deploy-auth-postgres deploy-common deploy-services deploy-ingress delete-blue delete-green \
        release-blue release-green deploy-friendship-neo4j

all: build

clean:
	$(GOCMD) clean

test:
	$(GOCMD) test -v ./...

coverage:
	$(GOCMD) test -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

deps:
	$(GOCMD) mod tidy

fmt:
	$(GOCMD) fmt ./...

lint:
	golangci-lint run ./...

vet:
	$(GOCMD) vet ./...

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

# Деплой общих ресурсов
deploy-common:
	kubectl apply -f k8s/namespace.yaml
	kubectl apply -f k8s/common/configmap.yaml
	kubectl apply -f k8s/auth-service/secrets.yaml
	kubectl apply -f k8s/api-gateway/secrets.yaml
	kubectl apply -f k8s/friendship-service/secrets.yaml
	$(MAKE) deploy-auth-postgres
	$(MAKE) deploy-friendship-neo4j

# Деплой PostgreSQL для auth-service
deploy-auth-postgres:
	kubectl apply -f k8s/auth-service/persistent-volume.yaml
	kubectl apply -f k8s/auth-service/persistent-volume-claim.yaml
	kubectl apply -f k8s/auth-service/deployment-postgres.yaml
	kubectl apply -f k8s/auth-service/postgres-service.yaml

# Деплой Neo4j для friendship-service
deploy-friendship-neo4j:
	kubectl apply -f k8s/friendship-service/neo4j-volume.yaml
	kubectl apply -f k8s/friendship-service/neo4j-volume-claim.yaml
	kubectl apply -f k8s/friendship-service/neo4j-deployment.yaml
	kubectl apply -f k8s/friendship-service/neo4j-service.yaml

# Деплой сервисов
deploy-services:
	@for service in $(SERVICES); do \
		if [ -f k8s/$$service/service.yaml ]; then \
			echo "Deploying service for $$service..."; \
			kubectl apply -f k8s/$$service/service.yaml; \
		fi; \
	done

# Деплой blue версии
deploy-blue:
	@for service in $(SERVICES); do \
		if [ -f k8s/$$service/deployment-blue.yaml ]; then \
			echo "Deploying $$service to blue environment with version $(VERSION)..."; \
			VERSION=$(VERSION) envsubst < k8s/$$service/deployment-blue.yaml | kubectl apply -f -; \
		fi; \
	done

# Деплой green версии
deploy-green:
	@for service in $(SERVICES); do \
		if [ -f k8s/$$service/deployment-green.yaml ]; then \
			echo "Deploying $$service to green environment with version $(VERSION)..."; \
			VERSION=$(VERSION) envsubst < k8s/$$service/deployment-green.yaml | kubectl apply -f -; \
		fi; \
	done

# Переключение на blue версию
switch-to-blue:
	@for service in $(SERVICES); do \
		echo "Switching $$service to blue version..."; \
		kubectl patch service $$service -n $(K8S_NAMESPACE) -p '{"spec":{"selector":{"app":"'"$$service"'","version":"blue"}}}'; \
	done

# Переключение на green версию
switch-to-green:
	@for service in $(SERVICES); do \
		echo "Switching $$service to green version..."; \
		kubectl patch service $$service -n $(K8S_NAMESPACE) -p '{"spec":{"selector":{"app":"'"$$service"'","version":"green"}}}'; \
	done

# Удаление синих деплойментов
delete-blue:
	@for service in $(SERVICES); do \
		echo "Deleting $$service blue deployment..."; \
		kubectl delete deployment $$service-blue -n $(K8S_NAMESPACE) --ignore-not-found; \
	done

# Удаление зеленых деплойментов
delete-green:
	@for service in $(SERVICES); do \
		echo "Deleting $$service green deployment..."; \
		kubectl delete deployment $$service-green -n $(K8S_NAMESPACE) --ignore-not-found; \
	done

# Применение Ingress и других ресурсов
deploy-ingress:
	kubectl apply -f k8s/ingress.yaml
	kubectl apply -f k8s/registry/deployment.yaml

# Общий деплой ресурсов Kubernetes
deploy: deploy-common deploy-services deploy-ingress

# Полный цикл релиза для blue версии
release-blue: build-and-push deploy deploy-blue
	@echo "Waiting for blue deployment to be ready..."
	@kubectl rollout status deployment/api-gateway-blue -n $(K8S_NAMESPACE)
	@kubectl rollout status deployment/auth-service-blue -n $(K8S_NAMESPACE)
	@kubectl rollout status deployment/friendship-service-blue -n $(K8S_NAMESPACE)
	@echo "Switching traffic to blue version..."
	@$(MAKE) switch-to-blue
	@echo "Waiting for traffic to stabilize..."
	@sleep 30
	@echo "Deleting green deployment..."
	@$(MAKE) delete-green

# Полный цикл релиза для green версии
release-green: build-and-push deploy deploy-green
	@echo "Waiting for green deployment to be ready..."
	@kubectl rollout status deployment/api-gateway-green -n $(K8S_NAMESPACE)
	@kubectl rollout status deployment/auth-service-green -n $(K8S_NAMESPACE)
	@kubectl rollout status deployment/friendship-service-green -n $(K8S_NAMESPACE)
	@echo "Switching traffic to green version..."
	@$(MAKE) switch-to-green
	@echo "Waiting for traffic to stabilize..."
	@sleep 30
	@echo "Deleting blue deployment..."
	@$(MAKE) delete-blue

log:
	kubectl get pods -n $(K8S_NAMESPACE)
	kubectl get services -n $(K8S_NAMESPACE)
	kubectl get deployments -n $(K8S_NAMESPACE)
	kubectl get ingress -n $(K8S_NAMESPACE)
	kubectl get configmaps,secrets -n $(K8S_NAMESPACE)
	$(MAKE) check-events

check-events:
	kubectl get events -n $(K8S_NAMESPACE) --sort-by='.metadata.creationTimestamp'