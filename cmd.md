docker build --platform=linux/amd64 -f auth-service/Dockerfile -t auth-service .
docker build --platform=linux/amd64 -f api-gateway/Dockerfile -t api-gateway .
