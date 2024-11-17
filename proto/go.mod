module github.com/malytinKonstantin/go-messenger-mono/proto

go 1.22.7

toolchain go1.23.2

require (
	github.com/envoyproxy/protoc-gen-validate v1.1.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.22.0
	google.golang.org/genproto/googleapis/api v0.0.0-20241021214115-324edc3d5d38
	google.golang.org/grpc v1.68.0
	google.golang.org/protobuf v1.35.2
)

require (
	golang.org/x/net v0.31.0 // indirect
	golang.org/x/sys v0.27.0 // indirect
	golang.org/x/text v0.20.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241113202542-65e8d215514f // indirect
)

replace (
	github.com/malytinKonstantin/go-messenger-mono/proto/auth-service => ./pkg/api/auth_service/v1
	github.com/malytinKonstantin/go-messenger-mono/proto/friendship-service => ./pkg/api/friendship_service/v1
	github.com/malytinKonstantin/go-messenger-mono/proto/messaging-service => ./pkg/api/messaging_service/v1
	github.com/malytinKonstantin/go-messenger-mono/proto/notification-service => ./pkg/api/notification_service/v1
	github.com/malytinKonstantin/go-messenger-mono/proto/user-service => ./pkg/api/user_service/v1
)
