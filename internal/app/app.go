package app

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mithcs/probox-api/internal/db"
	"github.com/mithcs/probox-api/internal/globals"
	"github.com/mithcs/probox-api/internal/problems"
	"github.com/mithcs/probox-api/internal/solutions"
	"github.com/mithcs/probox-api/internal/tokens"
	"github.com/mithcs/probox-api/internal/users"
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
	conn, err := sql.Open(app.db.Driver, app.db.DSN)
	if err != nil {
		return err
	}
	defer conn.Close()

	globals.Queries = db.New(conn)

	// err = conn.Ping()
	// if err != nil {
	// 	panic(err)
	// }

	registerRoutes(app.router)

	return http.ListenAndServe(app.addr, app.router)
}

func registerRoutes(router *chi.Mux) {
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	tokens.SetupTokens(router)
	users.SetupUsers(router)
	problems.SetupProblems(router)
	solutions.SetupSolutions(router)
}
