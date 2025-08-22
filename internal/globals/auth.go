package globals

import (
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
)

var AccessTokenAuth *jwtauth.JWTAuth
var RefreshTokenAuth *jwtauth.JWTAuth

func init() {
	atAlg := "HS256"
	atSignKey := []byte("a-string-secret-256-bits-long-boom")
	AccessTokenAuth = jwtauth.New(atAlg, atSignKey, nil)

	rtAlg := "HS256"
	rtSignKey := []byte("b-string-secret-256-bits-long-boom")
	RefreshTokenAuth = jwtauth.New(rtAlg, rtSignKey, nil)
}

func GenerateAccessToken(userId int64) (string, error) {
	claims := map[string]any{
		"uid": userId,
	}

	jwtauth.SetIssuedNow(claims)
	jwtauth.SetExpiryIn(claims, 10*time.Minute)

	_, accessToken, err := AccessTokenAuth.Encode(claims)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func GenerateRefreshToken(userId int64) (string, error) {
	claims := map[string]any{
		"uid": 1,
	}

	jwtauth.SetIssuedNow(claims)
	jwtauth.SetExpiryIn(claims, 24*time.Hour)

	_, refreshToken, err := RefreshTokenAuth.Encode(claims)
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func GenerateAccessAndRefreshTokens(userId int64) (string, string, error) {
	accessToken, err := GenerateAccessToken(userId)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := GenerateRefreshToken(userId)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, err
}

func AuthenticatorMiddleware(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			token, _, err := jwtauth.FromContext(r.Context())

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write(ReturnErrorResponse(
					"Bad Request",
					err.Error(),
				))
				return
			}

			if token == nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write(ReturnErrorResponse(
					"Bad Request",
					"Token is empty.",
				))
				return
			}

			// Token is authenticated, pass it through
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hfn)
	}
}
