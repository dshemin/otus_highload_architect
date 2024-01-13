package application

import (
	"context"
	"github.com/dshemin/otus_highload_architect/internal/config"
	"github.com/dshemin/otus_highload_architect/internal/http"
	"github.com/pkg/errors"
	"log"
	"os"
	"os/signal"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return errors.Wrap(err, "load config")
	}

	log.Printf("Config %#v", cfg)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	c, err := initContainer(ctx, cfg)
	if err != nil {
		return errors.Wrap(err, "init container")
	}

	err = http.Serve(
		ctx,
		c.UserService,
		c.TokenService,
		cfg.HTTP.Host,
		cfg.HTTP.Port,
	)
	return errors.Wrap(err, "http serve")
}
