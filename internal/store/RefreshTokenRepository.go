package store

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dorik33/medods_test/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type RefreshTokenRepository struct {
	store  *Store
	logger *logrus.Logger
}

func (r *RefreshTokenRepository) Create(ctx context.Context, refreshToken *models.RefreshToken) error {
	query := `
	 INSERT INTO refresh_tokens (user_id, access_token_jti, token_hash, expires_at, ip_address)
	 VALUES ($1, $2, $3, $4, $5)
	`
	r.store.logger.Debugf("Executing query: %s", query)

	_, err := r.store.pool.Exec(ctx, query,
		refreshToken.UserID,
		refreshToken.AccessTokenJTI,
		refreshToken.TokenHash,
		refreshToken.ExpiresAt,
		refreshToken.IPAddress,
	)
	log.Println(query, err)
	if err != nil {
		return fmt.Errorf("failed to insert refresh token: %w", err)
	}

	return nil
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, id int) error {
	query := `
	 DELETE FROM refresh_tokens WHERE id = $1
	`
	r.store.logger.Debugf("Executing query: %s with params: id=%d", query, id)

	_, err := r.store.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	return nil
}

func (r *RefreshTokenRepository) FindByJTI(ctx context.Context, jti uuid.UUID) (*models.RefreshToken, error) {
	query := `
        SELECT id, user_id, access_token_jti, token_hash, expires_at, ip_address
        FROM refresh_tokens
        WHERE access_token_jti = $1 AND expires_at > $2
    `
	r.store.logger.Debugf("Executing query: %s with params: jti=%s", query, jti)

	row := r.store.pool.QueryRow(ctx, query, jti, time.Now())

	refreshToken := &models.RefreshToken{}
	err := row.Scan(
		&refreshToken.ID,
		&refreshToken.UserID,
		&refreshToken.AccessTokenJTI,
		&refreshToken.TokenHash,
		&refreshToken.ExpiresAt,
		&refreshToken.IPAddress,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("refresh token not found")
		}
		return nil, fmt.Errorf("failed to query refresh token: %w", err)
	}

	return refreshToken, nil
}
