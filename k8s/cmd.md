kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/common/configmap.yaml
kubectl apply -f k8s/auth-service/secrets.yaml
kubectl apply -f k8s/api-gateway/
kubectl apply -f k8s/auth-service/
kubectl apply -f k8s/ingress.yaml