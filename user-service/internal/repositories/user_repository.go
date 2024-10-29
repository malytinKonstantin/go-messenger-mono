package repositories

import (
	"context"
	"time"

	"github.com/gocql/gocql"
	"github.com/malytinKonstantin/go-messenger-mono/user-service/internal/models"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
)

type UserRepository interface {
	GetUser(ctx context.Context, userID gocql.UUID) (*models.UserProfile, error)
	CreateUserProfile(ctx context.Context, profile *models.UserProfile) error
	UpdateUserProfile(ctx context.Context, profile *models.UserProfile) error
	SearchUsers(ctx context.Context, query string, limit, offset int) ([]*models.UserProfile, error)
	NicknameExists(ctx context.Context, nickname string) (bool, error)
}

type userRepository struct {
	session *gocql.Session
}

// NewUserRepository создаёт новый экземпляр UserRepository.
func NewUserRepository(session *gocql.Session) UserRepository {
	return &userRepository{
		session: session,
	}
}

// GetUser получает информацию о пользователе по его user_id.
func (r *userRepository) GetUser(ctx context.Context, userID gocql.UUID) (*models.UserProfile, error) {
	var profile models.UserProfile

	stmt, names := qb.Select("user_profiles").
		Where(qb.Eq("user_id")).
		ToCql()

	q := gocqlx.Query(r.session.Query(stmt), names).WithContext(ctx)
	err := q.BindMap(qb.M{"user_id": userID}).GetRelease(&profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

// CreateUserProfile создаёт новый профиль пользователя.
func (r *userRepository) CreateUserProfile(ctx context.Context, profile *models.UserProfile) error {
	profile.CreatedAt = time.Now()
	profile.UpdatedAt = time.Now()

	stmt, names := qb.Insert("user_profiles").
		Columns("user_id", "nickname", "bio", "avatar_url", "created_at", "updated_at").
		ToCql()

	q := gocqlx.Query(r.session.Query(stmt), names).WithContext(ctx)
	return q.BindStruct(profile).ExecRelease()
}

// UpdateUserProfile обновляет существующий профиль пользователя.
func (r *userRepository) UpdateUserProfile(ctx context.Context, profile *models.UserProfile) error {
	profile.UpdatedAt = time.Now()

	stmt, names := qb.Update("user_profiles").
		Set("nickname", "bio", "avatar_url", "updated_at").
		Where(qb.Eq("user_id")).
		ToCql()

	q := gocqlx.Query(r.session.Query(stmt), names).WithContext(ctx)
	return q.BindStruct(profile).ExecRelease()
}

// SearchUsers ищет пользователей по заданному запросу.
func (r *userRepository) SearchUsers(ctx context.Context, query string, limit, offset int) ([]*models.UserProfile, error) {
	var profiles []*models.UserProfile

	stmt, names := qb.Select("user_profiles").
		Where(qb.Like("nickname")).
		Limit(uint(limit)).
		ToCql()

	q := gocqlx.Query(r.session.Query(stmt), names).WithContext(ctx)
	err := q.BindMap(qb.M{"nickname": "%" + query + "%"}).SelectRelease(&profiles)
	if err != nil {
		return nil, err
	}

	return profiles, nil
}

// NicknameExists проверяет, существует ли пользователь с данным никнеймом
func (r *userRepository) NicknameExists(ctx context.Context, nickname string) (bool, error) {
	var count int
	stmt, names := qb.Select("user_profiles").
		Columns("count(*)").
		Where(qb.Eq("nickname")).
		Limit(1).
		ToCql()

	q := gocqlx.Query(r.session.Query(stmt), names).WithContext(ctx)
	err := q.BindMap(qb.M{"nickname": nickname}).GetRelease(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
