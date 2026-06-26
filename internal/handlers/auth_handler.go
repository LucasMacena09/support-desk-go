package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"support-desk-go/internal/auth"
	"support-desk-go/internal/httperr"
	"support-desk-go/internal/models"
	"support-desk-go/internal/repository"
)

type AuthHandler struct {
	userRepo   *repository.UserRepository
	jwtManager *auth.JWTManager
}

func NewAuthHandler(userRepo *repository.UserRepository, jwtManager *auth.JWTManager) *AuthHandler {
	return &AuthHandler{userRepo: userRepo, jwtManager: jwtManager}
}

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperr.BadRequest(w, "corpo da requisição inválido")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.Name == "" || req.Email == "" || req.Password == "" {
		httperr.BadRequest(w, "nome, email e senha são obrigatórios")
		return
	}

	hashed, err := auth.HashPassword(req.Password)
	if err != nil {
		httperr.Internal(w, "erro ao processar senha")
		return
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashed,
	}

	ctx := r.Context()
	if err := h.userRepo.Create(ctx, user); err != nil {
		if err == repository.ErrUserAlreadyExists {
			httperr.Conflict(w, "usuário já cadastrado com esse email")
			return
		}
		httperr.Internal(w, "erro de conexão com o banco de dados")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "usuário cadastrado com sucesso"})
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperr.BadRequest(w, "corpo da requisição inválido")
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.Email == "" || req.Password == "" {
		httperr.BadRequest(w, "email e senha são obrigatórios")
		return
	}

	ctx := r.Context()
	user, err := h.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if err == repository.ErrUserNotFound {
			httperr.Unauthorized(w, "usuário não cadastrado")
			return
		}
		httperr.Internal(w, "erro de conexão com o banco de dados")
		return
	}

	if !auth.CheckPassword(user.Password, req.Password) {
		httperr.Unauthorized(w, "senha incorreta")
		return
	}

	token, err := h.jwtManager.Generate(user.Email)
	if err != nil {
		httperr.Internal(w, "erro ao gerar token de autenticação")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}