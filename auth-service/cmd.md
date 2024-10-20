run dev local
env $(cat .env.dev | xargs) go run cmd/server/main.go