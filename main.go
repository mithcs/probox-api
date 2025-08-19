package main

import (
	"context"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/mithcs/probox-api/internal/app"
)

func main() {
	ctx := context.Background()
	addr := ":8080"
	db := app.Database{
		Driver: "pgx",
		DSN:    os.Getenv("DATABASE_DSN"),
	}
	router := chi.NewRouter()

	app := app.New(ctx, addr, db, router)

	err := app.Run()
	if err != nil {
		panic(err)
	}
}
