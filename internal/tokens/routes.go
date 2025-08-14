package tokens

import "github.com/go-chi/chi/v5"

func SetupTokens(router *chi.Mux) {
	router.Post("/tokens", CreateTokens)
	router.Put("/tokens", RefreshTokens)
}
