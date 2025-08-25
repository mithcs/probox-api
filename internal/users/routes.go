package users

import "github.com/go-chi/chi/v5"

func SetupUsers(router *chi.Mux) {
	router.Get("/users", GetUsers)
	router.Post("/users", CreateUser)
}
