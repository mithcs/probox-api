package app

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Database struct {
	Driver string
	DSN    string
}

type App struct {
	ctx    context.Context
	addr   string
	db     Database
	router *chi.Mux
}

func New(ctx context.Context, addr string, db Database, router *chi.Mux) App {
	app := App{
		ctx:    ctx,
		addr:   addr,
		db:     db,
		router: router,
	}

	return app
}

func (app *App) Run() error {
	// setup db connection here

	registerRoutes(app.router)

	http.ListenAndServe(app.addr, app.router)

	return nil
}

func registerRoutes(router *chi.Mux) {
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
}
