package handlers

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jackc/pgx/v5/pgxpool"
	pb "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/auth_service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	producer *kafka.Producer
	db       *pgxpool.Pool
}

func NewAuthHandler(producer *kafka.Producer, db *pgxpool.Pool) *AuthHandler {
	return &AuthHandler{producer: producer, db: db}
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation error: %v", err)
	}
	// Моковый ответ
	return &pb.RegisterResponse{UserId: "mock-user-id-12345"}, nil
}

func (h *AuthHandler) Authenticate(ctx context.Context, req *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation error: %v", err)
	}
	// Моковый ответ
	return &pb.AuthenticateResponse{Token: "mock-jwt-token-12345"}, nil
}

func (h *AuthHandler) OAuthAuthenticate(ctx context.Context, req *pb.OAuthAuthenticateRequest) (*pb.OAuthAuthenticateResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation error: %v", err)
	}
	// Моковый ответ
	return &pb.OAuthAuthenticateResponse{Token: "mock-oauth-token-12345"}, nil
}

func (h *AuthHandler) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation error: %v", err)
	}
	// Моковый ответ
	return &pb.VerifyEmailResponse{Success: true}, nil
}

func (h *AuthHandler) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation error: %v", err)
	}
	// Моковый ответ
	return &pb.ResetPasswordResponse{Success: true}, nil
}

func (h *AuthHandler) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation error: %v", err)
	}
	// Моковый ответ
	return &pb.ChangePasswordResponse{Success: true}, nil
}
