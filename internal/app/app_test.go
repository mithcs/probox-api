package app

import (
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestNew(t *testing.T) {
	ctx := t.Context()
	addr := ":8080"
	db := Database{
		Driver: "sqlite3",
		DSN:    ":memory:",
	}
	router := chi.NewRouter()

	app := New(ctx, addr, db, router)

	if app.ctx != ctx {
		t.Errorf("ctx is not same, expected to be same")
	}

	if app.addr != addr {
		t.Errorf("addr is not same, expected to be same")
	}

	if app.db != db {
		t.Errorf("db is not same, expected to be same")
	}

	if app.router != router {
		t.Errorf("router is not same, expected to be same")
	}
}
