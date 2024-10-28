package repositories

import (
	"context"

	"github.com/malytinKonstantin/go-messenger-mono/friendship-service/internal/models"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, userID string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, userID string) error
	AddFriend(ctx context.Context, userID, friendID string) error
	RemoveFriend(ctx context.Context, userID, friendID string) error
	GetFriends(ctx context.Context, userID string) ([]*models.User, error)
}

type userRepository struct {
	driver neo4j.DriverWithContext
}

func NewUserRepository(driver neo4j.DriverWithContext) UserRepository {
	return &userRepository{
		driver: driver,
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
            CREATE (u:User {user_id: $userID, nickname: $nickname, avatar_url: $avatarURL})
        `
		params := map[string]interface{}{
			"userID":    user.UserID,
			"nickname":  user.Nickname,
			"avatarURL": user.AvatarURL,
		}
		_, err := tx.Run(ctx, query, params)
		return nil, err
	})
	return err
}

func (r *userRepository) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
            MATCH (u:User {user_id: $userID})
            RETURN u.user_id AS userID, u.nickname AS nickname, u.avatar_url AS avatarURL
        `
		params := map[string]interface{}{
			"userID": userID,
		}
		record, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}
		if record.Next(ctx) {
			values := record.Record().Values
			return &models.User{
				UserID:    values[0].(string),
				Nickname:  values[1].(string),
				AvatarURL: values[2].(string),
			}, nil
		}
		return nil, neo4j.ErrNoResults
	})
	if err != nil {
		return nil, err
	}
	return result.(*models.User), nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *models.User) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
            MATCH (u:User {user_id: $userID})
            SET u.nickname = $nickname, u.avatar_url = $avatarURL
        `
		params := map[string]interface{}{
			"userID":    user.UserID,
			"nickname":  user.Nickname,
			"avatarURL": user.AvatarURL,
		}
		_, err := tx.Run(ctx, query, params)
		return nil, err
	})
	return err
}

func (r *userRepository) DeleteUser(ctx context.Context, userID string) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
            MATCH (u:User {user_id: $userID})
            DETACH DELETE u
        `
		params := map[string]interface{}{
			"userID": userID,
		}
		_, err := tx.Run(ctx, query, params)
		return nil, err
	})
	return err
}

func (r *userRepository) AddFriend(ctx context.Context, userID, friendID string) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
            MATCH (u1:User {user_id: $userID}), (u2:User {user_id: $friendID})
            MERGE (u1)-[:FRIEND_WITH]->(u2)
        `
		params := map[string]interface{}{
			"userID":   userID,
			"friendID": friendID,
		}
		_, err := tx.Run(ctx, query, params)
		return nil, err
	})
	return err
}

func (r *userRepository) RemoveFriend(ctx context.Context, userID, friendID string) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
            MATCH (u1:User {user_id: $userID})-[f:FRIEND_WITH]-(u2:User {user_id: $friendID})
            DELETE f
        `
		params := map[string]interface{}{
			"userID":   userID,
			"friendID": friendID,
		}
		_, err := tx.Run(ctx, query, params)
		return nil, err
	})
	return err
}

func (r *userRepository) GetFriends(ctx context.Context, userID string) ([]*models.User, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
            MATCH (u:User {user_id: $userID})-[:FRIEND_WITH]-(friend:User)
            RETURN friend.user_id AS userID, friend.nickname AS nickname, friend.avatar_url AS avatarURL
        `
		params := map[string]interface{}{
			"userID": userID,
		}
		records, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		var friends []*models.User
		for records.Next(ctx) {
			values := records.Record().Values
			friend := &models.User{
				UserID:    values[0].(string),
				Nickname:  values[1].(string),
				AvatarURL: values[2].(string),
			}
			friends = append(friends, friend)
		}
		return friends, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*models.User), nil
}
