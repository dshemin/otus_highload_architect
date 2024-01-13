package handlers

import (
	"encoding/json"
	"github.com/dshemin/otus_highload_architect/internal/user"
	"github.com/dshemin/otus_highload_architect/internal/validator"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type GetByID struct {
	common

	service *user.Service
}

func NewGetByID(service *user.Service) *GetByID {
	return &GetByID{
		service: service,
	}
}

func (l *GetByID) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		l.internalServerError(w, err)
		return
	}

	id := chi.URLParam(r, "id")

	err = l.validate(id)
	if err != nil {
		l.badRequest(w, err)
		return
	}

	usr, err := l.service.GetByID(r.Context(), id)
	if errors.Is(err, user.ErrNotFound) {
		l.badRequest(w, err)
		return
	}
	if err != nil {
		l.internalServerError(w, err)
		return
	}

	l.json(w, GetByIDResponse{
		ID:         usr.ID,
		FirstName:  usr.FirstName,
		SecondName: usr.SecondName,
		Birthdate:  usr.Birthdate.Format(time.DateOnly),
		Biography:  usr.Biography,
		City:       usr.City,
	})
}

func (*GetByID) validate(id string) error {
	return validator.Var(id, "required,uuid")
}

type GetByIDResponse struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Birthdate  string `json:"birthdate"`
	Biography  string `json:"biography"`
	City       string `json:"city"`
}
