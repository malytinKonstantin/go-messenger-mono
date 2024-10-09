make all
make update-images

make build-service SERVICE=auth-service
make push-service SERVICE=auth-service

kubectl delete deployment auth-service -n go-messenger
kubectl apply -f k8s/auth-service/deployment.yaml