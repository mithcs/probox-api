package solutions

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/mithcs/probox-api/internal/globals"
)

func SetupSolutions(router *chi.Mux) {
	router.Get("/solutions", GetSolutions)

	router.
		With(jwtauth.Verifier(globals.AccessTokenAuth)).
		With(globals.AuthenticatorMiddleware(globals.AccessTokenAuth)).
		Post("/solutions", CreateSolution)
}
