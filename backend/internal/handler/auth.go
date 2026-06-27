package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/maxim-programmer/reimagined-journey/backend/internal/middleware"
	"github.com/maxim-programmer/reimagined-journey/backend/internal/model"
)

type authService interface {
	Register(ctx context.Context, login, password string) (*model.User, error)
	Login(ctx context.Context, login, password string) (string, *model.User, error)
	Logout(ctx context.Context, token string) error
	GetUserByToken(ctx context.Context, token string) (*model.User, error)
}

type AuthHandler struct {
	svc authService
}

func NewAuthHandler(svc authService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

type authRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.Login = strings.TrimSpace(req.Login)
	if req.Login == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "login and password are required")
		return
	}
	if len(req.Password) < 6 {
		writeError(w, http.StatusBadRequest, "password must be at least 6 characters")
		return
	}

	user, err := h.svc.Register(r.Context(), req.Login, req.Password)
	if err != nil {
		if strings.Contains(err.Error(), "already taken") {
			writeError(w, http.StatusConflict, "login already taken")
			return
		}
		log.Printf("register error: %v", err)
		writeError(w, http.StatusInternalServerError, "registration failed")
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	token, user, err := h.svc.Login(r.Context(), req.Login, req.Password)
	if err != nil {
		if strings.Contains(err.Error(), "invalid credentials") {
			writeError(w, http.StatusUnauthorized, "invalid login or password")
			return
		}
		log.Printf("login error: %v", err)
		writeError(w, http.StatusInternalServerError, "login failed")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"token": token,
		"user":  user,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	bearer := r.Header.Get("Authorization")
	token := strings.TrimPrefix(bearer, "Bearer ")
	if token == "" {
		writeError(w, http.StatusBadRequest, "missing token")
		return
	}

	if err := h.svc.Logout(r.Context(), token); err != nil {
		log.Printf("logout error: %v", err)
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	bearer := r.Header.Get("Authorization")
	token := strings.TrimPrefix(bearer, "Bearer ")

	user, err := h.svc.GetUserByToken(r.Context(), token)
	if err != nil || user == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	writeJSON(w, http.StatusOK, user)
}