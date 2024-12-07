package grpc

import (
	"context"

	"github.com/gocql/gocql"
	pb "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/user_service/v1"
	"github.com/malytinKonstantin/go-messenger-mono/user-service/internal/models"
	"github.com/malytinKonstantin/go-messenger-mono/user-service/internal/usecases/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	getUserUsecase           user.GetUserUsecase
	createUserProfileUsecase user.CreateUserProfileUsecase
	updateUserProfileUsecase user.UpdateUserProfileUsecase
	searchUsersUsecase       user.SearchUsersUsecase
}

func NewUserHandler(
	getUserUC user.GetUserUsecase,
	createUserUC user.CreateUserProfileUsecase,
	updateUserUC user.UpdateUserProfileUsecase,
	searchUsersUC user.SearchUsersUsecase,
) *UserHandler {
	return &UserHandler{
		getUserUsecase:           getUserUC,
		createUserProfileUsecase: createUserUC,
		updateUserProfileUsecase: updateUserUC,
		searchUsersUsecase:       searchUsersUC,
	}
}

func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	userID, err := gocql.ParseUUID(req.UserId)
	if err != nil {
		return nil, err
	}

	profile, err := h.getUserUsecase.Execute(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserResponse{
		Profile: mapProfileToProto(profile),
	}, nil
}

func (h *UserHandler) CreateUserProfile(ctx context.Context, req *pb.CreateUserProfileRequest) (*pb.CreateUserProfileResponse, error) {
	profile := &models.UserProfile{
		Nickname:  req.Nickname,
		Bio:       req.Bio,
		AvatarURL: req.AvatarUrl,
	}

	err := h.createUserProfileUsecase.Execute(ctx, profile)
	if err != nil {
		if err == user.ErrNicknameAlreadyExists {
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}
		return nil, err
	}

	return &pb.CreateUserProfileResponse{
		Profile: mapProfileToProto(profile),
	}, nil
}

func (h *UserHandler) UpdateUserProfile(ctx context.Context, req *pb.UpdateUserProfileRequest) (*pb.UpdateUserProfileResponse, error) {
	userID, err := gocql.ParseUUID(req.UserId)
	if err != nil {
		return nil, err
	}

	profile := &models.UserProfile{
		UserID:    userID,
		Nickname:  req.Nickname,
		Bio:       req.Bio,
		AvatarURL: req.AvatarUrl,
	}

	err = h.updateUserProfileUsecase.Execute(ctx, profile)
	if err != nil {
		if err == user.ErrNicknameAlreadyExists {
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}
		return nil, err
	}

	return &pb.UpdateUserProfileResponse{
		Profile: mapProfileToProto(profile),
	}, nil
}

func (h *UserHandler) SearchUsers(ctx context.Context, req *pb.SearchUsersRequest) (*pb.SearchUsersResponse, error) {
	profiles, err := h.searchUsersUsecase.Execute(ctx, req.Query, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, err
	}

	var protoProfiles []*pb.UserProfile
	for _, profile := range profiles {
		protoProfiles = append(protoProfiles, mapProfileToProto(profile))
	}

	return &pb.SearchUsersResponse{
		Users: protoProfiles,
	}, nil
}

func mapProfileToProto(profile *models.UserProfile) *pb.UserProfile {
	return &pb.UserProfile{
		UserId:    profile.UserID.String(),
		Nickname:  profile.Nickname,
		Bio:       profile.Bio,
		AvatarUrl: profile.AvatarURL,
		CreatedAt: profile.CreatedAt.Unix(),
		UpdatedAt: profile.UpdatedAt.Unix(),
	}
}
