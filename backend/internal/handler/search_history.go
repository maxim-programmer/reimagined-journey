package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/maxim-programmer/reimagined-journey/backend/internal/middleware"
	"github.com/maxim-programmer/reimagined-journey/backend/internal/model"
)

type searchHistoryService interface {
	ListHistory(ctx context.Context, userID string) ([]model.SearchHistory, error)
	ClearHistory(ctx context.Context, userID string) error
}

type SearchHistoryHandler struct {
	svc searchHistoryService
}

func NewSearchHistoryHandler(svc searchHistoryService) *SearchHistoryHandler {
	return &SearchHistoryHandler{svc: svc}
}

func (h *SearchHistoryHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	items, err := h.svc.ListHistory(r.Context(), userID)
	if err != nil {
		log.Printf("list history error: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to retrieve history")
		return
	}

	if items == nil {
		items = []model.SearchHistory{}
	}

	writeJSON(w, http.StatusOK, items)
}

func (h *SearchHistoryHandler) Clear(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.svc.ClearHistory(r.Context(), userID); err != nil {
		log.Printf("clear history error: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to clear history")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}