run dev local
env $(cat .env | xargs) go run cmd/server/main.go