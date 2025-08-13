package build

import (
	"context"
	"net/http"

	"otusgruz/config"

	"github.com/gorilla/mux"
)

type Builder struct {
	config config.Config

	shutdown shutdown

	http struct {
		router *mux.Router
		server *http.Server
	}
}

func New(ctx context.Context, conf config.Config) *Builder {
	b := Builder{config: conf} //nolint:exhaustruct

	return &b
}
