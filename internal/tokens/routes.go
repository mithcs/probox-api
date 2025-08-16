package tokens

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/mithcs/probox-api/internal/globals"
)

func SetupTokens(router *chi.Mux) {
	router.Post("/tokens", CreateTokens)

	router.
		With(jwtauth.Verifier(globals.RefreshTokenAuth)).
		With(globals.AuthenticatorMiddleware(globals.RefreshTokenAuth)).
		Put("/tokens", RefreshTokens)
}
