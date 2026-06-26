package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"support-desk-go/internal/chatbot"
	"support-desk-go/internal/httperr"
	"support-desk-go/internal/models"
	"support-desk-go/internal/repository"
)

type MessageHandler struct {
	messageRepo *repository.MessageRepository
	ticketRepo  *repository.TicketRepository
	bot         *chatbot.Bot
}

func NewMessageHandler(messageRepo *repository.MessageRepository, ticketRepo *repository.TicketRepository, bot *chatbot.Bot) *MessageHandler {
	return &MessageHandler{messageRepo: messageRepo, ticketRepo: ticketRepo, bot: bot}
}

type sendMessageRequest struct {
	Content string `json:"content"`
}

func (h *MessageHandler) Send(w http.ResponseWriter, r *http.Request) {
	ticketID := r.PathValue("id")

	ctx := r.Context()
	_, err := h.ticketRepo.FindByID(ctx, ticketID)
	if err != nil {
		if err == repository.ErrTicketNotFound {
			httperr.NotFound(w, "chamado não encontrado")
			return
		}
		httperr.Internal(w, "erro de conexão com o banco de dados")
		return
	}

	var req sendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperr.BadRequest(w, "corpo da requisição inválido")
		return
	}

	req.Content = strings.TrimSpace(req.Content)
	if req.Content == "" {
		httperr.BadRequest(w, "mensagem não pode estar em branco")
		return
	}

	userMsg := &models.Message{
		TicketID:  ticketID,
		Sender:    "user",
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	if err := h.messageRepo.Create(ctx, userMsg); err != nil {
		httperr.Internal(w, "erro de conexão com o banco de dados")
		return
	}

	botReply := h.bot.Respond(req.Content)
	botMsg := &models.Message{
		TicketID:  ticketID,
		Sender:    "bot",
		Content:   botReply,
		CreatedAt: time.Now(),
	}

	if err := h.messageRepo.Create(ctx, botMsg); err != nil {
		httperr.Internal(w, "erro de conexão com o banco de dados")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]*models.Message{
		"user_message": userMsg,
		"bot_message":  botMsg,
	})
}

func (h *MessageHandler) List(w http.ResponseWriter, r *http.Request) {
	ticketID := r.PathValue("id")

	ctx := r.Context()
	messages, err := h.messageRepo.ListByTicket(ctx, ticketID)
	if err != nil {
		httperr.Internal(w, "erro de conexão com o banco de dados")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}