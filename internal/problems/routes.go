package problems

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/mithcs/probox-api/internal/globals"
)

func SetupProblems(router *chi.Mux) {
	router.
		With(jwtauth.Verifier(globals.AccessTokenAuth)).
		With(globals.AuthenticatorMiddleware(globals.AccessTokenAuth)).
		Post("/problems", CreateProblem)
}
