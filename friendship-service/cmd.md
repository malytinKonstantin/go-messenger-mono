run dev local
env $(cat .env.dev | xargs) go run cmd/server/main.go


$GOPATH/bin/mockery --dir=./internal/repositories --output=./internal/usecase/friendship/mocks --name=FriendRequestRepository


$GOPATH/bin/mockery --dir=./internal/repositories --output=./internal/usecase/friendship/mocks --name=UserRepository
