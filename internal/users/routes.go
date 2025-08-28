package users

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/mithcs/probox-api/internal/globals"
)

func SetupUsers(router *chi.Mux) {
	router.Get("/users", GetUsers)
	router.Get("/users/{identifier}", GetUser)
	router.Post("/users", CreateUser)

	router.
		With(jwtauth.Verifier(globals.AccessTokenAuth)).
		With(globals.AuthenticatorMiddleware(globals.AccessTokenAuth)).
		Delete("/users/{id}", DeleteUser)
}
