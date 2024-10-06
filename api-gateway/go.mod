module github.com/malytinKonstantin/go-messenger-mono/api-gateway

go 1.23.1

require (
	github.com/gofiber/fiber/v2 v2.52.5
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.1
	github.com/malytinKonstantin/go-messenger-mono/proto v0.0.0
	google.golang.org/grpc v1.67.1
)

require (
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.17.10 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.56.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240814211410-ddb44dafa142 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240930140551-af27646dc61f // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)

replace github.com/malytinKonstantin/go-messenger-mono/proto => ../proto
