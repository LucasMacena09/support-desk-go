package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"support-desk-go/internal/auth"
	"support-desk-go/internal/httperr"
	"support-desk-go/internal/models"
	"support-desk-go/internal/repository"
)

type TicketHandler struct {
	ticketRepo *repository.TicketRepository
}

func NewTicketHandler(ticketRepo *repository.TicketRepository) *TicketHandler {
	return &TicketHandler{ticketRepo: ticketRepo}
}

type createTicketRequest struct {
	Subject string `json:"subject"`
}

func (h *TicketHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(auth.UserIDKey).(string)

	var req createTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperr.BadRequest(w, "corpo da requisição inválido")
		return
	}

	req.Subject = strings.TrimSpace(req.Subject)
	if req.Subject == "" {
		httperr.BadRequest(w, "assunto do chamado é obrigatório")
		return
	}

	ticket := &models.Ticket{
		UserID:    userID,
		Subject:   req.Subject,
		Status:    "open",
		CreatedAt: time.Now(),
	}

	ctx := r.Context()
	if err := h.ticketRepo.Create(ctx, ticket); err != nil {
		httperr.Internal(w, "erro de conexão com o banco de dados")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ticket)
}

func (h *TicketHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(auth.UserIDKey).(string)

	ctx := r.Context()
	tickets, err := h.ticketRepo.ListByUser(ctx, userID)
	if err != nil {
		httperr.Internal(w, "erro de conexão com o banco de dados")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tickets)
}

func (h *TicketHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ctx := r.Context()
	ticket, err := h.ticketRepo.FindByID(ctx, id)
	if err != nil {
		if err == repository.ErrTicketNotFound {
			httperr.NotFound(w, "chamado não encontrado")
			return
		}
		httperr.Internal(w, "erro de conexão com o banco de dados")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ticket)
}