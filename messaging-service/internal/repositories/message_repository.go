package repositories

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/malytinKonstantin/go-messenger-mono/messaging-service/internal/models"
)

type MessageRepository interface {
	SaveMessage(ctx context.Context, message *models.Message) error
	GetMessages(ctx context.Context, conversationID gocql.UUID, limit int) ([]*models.Message, error)
	UpdateMessageStatus(ctx context.Context, messageID gocql.UUID, status models.MessageStatus) error
}

type messageRepository struct {
	session *gocql.Session
}

func NewMessageRepository(session *gocql.Session) MessageRepository {
	return &messageRepository{
		session: session,
	}
}

func (r *messageRepository) SaveMessage(ctx context.Context, message *models.Message) error {
	query := `INSERT INTO messages (
        message_id, sender_id, recipient_id, conversation_id, content, status, timestamp
    ) VALUES (?, ?, ?, ?, ?, ?, ?)`
	return r.session.Query(query,
		message.MessageID,
		message.SenderID,
		message.RecipientID,
		message.ConversationID,
		message.Content,
		int32(message.Status),
		message.Timestamp,
	).WithContext(ctx).Exec()
}

func (r *messageRepository) GetMessages(ctx context.Context, conversationID gocql.UUID, limit int) ([]*models.Message, error) {
	query := `SELECT message_id, sender_id, recipient_id, content, status, timestamp FROM messages WHERE conversation_id = ? LIMIT ?`
	iter := r.session.Query(query, conversationID, limit).WithContext(ctx).Iter()

	var messages []*models.Message
	var msg models.Message
	for iter.Scan(&msg.MessageID, &msg.SenderID, &msg.RecipientID, &msg.Content, &msg.Status, &msg.Timestamp) {
		messages = append(messages, &msg)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *messageRepository) UpdateMessageStatus(ctx context.Context, messageID gocql.UUID, status models.MessageStatus) error {
	query := `UPDATE messages SET status = ? WHERE message_id = ?`
	return r.session.Query(query, int32(status), messageID).WithContext(ctx).Exec()
}
