package handlers

import (
	"html/template"
	"net/http"
)

type PageHandler struct{}

func NewPageHandler() *PageHandler {
	return &PageHandler{}
}

func (h *PageHandler) Home(w http.ResponseWriter, r *http.Request) {
	render(w, "web/templates/home.html")
}

func (h *PageHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	render(w, "web/templates/login.html")
}

func (h *PageHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	render(w, "web/templates/register.html")
}

func (h *PageHandler) ChatPage(w http.ResponseWriter, r *http.Request) {
	render(w, "web/templates/chat.html")
}

func render(w http.ResponseWriter, path string) {
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, "erro ao carregar página", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}