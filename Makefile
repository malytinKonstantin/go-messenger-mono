include k8s.mk

lint:
	@for dir in api-gateway auth-service friendship-service messaging-service notification-service proto shared user-service; do \
		echo "Linting $$dir..."; \
		(cd $$dir && golangci-lint run --fix) || exit 1; \
	done