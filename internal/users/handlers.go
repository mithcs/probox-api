package users

import "net/http"

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create new user and generate access/refresh tokens"))
}
