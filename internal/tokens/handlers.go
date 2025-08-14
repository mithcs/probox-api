package tokens

import "net/http"

func CreateTokens(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("generate new pair of access and refresh tokens from refresh token"))
}

func RefreshTokens(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("generate new access token using refresh token"))
}
