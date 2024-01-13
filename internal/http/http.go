package http

import (
	"context"
	"github.com/dshemin/otus_highload_architect/internal/http/handlers"
	"github.com/dshemin/otus_highload_architect/internal/http/middleware"
	"github.com/dshemin/otus_highload_architect/internal/token"
	"github.com/dshemin/otus_highload_architect/internal/user"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

func Serve(
	ctx context.Context,
	userService *user.Service,
	tokeService *token.Service,
	host string,
	port uint16,
) error {
	mx := initHandlers(userService, tokeService)

	server := createServer(host, port, mx)
	go func() {
		<-ctx.Done()

		sdCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := server.Shutdown(sdCtx)
		if err != nil {
			log.Printf("failed to shutdown: %q", err)
		}
	}()

	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		err = nil
	}

	return errors.Wrap(err, "listen and serve")
}

func initHandlers(
	userService *user.Service,
	tokenService *token.Service,
) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.AllowContentType("application/json"))
	r.Use(chimiddleware.Heartbeat("/healthz"))
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	r.Post("/login", toHandlerFunc(handlers.NewLogin(tokenService)))

	r.Post("/user/register", toHandlerFunc(handlers.NewRegister(userService)))
	r.Route("/user", func(r chi.Router) {
		r.Use(middleware.TokenAuth(tokenService))
		r.Get("/get/{id}", toHandlerFunc(handlers.NewGetByID(userService)))
	})

	return r
}

func createServer(host string, port uint16, h http.Handler) *http.Server {
	addr := net.JoinHostPort(host, strconv.FormatUint(uint64(port), 10))
	return &http.Server{Addr: addr, Handler: h}
}

func toHandlerFunc(h http.Handler) http.HandlerFunc {
	return h.ServeHTTP
}
