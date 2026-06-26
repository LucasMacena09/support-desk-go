package models

import "time"

type Message struct {
	ID        string    `json:"id"`
	TicketID  string    `json:"ticket_id"`
	Sender    string    `json:"sender"` // "user" ou "bot"
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}