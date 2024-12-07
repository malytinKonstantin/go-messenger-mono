package grpc

import (
	"context"

	"github.com/google/uuid"
	passwordUsecase "github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/usecase/password"
	pb "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/auth_service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PasswordHandler struct {
	pb.UnimplementedAuthServiceServer
	resetPasswordRequestUsecase passwordUsecase.ResetPasswordRequestUsecase
	changePasswordUsecase       passwordUsecase.ChangePasswordUsecase
}

func NewPasswordHandler(resetUC passwordUsecase.ResetPasswordRequestUsecase, changeUC passwordUsecase.ChangePasswordUsecase) *PasswordHandler {
	return &PasswordHandler{
		resetPasswordRequestUsecase: resetUC,
		changePasswordUsecase:       changeUC,
	}
}

func (h *PasswordHandler) ResetPasswordRequest(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	_, err := h.resetPasswordRequestUsecase.Execute(ctx, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "password reset error: %v", err)
	}

	return &pb.ResetPasswordResponse{Success: true}, nil
}

func (h *PasswordHandler) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	token, err := uuid.Parse(req.Token)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid token: %v", err)
	}

	err = h.changePasswordUsecase.Execute(ctx, token, req.NewPassword)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "password change error: %v", err)
	}

	return &pb.ChangePasswordResponse{Success: true}, nil
}
