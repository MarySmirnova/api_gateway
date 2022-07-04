package internal

import (
	"net/http"

	"github.com/MarySmirnova/api_gateway/internal/api"
	"github.com/MarySmirnova/api_gateway/internal/config"

	log "github.com/sirupsen/logrus"
)

type Application struct {
	cfg config.Application
}

func NewApplication(cfg config.Application) *Application {
	return &Application{
		cfg: cfg,
	}
}

func (a *Application) StartServer() {
	srv := api.NewGateway(a.cfg.Server)
	s := srv.GetHTTPServer()

	log.WithField("listen", s.Addr).Info("start server")

	err := s.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		log.WithError(err).Error("the channel raised an error")
		return
	}
}
