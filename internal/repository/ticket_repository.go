package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"support-desk-go/internal/models"
)

var ErrTicketNotFound = errors.New("chamado não encontrado")

type TicketRepository struct {
	rdb *redis.Client
}

func NewTicketRepository(rdb *redis.Client) *TicketRepository {
	return &TicketRepository{rdb: rdb}
}

func ticketKey(id string) string {
	return fmt.Sprintf("ticket:%s", id)
}

func userTicketsKey(userID string) string {
	return fmt.Sprintf("user:%s:tickets", userID)
}

func (r *TicketRepository) Create(ctx context.Context, ticket *models.Ticket) error {
	ticket.ID = uuid.NewString()

	data, err := json.Marshal(ticket)
	if err != nil {
		return err
	}

	pipe := r.rdb.TxPipeline()
	pipe.Set(ctx, ticketKey(ticket.ID), data, 0)
	pipe.SAdd(ctx, userTicketsKey(ticket.UserID), ticket.ID)
	_, err = pipe.Exec(ctx)
	return err
}

func (r *TicketRepository) FindByID(ctx context.Context, id string) (*models.Ticket, error) {
	data, err := r.rdb.Get(ctx, ticketKey(id)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrTicketNotFound
		}
		return nil, err
	}

	var ticket models.Ticket
	if err := json.Unmarshal(data, &ticket); err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (r *TicketRepository) ListByUser(ctx context.Context, userID string) ([]models.Ticket, error) {
	ids, err := r.rdb.SMembers(ctx, userTicketsKey(userID)).Result()
	if err != nil {
		return nil, err
	}

	tickets := make([]models.Ticket, 0, len(ids))
	for _, id := range ids {
		t, err := r.FindByID(ctx, id)
		if err != nil {
			continue
		}
		tickets = append(tickets, *t)
	}

	return tickets, nil
}