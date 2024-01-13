package postgresql

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type TokenRepository struct {
	db *pgxpool.Pool
}

func NewTokenRepository(db *pgxpool.Pool) *TokenRepository {
	return &TokenRepository{
		db: db,
	}
}

func (r *TokenRepository) Save(ctx context.Context, token, userID string) error {
	_, err := r.db.Exec(ctx, `INSERT INTO "tokens" ("token", "user_id") VALUES ($1, $2)`, token, userID)
	return errors.Wrap(err, "insert")
}

func (r *TokenRepository) IsExists(ctx context.Context, token string) (bool, error) {
	row := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM "tokens" WHERE "token" = $1`, token)

	var cnt int
	err := row.Scan(&cnt)
	return cnt > 0, errors.Wrap(err, "query and scan row")
}
