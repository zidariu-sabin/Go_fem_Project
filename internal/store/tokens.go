package store

import (
	"database/sql"
	"time"

	"github.com/zidariu-sabin/femProject/internal/tokens"
)

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{
		db: db,
	}
}

type TokenStore interface {
	Insert(token *tokens.Token) error
	CreateNewToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error)
	DeleteAllTokensForUser(userID int, scope string) error
}

func (th *PostgresTokenStore) CreateNewToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateToken(userID, ttl, scope)

	if err != nil {
		return nil, err
	}

	err = th.Insert(token)

	return token, err
}

func (th *PostgresTokenStore) Insert(token *tokens.Token) error {
	query := `
	INSERT INTO tokens (hash, user_id, expiry, scope)
	VALUES ($1, $2, $3, $4)`

	_, err := th.db.Exec(query, token.Hash, token.UserID, token.Expiry, token.Scope)

	return err
}

func (th *PostgresTokenStore) DeleteAllTokensForUser(userID int, scope string) error {
	query := `
	DELETE FROM tokens
	WHERE user_id = $1 AND scope = $2`

	_, err := th.db.Exec(query, userID, scope)

	return err
}
