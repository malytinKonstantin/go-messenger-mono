docker build --platform=linux/amd64 -f auth-service/Dockerfile -t auth-service .
docker build --platform=linux/amd64 -f api-gateway/Dockerfile -t api-gateway .

echo -n "postgres" | base64

kubectl logs api-gateway-7df58c48dd-9fhcp -n go-messenger

kubectl create configmap api-gateway-env-file --from-file=.env -n go-messenger

kubectl create secret generic api-gateway-env-secret --from-file=api-gateway/.env -n go-messenger --dry-run=client -o yaml > k8s/api-gateway/secret-env.yaml

docker-compose --env-file .env.local -f docker-compose-dev.yaml up -d

docker-compose --env-file .env.local -f docker-compose-dev.yaml up -d --build