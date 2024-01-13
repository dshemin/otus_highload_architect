package token

import (
	"context"
	"github.com/dshemin/otus_highload_architect/internal/user"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Service struct {
	user       *user.Service
	repo       repository
	pwdChecker passwordChecker
}

type repository interface {
	Save(ctx context.Context, token, userID string) error
	IsExists(ctx context.Context, token string) (bool, error)
}

type passwordChecker interface {
	Check(given, hashed string) error
}

func NewService(
	user *user.Service,
	repo repository,
	pwdChecker passwordChecker,
) *Service {
	return &Service{
		user:       user,
		repo:       repo,
		pwdChecker: pwdChecker,
	}
}

func (s *Service) Authenticate(ctx context.Context, userID, password string) (token string, err error) {
	u, err := s.user.GetByID(ctx, userID)
	if err != nil {
		return "", err
	}

	err = s.pwdChecker.Check(password, u.HashedPassword)
	if err != nil {
		return "", errors.Wrap(err, "compare passwords")
	}

	token = uuid.New().String()

	err = s.repo.Save(ctx, token, userID)
	if err != nil {
		return "", errors.Wrap(err, "save token")
	}

	return token, nil
}

func (s *Service) IsExists(ctx context.Context, token string) (bool, error) {
	return s.repo.IsExists(ctx, token)
}
