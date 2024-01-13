package user

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
)

type Service struct {
	repo      repository
	pwdHasher passwordHasher
}

type repository interface {
	Save(ctx context.Context, u *User) error
	GetByID(ctx context.Context, id string) (User, error)
}

type passwordHasher interface {
	Hash(p string) (string, error)
}

func NewService(
	repo repository,
	pwdHasher passwordHasher,
) *Service {
	return &Service{
		repo:      repo,
		pwdHasher: pwdHasher,
	}
}

var ErrNotFound = errors.New("user not found")

func (s *Service) Register(ctx context.Context, dto RegisterDTO) (id string, err error) {
	birthdate, err := time.Parse(time.DateOnly, dto.Birthdate)
	if err != nil {
		return "", errors.Wrap(err, "parse birthdate")
	}

	hashedPassword, err := s.pwdHasher.Hash(dto.Password)
	if err != nil {
		return "", errors.Wrap(err, "hash password")
	}

	u := User{
		ID:             uuid.New().String(),
		FirstName:      dto.FirstName,
		SecondName:     dto.SecondName,
		Birthdate:      birthdate,
		Biography:      dto.Biography,
		City:           dto.City,
		HashedPassword: hashedPassword,
	}

	err = s.repo.Save(ctx, &u)
	return u.ID, errors.Wrap(err, "save")
}

type RegisterDTO struct {
	FirstName  string `validate:"required"`
	SecondName string `validate:"required"`
	Birthdate  string `validate:"required,datetime=2006-01-02"`
	Biography  string
	City       string
	Password   string `validate:"required,"`
}

func (s *Service) GetByID(ctx context.Context, id string) (User, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return User{}, errors.Wrap(err, "get by id")
	}
	return u, nil
}
