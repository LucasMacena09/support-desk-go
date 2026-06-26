package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"support-desk-go/internal/models"
)

type MessageRepository struct {
	rdb *redis.Client
}

func NewMessageRepository(rdb *redis.Client) *MessageRepository {
	return &MessageRepository{rdb: rdb}
}

func ticketMessagesKey(ticketID string) string {
	return fmt.Sprintf("ticket:%s:messages", ticketID)
}

func (r *MessageRepository) Create(ctx context.Context, msg *models.Message) error {
	msg.ID = uuid.NewString()

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return r.rdb.RPush(ctx, ticketMessagesKey(msg.TicketID), data).Err()
}

func (r *MessageRepository) ListByTicket(ctx context.Context, ticketID string) ([]models.Message, error) {
	results, err := r.rdb.LRange(ctx, ticketMessagesKey(ticketID), 0, -1).Result()
	if err != nil {
		return nil, err
	}

	messages := make([]models.Message, 0, len(results))
	for _, raw := range results {
		var m models.Message
		if err := json.Unmarshal([]byte(raw), &m); err != nil {
			continue
		}
		messages = append(messages, m)
	}

	return messages, nil
}