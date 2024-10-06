package handlers

import (
	"context"

	"github.com/google/uuid"
	"github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/services"
	pb "github.com/malytinKonstantin/go-messenger-mono/proto/auth-service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user, err := h.authService.RegisterUser(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to register user: %v", err)
	}
	return &pb.RegisterResponse{UserId: user.UserID.String()}, nil
}

func (h *AuthHandler) Authenticate(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	_, err := h.authService.AuthenticateUser(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Authentication failed: %v", err)
	}
	// Здесь нужно реализовать генерацию JWT токена
	token := "dummy_token" // Замените на реальную генерацию токена
	return &pb.AuthResponse{Token: token}, nil
}

func (h *AuthHandler) OAuthAuthenticate(ctx context.Context, req *pb.OAuthRequest) (*pb.AuthResponse, error) {
	// Здесь нужно реализовать OAuth аутентификацию
	return nil, status.Errorf(codes.Unimplemented, "OAuth authentication not implemented")
}

func (h *AuthHandler) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid user ID: %v", err)
	}
	err = h.authService.VerifyUser(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to verify email: %v", err)
	}
	return &pb.VerifyEmailResponse{Success: true}, nil
}

func (h *AuthHandler) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	// Здесь нужно реализовать логику сброса пароля
	return nil, status.Errorf(codes.Unimplemented, "Reset password not implemented")
}

func (h *AuthHandler) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	// Здесь нужно реализовать логику изменения пароля
	return nil, status.Errorf(codes.Unimplemented, "Change password not implemented")
}
