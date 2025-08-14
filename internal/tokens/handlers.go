package tokens

import "net/http"

func CreateTokens(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create new pair of access and refresh tokens from username/password"))
}

func RefreshTokens(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create new pair of access and refresh tokens from refresh token"))
}
