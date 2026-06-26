package main

import (
	"context"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"

	"support-desk-go/internal/auth"
	"support-desk-go/internal/chatbot"
	"support-desk-go/internal/config"
	"support-desk-go/internal/handlers"
	"support-desk-go/internal/repository"
)

func main() {
	cfg := config.Load()

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("erro ao conectar no Redis: %v", err)
	}
	log.Println("conectado ao Redis com sucesso")

	jwtManager := auth.NewJWTManager(cfg.JWTSecret)

	userRepo := repository.NewUserRepository(rdb)
	ticketRepo := repository.NewTicketRepository(rdb)
	messageRepo := repository.NewMessageRepository(rdb)

	bot := chatbot.NewBot()

	authHandler := handlers.NewAuthHandler(userRepo, jwtManager)
	ticketHandler := handlers.NewTicketHandler(ticketRepo)
	messageHandler := handlers.NewMessageHandler(messageRepo, ticketRepo, bot)
	pageHandler := handlers.NewPageHandler()

	mux := http.NewServeMux()

	// Arquivos estáticos (css/js)
	fs := http.FileServer(http.Dir("./web/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	// Páginas HTML (públicas)
	mux.HandleFunc("GET /", pageHandler.Home)
	mux.HandleFunc("GET /login", pageHandler.LoginPage)
	mux.HandleFunc("GET /register", pageHandler.RegisterPage)
	mux.HandleFunc("GET /chat", pageHandler.ChatPage)

	// API pública
	mux.HandleFunc("POST /api/register", authHandler.Register)
	mux.HandleFunc("POST /api/login", authHandler.Login)

	// API protegida (exige token)
	mux.Handle("POST /api/tickets", auth.Middleware(jwtManager, http.HandlerFunc(ticketHandler.Create)))
	mux.Handle("GET /api/tickets", auth.Middleware(jwtManager, http.HandlerFunc(ticketHandler.List)))
	mux.Handle("GET /api/tickets/{id}", auth.Middleware(jwtManager, http.HandlerFunc(ticketHandler.Get)))
	mux.Handle("POST /api/tickets/{id}/messages", auth.Middleware(jwtManager, http.HandlerFunc(messageHandler.Send)))
	mux.Handle("GET /api/tickets/{id}/messages", auth.Middleware(jwtManager, http.HandlerFunc(messageHandler.List)))

	log.Printf("servidor rodando na porta %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}