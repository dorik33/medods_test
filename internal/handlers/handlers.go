package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dorik33/medods_test/internal/auth"
	"github.com/dorik33/medods_test/internal/models"
	"github.com/dorik33/medods_test/internal/store"
	"github.com/google/uuid"
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


func (h *Handlers) GenerateTokensHandler(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.URL.Query().Get("guid")
	if userIdStr == "" {
		sendErrorResponse(w, http.StatusBadRequest, "user_id is required")
		return
	}

	userID, err := uuid.Parse(userIdStr)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "invalid user_id format")
		return
	}

	ip := r.RemoteAddr
	accessToken, refreshToken, err := h.auth.GenerateTokens(userID, ip)
	if err != nil {
		h.logger.WithError(err).Error("failed to generate tokens")
		sendErrorResponse(w, http.StatusInternalServerError, "failed to generate tokens")
		return
	}

	response := models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.WithError(err).Error("failed to encode response")
	}
}

func (h *Handlers) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		h.logger.WithError(err).Error("failed to decode request body")
		sendErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	newAccessToken, newRefreshToken, err := h.auth.RefreshToken(
		requestBody.AccessToken,
		requestBody.RefreshToken,
		r.RemoteAddr,
	)

	if err != nil {
		h.logger.WithError(err).Error("failed to refresh tokens")
		sendErrorResponse(w, http.StatusBadRequest, "failed to refresh tokens")
		return
	}

	response := models.TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.WithError(err).Error("failed to encode response")
	}
}
