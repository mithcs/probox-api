package globals

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/jwtauth/v5"
)

var AccessTokenAuth *jwtauth.JWTAuth
var RefreshTokenAuth *jwtauth.JWTAuth

func init() {
	atAlg := "HS256"
	atSignKey := os.Getenv("ACCESS_TOKEN_SECRET")
	AccessTokenAuth = jwtauth.New(atAlg, []byte(atSignKey), nil)

	rtAlg := "HS256"
	rtSignKey := os.Getenv("REFRESH_TOKEN_SECRET")
	RefreshTokenAuth = jwtauth.New(rtAlg, []byte(rtSignKey), nil)
}

func GenerateAccessAndRefreshTokens(userID int64) (string, string, error) {
	accessToken, err := GenerateAccessToken(userID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, err
}

func GenerateAccessToken(userID int64) (string, error) {
	accessToken, err := generateToken(userID, AccessTokenAuth)

	return accessToken, err
}

func GenerateRefreshToken(userID int64) (string, error) {
	refreshToken, err := generateToken(userID, RefreshTokenAuth)

	return refreshToken, err
}

func generateToken(userID int64, auth *jwtauth.JWTAuth) (string, error) {
	claims := map[string]any{
		"uid": userID,
	}

	jwtauth.SetIssuedNow(claims)
	jwtauth.SetExpiryIn(claims, 24*time.Hour)

	_, token, err := auth.Encode(claims)
	if err != nil {
		return "", err
	}

	return token, err
}

func GetUserIDFromContext(ctx context.Context) (int64, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return -1, err
	}

	uid_s := fmt.Sprintf("%v", claims["uid"])
	uid, err := strconv.ParseInt(uid_s, 10, 64)

	return uid, err
}

func AuthenticatorMiddleware(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			token, _, err := jwtauth.FromContext(r.Context())

			if err != nil {
				WriteErrorResponse(w, Error{
					Status:  http.StatusUnauthorized,
					Title:   "Bad Request.",
					Details: err.Error(),
				})

				return
			}

			if token == nil {
				WriteErrorResponse(w, Error{
					Status:  http.StatusUnauthorized,
					Title:   "Bad Request.",
					Details: "Token is empty.",
				})

				return
			}

			// Token is authenticated, pass it through
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hfn)
	}
}
