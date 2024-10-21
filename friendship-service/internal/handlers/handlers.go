package handlers

import (
	"context"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	pb "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/friendship_service/v1"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FriendshipHandler struct {
	pb.UnimplementedFriendshipServiceServer
	producer *kafka.Producer
	driver   neo4j.Driver
}

func NewFriendshipHandler(producer *kafka.Producer, driver neo4j.Driver) *FriendshipHandler {
	return &FriendshipHandler{producer: producer, driver: driver}
}

func (h *FriendshipHandler) SendFriendRequest(ctx context.Context, req *pb.SendFriendRequestRequest) (*pb.SendFriendRequestResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Ошибка валидации: %v", err)
	}
	// Моковый ответ с учетом типа pb.SendFriendRequestResponse
	return &pb.SendFriendRequestResponse{
		RequestId: "550e8400-e29b-41d4-a716-446655440000",
	}, nil
}

func (h *FriendshipHandler) AcceptFriendRequest(ctx context.Context, req *pb.AcceptFriendRequestRequest) (*pb.AcceptFriendRequestResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Ошибка валидации: %v", err)
	}
	// Моковый ответ с учетом типа pb.AcceptFriendRequestResponse
	return &pb.AcceptFriendRequestResponse{
		Success: true,
	}, nil
}

func (h *FriendshipHandler) RejectFriendRequest(ctx context.Context, req *pb.RejectFriendRequestRequest) (*pb.RejectFriendRequestResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Ошибка валидации: %v", err)
	}
	// Моковый ответ с учетом типа pb.RejectFriendRequestResponse
	return &pb.RejectFriendRequestResponse{
		Success: true,
	}, nil
}

func (h *FriendshipHandler) RemoveFriend(ctx context.Context, req *pb.RemoveFriendRequest) (*pb.RemoveFriendResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Ошибка валидации: %v", err)
	}
	// Моковый ответ
	return &pb.RemoveFriendResponse{}, nil
}

func (h *FriendshipHandler) GetFriendsList(ctx context.Context, req *pb.GetFriendsListRequest) (*pb.GetFriendsListResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Ошибка валидации: %v", err)
	}
	// Моковый список друзей с учетом типа pb.Friend
	friends := []*pb.Friend{
		{
			UserId:    "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			Nickname:  "Иван Иванов",
			AvatarUrl: "https://example.com/avatar1.png",
			AddedAt:   time.Now().Unix(),
		},
		{
			UserId:    "6ba7b811-9dad-11d1-80b4-00c04fd430c8",
			Nickname:  "Петр Петров",
			AvatarUrl: "https://example.com/avatar2.png",
			AddedAt:   time.Now().Unix(),
		},
	}
	return &pb.GetFriendsListResponse{Friends: friends}, nil
}

func (h *FriendshipHandler) GetPendingRequests(ctx context.Context, req *pb.GetPendingRequestsRequest) (*pb.GetPendingRequestsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Ошибка валидации: %v", err)
	}
	// Моковые входящие и исходящие запросы с учетом структуры FriendRequest
	incomingRequests := []*pb.FriendRequest{
		{
			RequestId:  "6ba7b812-9dad-11d1-80b4-00c04fd430c8",
			SenderId:   "6ba7b813-9dad-11d1-80b4-00c04fd430c8",
			ReceiverId: req.UserId,
			Status:     "pending",
			CreatedAt:  time.Now().Unix(),
			UpdatedAt:  time.Now().Unix(),
		},
	}
	outgoingRequests := []*pb.FriendRequest{
		{
			RequestId:  "6ba7b814-9dad-11d1-80b4-00c04fd430c8",
			SenderId:   req.UserId,
			ReceiverId: "6ba7b815-9dad-11d1-80b4-00c04fd430c8",
			Status:     "pending",
			CreatedAt:  time.Now().Unix(),
			UpdatedAt:  time.Now().Unix(),
		},
	}
	return &pb.GetPendingRequestsResponse{
		IncomingRequests: incomingRequests,
		OutgoingRequests: outgoingRequests,
	}, nil
}
