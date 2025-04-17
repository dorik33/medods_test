package auth

import (
	"time"

	"github.com/dorik33/medods_test/internal/models"
	"github.com/dorik33/medods_test/internal/store"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Service struct {
	store           *store.Store
	jwtSecret       []byte
	accessDuration  time.Duration
	refreshDuration time.Duration
}

func NewService(store *store.Store, jwtSecret string, accessTTL, refreshTTL time.Duration) *Service {
	return &Service{
		store:           store,
		jwtSecret:       []byte(jwtSecret),
		accessDuration:  accessTTL,
		refreshDuration: refreshTTL,
	}
}

func (s *Service) generateAccessToken(userID uuid.UUID, ip string) (string, error) {
	claims := &models.Claims{
		UserID:    userID,
		IPAddress: ip,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(s.jwtSecret)
}
