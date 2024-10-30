package grpc

import (
	"context"

	pb "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/auth_service/v1"
)

type CombinedAuthHandler struct {
	pb.UnimplementedAuthServiceServer
	authHandler     *AuthHandler
	oauthHandler    *OAuthHandler
	passwordHandler *PasswordHandler
}

func NewCombinedAuthHandler(authHandler *AuthHandler, oauthHandler *OAuthHandler, passwordHandler *PasswordHandler) *CombinedAuthHandler {
	return &CombinedAuthHandler{
		authHandler:     authHandler,
		oauthHandler:    oauthHandler,
		passwordHandler: passwordHandler,
	}
}

func (h *CombinedAuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return h.authHandler.Register(ctx, req)
}

func (h *CombinedAuthHandler) Authenticate(ctx context.Context, req *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	return h.authHandler.Authenticate(ctx, req)
}

func (h *CombinedAuthHandler) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	return h.authHandler.VerifyEmail(ctx, req)
}

func (h *CombinedAuthHandler) OAuthAuthenticate(ctx context.Context, req *pb.OAuthAuthenticateRequest) (*pb.OAuthAuthenticateResponse, error) {
	return h.oauthHandler.OAuthAuthenticate(ctx, req)
}

func (h *CombinedAuthHandler) ResetPasswordRequest(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	return h.passwordHandler.ResetPasswordRequest(ctx, req)
}

func (h *CombinedAuthHandler) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	return h.passwordHandler.ChangePassword(ctx, req)
}
