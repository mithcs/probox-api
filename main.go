package main

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/mithcs/probox-api/internal/app"
)

func main() {
	ctx := context.Background()
	addr := ":8080"
	db := app.Database{
		Driver: "sample",
		DSN:    "sample",
	}
	router := chi.NewRouter()

	app := app.New(ctx, addr, db, router)

	err := app.Run()
	if err != nil {
		panic(err)
	}
}
