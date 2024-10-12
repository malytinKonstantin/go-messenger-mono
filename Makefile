# Переменные

DOCKER_HOST := unix:///var/run/docker.sock
DOCKER_REGISTRY := constmalytin
SERVICES := api-gateway auth-service friendship-service messaging-service
K8S_NAMESPACE := go-messenger
VERSION ?= $(shell git rev-parse --short HEAD)
MESSAGING_SERVICE_DIR := ./messaging-service

GOCMD = go
GOBUILD = $(GOCMD) build

# Цели
.PHONY: all build clean test coverage deps fmt lint vet check push deploy update-images deploy-all \
        deploy-blue deploy-green switch-to-blue switch-to-green log \
        check-events deploy-auth-postgres deploy-friendship-neo4j deploy-cassandra deploy-common deploy-services deploy-ingress delete-blue delete-green \
        release-blue release-green deploy-databases \
        build-cassandra push-cassandra run-cassandra stop-cassandra restart-cassandra \
        build-messaging-service push-messaging-service

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

# Сборка и пуш
build-and-push: build push

# Обновление образов в деплойментах
update-images:
	@for service in $(SERVICES); do \
		echo "Updating image for $$service to version $(VERSION)..."; \
		kubectl set image deployment/$$service $$service=$(DOCKER_REGISTRY)/$$service:$(VERSION) -n $(K8S_NAMESPACE); \
	done

# Деплой общих ресурсов
deploy-common:
	kubectl apply -f k8s/namespace.yaml
	kubectl apply -f k8s/common/configmap.yaml -n $(K8S_NAMESPACE)
	kubectl apply -f k8s/auth-service/secrets.yaml -n $(K8S_NAMESPACE)
	kubectl apply -f k8s/friendship-service/secrets.yaml -n $(K8S_NAMESPACE)
	kubectl apply -f k8s/messaging-service/configmap.yaml -n $(K8S_NAMESPACE)
	# kubectl apply -f k8s/common/secrets.yaml -n $(K8S_NAMESPACE)

# Деплой баз данных
deploy-auth-postgres:
	kubectl apply -f k8s/auth-service/persistent-volume.yaml -n $(K8S_NAMESPACE)
	kubectl apply -f k8s/auth-service/persistent-volume-claim.yaml -n $(K8S_NAMESPACE)
	kubectl apply -f k8s/auth-service/deployment-postgres.yaml -n $(K8S_NAMESPACE)
	kubectl apply -f k8s/auth-service/postgres-nodeport-service.yaml -n $(K8S_NAMESPACE)
	kubectl apply -f k8s/auth-service/postgres-service.yaml -n $(K8S_NAMESPACE)

deploy-friendship-neo4j:
	kubectl apply -f k8s/friendship-service/neo4j-volume.yaml -n $(K8S_NAMESPACE)
	kubectl apply -f k8s/friendship-service/neo4j-volume-claim.yaml -n $(K8S_NAMESPACE)
	kubectl apply -f k8s/friendship-service/neo4j-deployment.yaml -n $(K8S_NAMESPACE)
	kubectl apply -f k8s/friendship-service/neo4j-service.yaml -n $(K8S_NAMESPACE)

deploy-cassandra:
	kubectl apply -f k8s/messaging-service/cassandra-pv.yaml -n $(K8S_NAMESPACE)
	kubectl apply -f k8s/messaging-service/cassandra-pvc.yaml -n $(K8S_NAMESPACE)
	kubectl apply -f k8s/messaging-service/cassandra-deployment.yaml -n $(K8S_NAMESPACE)
	kubectl apply -f k8s/messaging-service/cassandra-service.yaml -n $(K8S_NAMESPACE)

deploy-databases: deploy-auth-postgres deploy-friendship-neo4j deploy-cassandra

# Деплой сервисов
deploy-services:
	@for service in $(SERVICES); do \
		kubectl apply -f k8s/$$service/service.yaml -n $(K8S_NAMESPACE); \
	done

# Деплой Ingress
deploy-ingress:
	kubectl apply -f k8s/ingress.yaml -n $(K8S_NAMESPACE)
	# kubectl apply -f k8s/registry/deployment.yaml -n $(K8S_NAMESPACE)

# Деплой blue версии
deploy-blue:
	@for service in $(SERVICES); do \
		echo "Deploying $$service blue version..."; \
		envsubst < k8s/$$service/deployment-blue.yaml | kubectl apply -n $(K8S_NAMESPACE) -f -; \
	done

# Деплой green версии
deploy-green:
	@for service in $(SERVICES); do \
		echo "Deploying $$service green version..."; \
		envsubst < k8s/$$service/deployment-green.yaml | kubectl apply -n $(K8S_NAMESPACE) -f -; \
	done

# Переключение на blue версию
switch-to-blue:
	@for service in $(SERVICES); do \
		echo "Switching $$service service to blue version..."; \
		kubectl patch service $$service -n $(K8S_NAMESPACE) -p "{\"spec\":{\"selector\":{\"app\":\"$$service\",\"version\":\"blue\"}}}"; \
	done

# Переключение на green версию
switch-to-green:
	@for service in $(SERVICES); do \
		echo "Switching $$service service to green version..."; \
		kubectl patch service $$service -n $(K8S_NAMESPACE) -p "{\"spec\":{\"selector\":{\"app\":\"$$service\",\"version\":\"green\"}}}"; \
	done

# Удаление blue версий деплойментов
delete-blue:
	@for service in $(SERVICES); do \
		echo "Deleting $$service blue deployment..."; \
		kubectl delete deployment $$service-blue -n $(K8S_NAMESPACE) --ignore-not-found; \
	done

# Удаление green версий деплойментов
delete-green:
	@for service in $(SERVICES); do \
		echo "Deleting $$service green deployment..."; \
		kubectl delete deployment $$service-green -n $(K8S_NAMESPACE) --ignore-not-found; \
	done

# Полный цикл релиза для blue версии
release-blue: build-and-push deploy-common deploy-databases deploy-services deploy-blue deploy-ingress
	@echo "Waiting for blue deployments to be ready..."
	@for service in $(SERVICES); do \
		kubectl rollout status deployment/$$service-blue -n $(K8S_NAMESPACE); \
	done
	@echo "Switching traffic to blue version..."
	@$(MAKE) switch-to-blue
	@echo "Waiting for traffic to stabilize..."
	@sleep 30
	@echo "Deleting green deployments..."
	@$(MAKE) delete-green

# Полный цикл релиза для green версии
release-green: build-and-push deploy-common deploy-databases deploy-services deploy-green deploy-ingress
	@echo "Waiting for green deployments to be ready..."
	@for service in $(SERVICES); do \
		kubectl rollout status deployment/$$service-green -n $(K8S_NAMESPACE); \
	done
	@echo "Switching traffic to green version..."
	@$(MAKE) switch-to-green
	@echo "Waiting for traffic to stabilize..."
	@sleep 30
	@echo "Deleting blue deployments..."
	@$(MAKE) delete-blue

# Логи и проверки
log:
	kubectl get pods -n $(K8S_NAMESPACE)
	kubectl get services -n $(K8S_NAMESPACE)
	kubectl get deployments -n $(K8S_NAMESPACE)
	kubectl get ingress -n $(K8S_NAMESPACE)
	kubectl get configmaps,secrets -n $(K8S_NAMESPACE)
# $(MAKE) check-events

check-events:
	kubectl get events -n $(K8S_NAMESPACE) --sort-by='.metadata.creationTimestamp'

# Команды для локального запуска Cassandra
build-cassandra:
	docker build -t $(DOCKER_REGISTRY)/messaging-service-cassandra:latest -f $(MESSAGING_SERVICE_DIR)/Dockerfile.cassandra $(MESSAGING_SERVICE_DIR)

push-cassandra:
	docker push $(DOCKER_REGISTRY)/messaging-service-cassandra:latest

run-cassandra:
	docker run --name my-cassandra \
		-p 9042:9042 \
		-d $(DOCKER_REGISTRY)/messaging-service-cassandra:latest

stop-cassandra:
	docker rm -f my-cassandra

restart-cassandra: stop-cassandra run-cassandra

# Команды для messaging-service
build-messaging-service:
	docker build -t $(DOCKER_REGISTRY)/messaging-service:$(VERSION) -f $(MESSAGING_SERVICE_DIR)/Dockerfile $(MESSAGING_SERVICE_DIR)

push-messaging-service:
	docker push $(DOCKER_REGISTRY)/messaging-service:$(VERSION)