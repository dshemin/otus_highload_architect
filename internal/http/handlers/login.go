package handlers

import (
	"encoding/json"
	"github.com/dshemin/otus_highload_architect/internal/passwd"
	"github.com/dshemin/otus_highload_architect/internal/token"
	"github.com/dshemin/otus_highload_architect/internal/user"
	"github.com/dshemin/otus_highload_architect/internal/validator"
	"github.com/pkg/errors"
	"net/http"
)

type Login struct {
	common

	service *token.Service
}

func NewLogin(service *token.Service) *Login {
	return &Login{
		service: service,
	}
}

func (l *Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		l.internalServerError(w, err)
		return
	}

	if err := validator.Struct(req); err != nil {
		l.badRequest(w, err)
		return
	}

	tokenValue, err := l.service.Authenticate(r.Context(), req.ID, req.Password)
	if errors.Is(err, user.ErrNotFound) {
		l.notFound(w, err)
		return
	}
	if errors.Is(err, passwd.ErrPasswordMismatch) {
		l.badRequest(w, err)
		return
	}
	if err != nil {
		l.internalServerError(w, err)
		return
	}

	l.json(w, LoginResponse{
		Token: tokenValue,
	})
}

type LoginRequest struct {
	ID       string `json:"id" validate:"required,uuid"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
