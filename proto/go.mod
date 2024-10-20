module github.com/malytinKonstantin/go-messenger-mono/proto

go 1.21

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.35.1-20240920164238-5a7b106cbb87.1
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.22.0
	google.golang.org/genproto/googleapis/api v0.0.0-20240930140551-af27646dc61f
	google.golang.org/grpc v1.67.1
	google.golang.org/protobuf v1.35.1
)

require (
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240930140551-af27646dc61f // indirect
)

replace (
    github.com/malytinKonstantin/go-messenger-mono/proto/auth-service => ./pkg/api/auth_service/v1
    github.com/malytinKonstantin/go-messenger-mono/proto/friendship-service => ./pkg/api/friendship_service/v1
    github.com/malytinKonstantin/go-messenger-mono/proto/messaging-service => ./pkg/api/messaging_service/v1
    github.com/malytinKonstantin/go-messenger-mono/proto/notification-service => ./pkg/api/notification_service/v1
    github.com/malytinKonstantin/go-messenger-mono/proto/user-service => ./pkg/api/user_service/v1
)