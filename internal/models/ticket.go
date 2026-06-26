package models

import "time"

type Ticket struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Subject   string    `json:"subject"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}