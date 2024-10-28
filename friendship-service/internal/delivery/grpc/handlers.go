package grpc

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/malytinKonstantin/go-messenger-mono/friendship-service/internal/usecase/friendship"
	pb "github.com/malytinKonstantin/go-messenger-mono/proto/pkg/api/friendship_service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FriendshipHandler struct {
	pb.UnimplementedFriendshipServiceServer
	producer              *kafka.Producer
	sendFriendRequestUC   friendship.SendFriendRequestUsecase
	acceptFriendRequestUC friendship.AcceptFriendRequestUsecase
	rejectFriendRequestUC friendship.RejectFriendRequestUsecase
	removeFriendUC        friendship.RemoveFriendUsecase
	getFriendsListUC      friendship.GetFriendsListUsecase
	getFriendRequestsUC   friendship.GetFriendRequestsUsecase
}

func NewFriendshipHandler(
	producer *kafka.Producer,
	sendFriendRequestUC friendship.SendFriendRequestUsecase,
	acceptFriendRequestUC friendship.AcceptFriendRequestUsecase,
	rejectFriendRequestUC friendship.RejectFriendRequestUsecase,
	removeFriendUC friendship.RemoveFriendUsecase,
	getFriendsListUC friendship.GetFriendsListUsecase,
	getFriendRequestsUC friendship.GetFriendRequestsUsecase,
) *FriendshipHandler {
	return &FriendshipHandler{
		producer:              producer,
		sendFriendRequestUC:   sendFriendRequestUC,
		acceptFriendRequestUC: acceptFriendRequestUC,
		rejectFriendRequestUC: rejectFriendRequestUC,
		removeFriendUC:        removeFriendUC,
		getFriendsListUC:      getFriendsListUC,
		getFriendRequestsUC:   getFriendRequestsUC,
	}
}

func (h *FriendshipHandler) SendFriendRequest(ctx context.Context, req *pb.SendFriendRequestRequest) (*pb.SendFriendRequestResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	request, err := h.sendFriendRequestUC.Execute(ctx, req.SenderId, req.ReceiverId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error sending friend request: %v", err)
	}

	// Отправка уведомления через Kafka (при необходимости)
	// ...

	return &pb.SendFriendRequestResponse{
		RequestId: request.RequestID,
	}, nil
}

func (h *FriendshipHandler) AcceptFriendRequest(ctx context.Context, req *pb.AcceptFriendRequestRequest) (*pb.AcceptFriendRequestResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	err := h.acceptFriendRequestUC.Execute(ctx, req.RequestId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error accepting friend request: %v", err)
	}

	// Отправка уведомления через Kafka (при необходимости)
	// ...

	return &pb.AcceptFriendRequestResponse{
		Success: true,
	}, nil
}

func (h *FriendshipHandler) RejectFriendRequest(ctx context.Context, req *pb.RejectFriendRequestRequest) (*pb.RejectFriendRequestResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	err := h.rejectFriendRequestUC.Execute(ctx, req.RequestId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error rejecting friend request: %v", err)
	}

	// Отправка уведомления через Kafka (при необходимости)
	// ...

	return &pb.RejectFriendRequestResponse{
		Success: true,
	}, nil
}

func (h *FriendshipHandler) RemoveFriend(ctx context.Context, req *pb.RemoveFriendRequest) (*pb.RemoveFriendResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	err := h.removeFriendUC.Execute(ctx, req.UserId, req.FriendId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error removing friend: %v", err)
	}

	// Отправка уведомления через Kafka (при необходимости)
	// ...

	return &pb.RemoveFriendResponse{
		Success: true,
	}, nil
}

func (h *FriendshipHandler) GetFriendsList(ctx context.Context, req *pb.GetFriendsListRequest) (*pb.GetFriendsListResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	friends, err := h.getFriendsListUC.Execute(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error getting friends list: %v", err)
	}

	var friendList []*pb.Friend
	for _, friend := range friends {
		friendList = append(friendList, &pb.Friend{
			UserId:    friend.UserID,
			Nickname:  friend.Nickname,
			AvatarUrl: friend.AvatarURL,
			AddedAt:   friend.AddedAt,
		})
	}

	return &pb.GetFriendsListResponse{
		Friends: friendList,
	}, nil
}

func (h *FriendshipHandler) GetPendingRequests(ctx context.Context, req *pb.GetPendingRequestsRequest) (*pb.GetPendingRequestsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	requests, err := h.getFriendRequestsUC.Execute(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error getting friend requests: %v", err)
	}

	var incomingRequests []*pb.FriendRequest
	var outgoingRequests []*pb.FriendRequest

	for _, request := range requests {
		pbRequest := &pb.FriendRequest{
			RequestId:  request.RequestID,
			SenderId:   request.SenderID,
			ReceiverId: request.ReceiverID,
			Status:     request.Status,
			CreatedAt:  request.CreatedAt,
			UpdatedAt:  request.UpdatedAt,
		}
		if request.ReceiverID == req.UserId {
			incomingRequests = append(incomingRequests, pbRequest)
		} else if request.SenderID == req.UserId {
			outgoingRequests = append(outgoingRequests, pbRequest)
		}
	}

	return &pb.GetPendingRequestsResponse{
		IncomingRequests: incomingRequests,
		OutgoingRequests: outgoingRequests,
	}, nil
}
