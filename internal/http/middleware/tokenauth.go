package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type tokenExistsChecker interface {
	IsExists(ctx context.Context, token string) (bool, error)
}

func TokenAuth(checker tokenExistsChecker) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			token := strings.TrimSpace(r.Header.Get("Authorization"))
			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{
"error": "Token not provided"
}`))
				return
			}

			ok, err := checker.IsExists(r.Context(), token)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(fmt.Sprintf(`{
"error": %q
}`, err)))
				return
			}

			if !ok {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
