package handlers

import (
	"net/http"

	"github.com/dorik33/medods_test/internal/auth"
	"github.com/dorik33/medods_test/internal/store"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Handlers struct {
	logger *logrus.Logger
	store  *store.Store
	auth   *auth.Service
}

func NewHandlers(logger *logrus.Logger, store *store.Store, authService *auth.Service) *Handlers {
	return &Handlers{
		logger: logger,
		store:  store,
		auth:   authService,
	}
}

type Message struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	IsError    bool   `json:"is_error"`
}

func (h *Handlers) GenerateTokenHandler(w http.ResponseWriter, r *http.Request) {
	userIdStr := mux.Vars(r)["guid"]
	if userIdStr == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		http.Error(w, "invalid user_id format", http.StatusBadRequest)
		return
	}

}
