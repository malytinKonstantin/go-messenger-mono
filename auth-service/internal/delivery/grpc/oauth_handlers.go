package grpc

import (
	"context"

	oauthUsecase "github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/usecase/oauth"
	pb "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/auth_service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OAuthHandler struct {
	pb.UnimplementedAuthServiceServer
	oauthAuthenticateUsecase oauthUsecase.OAuthAuthenticateUsecase
}

func NewOAuthHandler(oauthAuthUC oauthUsecase.OAuthAuthenticateUsecase) *OAuthHandler {
	return &OAuthHandler{
		oauthAuthenticateUsecase: oauthAuthUC,
	}
}

func (h *OAuthHandler) OAuthAuthenticate(ctx context.Context, req *pb.OAuthAuthenticateRequest) (*pb.OAuthAuthenticateResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	user, err := h.oauthAuthenticateUsecase.Execute(ctx, req.Provider, req.ProviderUserId, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "oauth authentication error: %v", err)
	}

	token, err := h.oauthAuthenticateUsecase.GenerateToken(user.UserID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "token generation error: %v", err)
	}

	return &pb.OAuthAuthenticateResponse{Token: token}, nil
}
