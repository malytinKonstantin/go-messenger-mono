module github.com/malytinKonstantin/go-messenger-mono

go 1.21

replace (
	github.com/malytinKonstantin/go-messenger-mono/api-gateway => ./api-gateway
	github.com/malytinKonstantin/go-messenger-mono/auth-service => ./auth-service
	github.com/malytinKonstantin/go-messenger-mono/friendship-service => ./friendship-service
	github.com/malytinKonstantin/go-messenger-mono/messaging-service => ./messaging-service
	github.com/malytinKonstantin/go-messenger-mono/notification-service => ./notification-service
	github.com/malytinKonstantin/go-messenger-mono/proto => ./proto
	github.com/malytinKonstantin/go-messenger-mono/user-service => ./user-service
)
