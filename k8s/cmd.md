 // Start of Selection
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

#### Получить IP-адрес Minikube:

```bash
minikube ip
```