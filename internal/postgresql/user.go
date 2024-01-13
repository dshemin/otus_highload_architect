package postgresql

import (
	"context"
	"github.com/dshemin/otus_highload_architect/internal/user"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Save(ctx context.Context, u *user.User) error {
	_, err := r.db.Exec(ctx, `
INSERT INTO "users" (
	"id",
	"password",
	"first_name",
	"second_name",
	"birthdate",
	"biography",
	"city"
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
`, u.ID, u.HashedPassword, u.FirstName, u.SecondName, u.Birthdate, u.Biography, u.City)
	return errors.Wrap(err, "insert")
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (user.User, error) {
	row := r.db.QueryRow(ctx, `
SELECT
	"password",
	"first_name",
	"second_name",
	"birthdate",
	"biography",
	"city"
FROM "users"
WHERE id = $1
`, id)

	u := user.User{
		ID: id,
	}
	err := row.Scan(
		&u.HashedPassword,
		&u.FirstName,
		&u.SecondName,
		&u.Birthdate,
		&u.Biography,
		&u.City,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return user.User{}, user.ErrNotFound
	}
	return u, errors.Wrap(err, "query and scan row")
}
