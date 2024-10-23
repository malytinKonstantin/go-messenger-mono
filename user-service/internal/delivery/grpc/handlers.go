package handlers

import (
	"context"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	pb "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/user_service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	session  *gocql.Session
	producer *kafka.Producer
}

func NewUserHandler(session *gocql.Session, producer *kafka.Producer) *UserHandler {
	return &UserHandler{
		session:  session,
		producer: producer,
	}
}

func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation error: %v", err)
	}

	profile := &pb.UserProfile{
		UserId:    req.UserId,
		Nickname:  "john_doe",
		Bio:       "Just a simple user.",
		AvatarUrl: "https://example.com/avatar.jpg",
		CreatedAt: time.Now().Add(-24 * time.Hour).Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	return &pb.GetUserResponse{
		Profile: profile,
	}, nil
}

func (h *UserHandler) CreateUserProfile(ctx context.Context, req *pb.CreateUserProfileRequest) (*pb.CreateUserProfileResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation error: %v", err)
	}

	userId := uuid.New().String()

	profile := &pb.UserProfile{
		UserId:    userId,
		Nickname:  req.Nickname,
		Bio:       req.Bio,
		AvatarUrl: req.AvatarUrl,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	return &pb.CreateUserProfileResponse{
		Profile: profile,
	}, nil
}

func (h *UserHandler) UpdateUserProfile(ctx context.Context, req *pb.UpdateUserProfileRequest) (*pb.UpdateUserProfileResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation error: %v", err)
	}

	profile := &pb.UserProfile{
		UserId:    req.UserId,
		Nickname:  req.Nickname,
		Bio:       req.Bio,
		AvatarUrl: req.AvatarUrl,
		CreatedAt: time.Now().Add(-48 * time.Hour).Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	return &pb.UpdateUserProfileResponse{
		Profile: profile,
	}, nil
}

func (h *UserHandler) SearchUsers(ctx context.Context, req *pb.SearchUsersRequest) (*pb.SearchUsersResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation error: %v", err)
	}

	users := []*pb.UserProfile{
		{
			UserId:    uuid.New().String(),
			Nickname:  "alice_smith",
			Bio:       "Love to travel.",
			AvatarUrl: "https://example.com/avatar1.jpg",
			CreatedAt: time.Now().Add(-72 * time.Hour).Unix(),
			UpdatedAt: time.Now().Add(-24 * time.Hour).Unix(),
		},
		{
			UserId:    uuid.New().String(),
			Nickname:  "bob_jones",
			Bio:       "Photographer and blogger.",
			AvatarUrl: "https://example.com/avatar2.jpg",
			CreatedAt: time.Now().Add(-96 * time.Hour).Unix(),
			UpdatedAt: time.Now().Add(-48 * time.Hour).Unix(),
		},
	}

	var start, end int32
	totalUsers := int32(len(users))

	if req.Offset > totalUsers {
		return &pb.SearchUsersResponse{
			Users: []*pb.UserProfile{},
		}, nil
	}

	start = req.Offset
	if req.Limit == 0 || req.Offset+req.Limit > totalUsers {
		end = totalUsers
	} else {
		end = req.Offset + req.Limit
	}

	paginatedUsers := users[start:end]

	return &pb.SearchUsersResponse{
		Users: paginatedUsers,
	}, nil
}
