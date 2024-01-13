package application

import (
	"context"
	"fmt"
	"github.com/dshemin/otus_highload_architect/internal/config"
	"github.com/dshemin/otus_highload_architect/internal/passwd"
	"github.com/dshemin/otus_highload_architect/internal/postgresql"
	"github.com/dshemin/otus_highload_architect/internal/token"
	"github.com/dshemin/otus_highload_architect/internal/user"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type container struct {
	UserService  *user.Service
	TokenService *token.Service
}

func initContainer(ctx context.Context, cfg config.Config) (container, error) {
	db, err := pgxpool.New(ctx, collectConnString(cfg.PostgreSQL))
	if err != nil {
		return container{}, errors.Wrap(err, "connect to db")
	}

	userRepo := postgresql.NewUserRepository(db)
	tokenRepo := postgresql.NewTokenRepository(db)

	pwdHasher := passwd.NewBcryptHasher(12)

	userService := user.NewService(userRepo, pwdHasher)
	tokenService := token.NewService(userService, tokenRepo, pwdHasher)

	return container{
		UserService:  userService,
		TokenService: tokenService,
	}, nil
}

func collectConnString(cfg config.PostgreSQL) string {
	sslMode := "disable"
	if cfg.SSLEnabled {
		sslMode = "require"
	}

	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		sslMode,
	)
}
