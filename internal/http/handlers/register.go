package handlers

import (
	"encoding/json"
	"github.com/dshemin/otus_highload_architect/internal/user"
	"github.com/dshemin/otus_highload_architect/internal/validator"
	"net/http"
)

type Register struct {
	common

	service *user.Service
}

func NewRegister(service *user.Service) *Register {
	return &Register{
		service: service,
	}
}

func (l *Register) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		l.internalServerError(w, err)
		return
	}

	if err := validator.Struct(req); err != nil {
		l.badRequest(w, err)
		return
	}

	id, err := l.service.Register(r.Context(), user.RegisterDTO{
		FirstName:  req.FirstName,
		SecondName: req.SecondName,
		Birthdate:  req.Birthdate,
		Biography:  req.Biography,
		City:       req.City,
		Password:   req.Password,
	})
	if err != nil {
		l.internalServerError(w, err)
		return
	}

	l.json(w, RegisterResponse{
		UserID: id,
	})
}

type RegisterRequest struct {
	FirstName  string `json:"first_name" validate:"required"`
	SecondName string `json:"second_name" validate:"required"`
	Birthdate  string `json:"birthdate" validate:"datetime=2006-01-02"`
	Biography  string `json:"biography"`
	City       string `json:"city"`
	Password   string `json:"password" validate:"required"`
}

type RegisterResponse struct {
	UserID string `json:"user_id"`
}
