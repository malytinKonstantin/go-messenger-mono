package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/malytinKonstantin/go-messenger-mono/friendship-service/internal/models"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var validStatuses = map[string]bool{
	"pending":  true,
	"accepted": true,
	"rejected": true,
}

type FriendRequestRepository interface {
	CreateFriendRequest(ctx context.Context, request *models.FriendRequest) error
	GetFriendRequestByID(ctx context.Context, requestID string) (*models.FriendRequest, error)
	UpdateFriendRequestStatus(ctx context.Context, requestID, status string, updatedAt int64) error
	DeleteFriendRequest(ctx context.Context, requestID string) error
	GetIncomingRequests(ctx context.Context, userID string) ([]*models.FriendRequest, error)
	GetOutgoingRequests(ctx context.Context, userID string) ([]*models.FriendRequest, error)
	DeleteFriendship(ctx context.Context, userID, friendID string) error
	GetIncomingAndOutgoingRequests(ctx context.Context, userID string) ([]*models.FriendRequest, error)
}

type friendRequestRepository struct {
	driver neo4j.DriverWithContext
}

func NewFriendRequestRepository(driver neo4j.DriverWithContext) FriendRequestRepository {
	return &friendRequestRepository{
		driver: driver,
	}
}

func (r *friendRequestRepository) CreateFriendRequest(ctx context.Context, request *models.FriendRequest) error {
	if !validStatuses[request.Status] {
		return fmt.Errorf("invalid status: %s", request.Status)
	}

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	// Устанавливаем текущие временные метки
	now := time.Now().Unix()
	request.CreatedAt = now
	request.UpdatedAt = now

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
            MATCH (sender:User {user_id: $senderID}), (receiver:User {user_id: $receiverID})
            CREATE (sender)-[r:FRIEND_REQUEST {
                request_id: $requestID,
                status: $status,
                created_at: $createdAt,
                updated_at: $updatedAt
            }]->(receiver)
        `
		params := map[string]interface{}{
			"requestID":  request.RequestID,
			"senderID":   request.SenderID,
			"receiverID": request.ReceiverID,
			"status":     request.Status,
			"createdAt":  request.CreatedAt,
			"updatedAt":  request.UpdatedAt,
		}
		_, err := tx.Run(ctx, query, params)
		return nil, err
	})
	return err
}

func (r *friendRequestRepository) GetFriendRequestByID(ctx context.Context, requestID string) (*models.FriendRequest, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
            MATCH (sender:User)-[r:FRIEND_REQUEST {request_id: $requestID}]->(receiver:User)
            RETURN r.request_id AS requestID, sender.user_id AS senderID, receiver.user_id AS receiverID,
                   r.status AS status, r.created_at AS createdAt, r.updated_at AS updatedAt
        `
		params := map[string]interface{}{
			"requestID": requestID,
		}
		record, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}
		if record.Next(ctx) {
			values := record.Record().Values
			return &models.FriendRequest{
				RequestID:  values[0].(string),
				SenderID:   values[1].(string),
				ReceiverID: values[2].(string),
				Status:     values[3].(string),
				CreatedAt:  values[4].(int64),
				UpdatedAt:  values[5].(int64),
			}, nil
		}
		return nil, errors.New("no results")
	})
	if err != nil {
		return nil, err
	}
	return result.(*models.FriendRequest), nil
}

func (r *friendRequestRepository) UpdateFriendRequestStatus(ctx context.Context, requestID, status string, updatedAt int64) error {
	if !validStatuses[status] {
		return fmt.Errorf("invalid status: %s", status)
	}

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
            MATCH (sender:User)-[r:FRIEND_REQUEST {request_id: $requestID}]->(receiver:User)
            SET r.status = $status, r.updated_at = $updatedAt
            RETURN r
        `
		params := map[string]interface{}{
			"requestID": requestID,
			"status":    status,
			"updatedAt": updatedAt,
		}
		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})
	return err
}

func (r *friendRequestRepository) DeleteFriendRequest(ctx context.Context, requestID string) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
            MATCH (sender:User)-[r:FRIEND_REQUEST {request_id: $requestID}]->(receiver:User)
            DELETE r
        `
		params := map[string]interface{}{
			"requestID": requestID,
		}
		_, err := tx.Run(ctx, query, params)
		return nil, err
	})
	return err
}

func (r *friendRequestRepository) GetIncomingRequests(ctx context.Context, userID string) ([]*models.FriendRequest, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
            MATCH (sender:User)-[r:FRIEND_REQUEST]->(receiver:User {user_id: $userID})
            RETURN r.request_id AS requestID, sender.user_id AS senderID, receiver.user_id AS receiverID,
                   r.status AS status, r.created_at AS createdAt, r.updated_at AS updatedAt
        `
		params := map[string]interface{}{
			"userID": userID,
		}
		records, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		var requests []*models.FriendRequest
		for records.Next(ctx) {
			values := records.Record().Values
			request := &models.FriendRequest{
				RequestID:  values[0].(string),
				SenderID:   values[1].(string),
				ReceiverID: values[2].(string),
				Status:     values[3].(string),
				CreatedAt:  values[4].(int64),
				UpdatedAt:  values[5].(int64),
			}
			requests = append(requests, request)
		}
		return requests, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*models.FriendRequest), nil
}

func (r *friendRequestRepository) GetOutgoingRequests(ctx context.Context, userID string) ([]*models.FriendRequest, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
            MATCH (sender:User {user_id: $userID})-[r:FRIEND_REQUEST]->(receiver:User)
            RETURN r.request_id AS requestID, sender.user_id AS senderID, receiver.user_id AS receiverID,
                   r.status AS status, r.created_at AS createdAt, r.updated_at AS updatedAt
        `
		params := map[string]interface{}{
			"userID": userID,
		}
		records, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		var requests []*models.FriendRequest
		for records.Next(ctx) {
			values := records.Record().Values
			request := &models.FriendRequest{
				RequestID:  values[0].(string),
				SenderID:   values[1].(string),
				ReceiverID: values[2].(string),
				Status:     values[3].(string),
				CreatedAt:  values[4].(int64),
				UpdatedAt:  values[5].(int64),
			}
			requests = append(requests, request)
		}
		return requests, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*models.FriendRequest), nil
}

func (r *friendRequestRepository) DeleteFriendship(ctx context.Context, userID, friendID string) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
            MATCH (u1:User {user_id: $userID})-[r:FRIEND_REQUEST {status: 'accepted'}]-(u2:User {user_id: $friendID})
            DELETE r
        `
		params := map[string]interface{}{
			"userID":   userID,
			"friendID": friendID,
		}
		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})
	return err
}

func (r *friendRequestRepository) GetIncomingAndOutgoingRequests(ctx context.Context, userID string) ([]*models.FriendRequest, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	requests := []*models.FriendRequest{}

	_, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
            MATCH (user:User {user_id: $userID})
            OPTIONAL MATCH (user)-[sent:FRIEND_REQUEST]->(receiver:User)
            OPTIONAL MATCH (sender:User)-[received:FRIEND_REQUEST]->(user)
            RETURN sent, receiver, received, sender
        `
		params := map[string]interface{}{
			"userID": userID,
		}
		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}
		for result.Next(ctx) {
			record := result.Record()

			if sent, ok := record.Get("sent"); ok && sent != nil {
				rel := sent.(neo4j.Relationship)
				request := &models.FriendRequest{
					RequestID:  rel.Props["request_id"].(string),
					SenderID:   userID,
					ReceiverID: record.Values[1].(neo4j.Node).Props["user_id"].(string),
					Status:     rel.Props["status"].(string),
					CreatedAt:  rel.Props["created_at"].(int64),
					UpdatedAt:  rel.Props["updated_at"].(int64),
				}
				requests = append(requests, request)
			}

			if received, ok := record.Get("received"); ok && received != nil {
				rel := received.(neo4j.Relationship)
				request := &models.FriendRequest{
					RequestID:  rel.Props["request_id"].(string),
					SenderID:   record.Values[3].(neo4j.Node).Props["user_id"].(string),
					ReceiverID: userID,
					Status:     rel.Props["status"].(string),
					CreatedAt:  rel.Props["created_at"].(int64),
					UpdatedAt:  rel.Props["updated_at"].(int64),
				}
				requests = append(requests, request)
			}
		}
		return nil, result.Err()
	})
	return requests, err
}
