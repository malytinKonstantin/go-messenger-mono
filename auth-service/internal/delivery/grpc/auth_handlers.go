package grpc

import (
	"context"

	"github.com/google/uuid"
	authUsecase "github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/usecase/auth"
	credentialsUsecase "github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/usecase/credentials"
	pb "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/auth_service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	registerUsecase     authUsecase.RegisterUserUsecase
	authenticateUsecase authUsecase.AuthenticateUserUsecase
	verifyEmailUsecase  credentialsUsecase.VerifyEmailUsecase
}

func NewAuthHandler(
	registerUC authUsecase.RegisterUserUsecase,
	authenticateUC authUsecase.AuthenticateUserUsecase,
	verifyEmailUC credentialsUsecase.VerifyEmailUsecase,
) *AuthHandler {
	return &AuthHandler{
		registerUsecase:     registerUC,
		authenticateUsecase: authenticateUC,
		verifyEmailUsecase:  verifyEmailUC,
	}
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	user, err := h.registerUsecase.Execute(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error registering user: %v", err)
	}

	return &pb.RegisterResponse{UserId: user.UserID}, nil
}

func (h *AuthHandler) Authenticate(ctx context.Context, req *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	token, err := h.authenticateUsecase.Execute(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication error: %v", err)
	}

	return &pb.AuthenticateResponse{Token: token}, nil
}

func (h *AuthHandler) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid UserId: %v", err)
	}

	err = h.verifyEmailUsecase.Execute(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error verifying email: %v", err)
	}

	return &pb.VerifyEmailResponse{Success: true}, nil
}
