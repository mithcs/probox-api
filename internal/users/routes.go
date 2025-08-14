package users

import "github.com/go-chi/chi/v5"

func SetupUsers(router *chi.Mux) {
	router.Post("/users", CreateUser)
}
