package store

import (
	"context"
	"fmt"

	"github.com/dorik33/medods_test/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type Store struct {
	pool           *pgxpool.Pool
	logger         *logrus.Logger
	SongRepository *RefreshTokenRepository
}

func NewConnection(cfg *config.Config, logger *logrus.Logger) (*Store, error) {
	connStr := cfg.DatabaseURL
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	store := &Store{
		pool:   pool,
		logger: logger,
	}

	store.SongRepository = &RefreshTokenRepository{store: store}

	return store, nil
}

func (s *Store) Close() {
	s.pool.Close()
}
