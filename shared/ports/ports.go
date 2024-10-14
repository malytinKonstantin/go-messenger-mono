package ports

type ServicePorts struct {
	HTTP string
	GRPC string
}

var (
	APIGateway   = ServicePorts{HTTP: "3000"}
	Auth         = ServicePorts{HTTP: "3001", GRPC: "50051"}
	User         = ServicePorts{HTTP: "3002", GRPC: "50052"}
	Friendship   = ServicePorts{HTTP: "3003", GRPC: "50053"}
	Messaging    = ServicePorts{HTTP: "3004", GRPC: "50054"}
	Notification = ServicePorts{HTTP: "3005", GRPC: "50055"}
)
