# Переменные
DOCKER_REGISTRY := localhost:5000
SERVICES := api-gateway auth-service
K8S_NAMESPACE := go-messenger

# Цели
.PHONY: all build push deploy

all: build push deploy

build:
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		docker build -t $(DOCKER_REGISTRY)/$$service:latest ./$$service; \
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