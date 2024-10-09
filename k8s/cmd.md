### Применение конфигураций
```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/common/configmap.yaml
kubectl apply -f k8s/auth-service/secrets.yaml
kubectl apply -f k8s/api-gateway/
kubectl apply -f k8s/auth-service/
kubectl apply -f k8s/ingress.yaml
```

### Проверка ресурсов

#### Проверка работающих подов
```bash
kubectl get pods -n go-messenger
```

#### Проверка сервисов
```bash
kubectl get services -n go-messenger
```

#### Проверка деплойментов
```bash
kubectl get deployments -n go-messenger
```

#### Проверка Ingress
```bash
kubectl get ingress -n go-messenger
```

#### Проверка конфигмапов и секретов
```bash
kubectl get configmaps,secrets -n go-messenger
```

### Дополнительные команды

#### Получение подробной информации о ресурсе
Используйте команду `describe`. Например:
```bash
kubectl describe deployment auth-service -n go-messenger
```

#### Просмотр логов пода
```bash
kubectl logs <имя-пода> -n go-messenger
```

#### Просмотр всех ресурсов в пространстве имен
```bash
kubectl get all -n go-messenger
```

#### Получение IP-адреса Minikube
```bash
minikube ip
```

#### Порт-форвардинг сервиса Postgres
```bash
kubectl port-forward svc/postgres 5432:5432 -n go-messenger
```

#### Выполнение миграций базы данных
```bash
migrate -path ./migrations -database "postgres://postgres:password@localhost:5432/auth_db?sslmode=disable" up
```

#### Перезапуск деплойментов
```bash
kubectl rollout restart deployment api-gateway -n go-messenger
kubectl rollout restart deployment auth-service -n go-messenger
```

#### Получение событий в пространстве имен
```bash
kubectl get events -n go-messenger --sort-by='.metadata.creationTimestamp'
```

#### Повторное получение списка сервисов
```bash
kubectl get services -n go-messenger
kubectl get services -n go-messenger
```

#### Работа с секретами базы данных
Получение всех секретов в формате YAML:
```bash
kubectl get secret auth-db-secrets -n go-messenger -o yaml
```
Декодирование имени пользователя:
```bash
kubectl get secret auth-db-secrets -n go-messenger -o jsonpath="{.data.username}" | base64 --decode
```
Декодирование пароля:
```bash
kubectl get secret auth-db-secrets -n go-messenger -o jsonpath="{.data.password}" | base64 --decode
```

#### Подключение к поду
```bash
kubectl exec -it auth-service-78799cc97-bqhjw -n go-messenger -- sh
```

#### Описание сервиса Postgres
```bash
kubectl describe service postgres -n go-messenger
```

#### Строка подключения к базе данных Postgres
```bash
postgres://postgres:postgres@<EXTERNAL_IP>:5432/auth_db
```

#### Дополнительные команды для работы с базой данных
Порт-форвардинг и подключение к psql:
```bash
kubectl port-forward svc/postgres 5432:5432 -n go-messenger
psql -h localhost -p 5432 -U postgres -d auth_db
```

#### Получение URL внешнего сервиса Postgres через Minikube
```bash
minikube service postgres-external -n go-messenger --url
```


kubectl patch service auth-service -n go-messenger -p '{"spec":{"selector":{"app":"auth-service","version":"green"}}}'

kubectl patch service api-gateway -n go-messenger -p '{"spec":{"selector":{"app":"api-gateway","version":"green"}}}'


kubectl delete deployment auth-service-blue -n go-messenger
kubectl delete deployment api-gateway-blue -n go-messenger