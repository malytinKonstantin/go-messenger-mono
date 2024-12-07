run dev local
env $(cat .env.dev | xargs) go run cmd/server/main.go


mockery --dir=../../repository --output=../mocks --name=UserCredentialsRepository
mockery --dir=../../repository --output=../mocks --name=ResetPasswordTokenRepository
mockery --dir=../../repository --output=../mocks --name=OauthAccountRepository

go install github.com/vektra/mockery/v2@latest
mockery --version

go test ./...
go test -v ./internal/usecase/auth/tests/...
go test -cover ./internal/usecase/auth/tests/...