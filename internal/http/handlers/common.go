package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"net/http"
)

type common struct{}

func (c common) json(w http.ResponseWriter, data any) {
	bb, err := json.Marshal(data)
	if err != nil {
		c.internalServerError(w, err)
		return
	}

	_, _ = w.Write(bb)
}

func (c common) internalServerError(w http.ResponseWriter, err error) {
	c.error(w, err, http.StatusInternalServerError)
}

func (c common) notFound(w http.ResponseWriter, err error) {
	c.error(w, err, http.StatusNotFound)
}

type validationErrorResponse struct {
	Errors []string `json:"errors"`
}

func (c common) badRequest(w http.ResponseWriter, err error) {
	var ee validator.ValidationErrors
	if !errors.As(err, &ee) {
		c.error(w, err, http.StatusBadRequest)
		return
	}

	r := validationErrorResponse{
		Errors: make([]string, 0, len(ee)),
	}
	for _, e := range ee {
		r.Errors = append(r.Errors, e.Error())
	}

	data, err := json.Marshal(r)
	if err != nil {
		c.internalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write(data)
}

func (common) error(w http.ResponseWriter, err error, statusCode int) {
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(fmt.Sprintf(`{
	"error": %q
}`, err)))
}
