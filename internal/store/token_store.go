package store

import (
	"database/sql"
	"time"

	"github.com/andras-szesztai/fem_fitness_project/internal/tokens"
)

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{db: db}
}

type TokenStore interface {
	InsertToken(token *tokens.Token) error
	CreateToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error)
	DeleteToken(userID int, scope string) error
}

func (s *PostgresTokenStore) InsertToken(token *tokens.Token) error {
	query := `
	INSERT INTO tokens (hash, user_id, expiry, scope)
	VALUES ($1, $2, $3, $4)
	`

	_, err := s.db.Exec(query, token.Hash, token.UserID, token.Expiry, token.Scope)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresTokenStore) CreateToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = s.InsertToken(token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *PostgresTokenStore) DeleteToken(userID int, scope string) error {
	query := `
	DELETE FROM tokens
	WHERE user_id = $1 AND scope = $2
	`

	_, err := s.db.Exec(query, userID, scope)
	if err != nil {
		return err
	}
	return nil
}
