package globals

import (
	"time"

	"github.com/go-chi/jwtauth/v5"
)

var accessTokenAuth *jwtauth.JWTAuth
var refreshTokenAuth *jwtauth.JWTAuth

func init() {
	atAlg := "HS256"
	atSignKey := []byte("a-string-secret-256-bits-long-boom")
	accessTokenAuth = jwtauth.New(atAlg, atSignKey, nil)

	rtAlg := "HS256"
	rtSignKey := []byte("b-string-secret-256-bits-long-boom")
	refreshTokenAuth = jwtauth.New(rtAlg, rtSignKey, nil)
}

func GenerateAccessToken(userId int) (string, error) {
	claims := map[string]any{
		"user_id": userId,
	}

	jwtauth.SetIssuedNow(claims)
	jwtauth.SetExpiryIn(claims, 10*time.Minute)

	_, accessToken, err := accessTokenAuth.Encode(claims)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func GenerateRefreshToken(userId int) (string, error) {
	claims := map[string]any{
		"user_id": 1,
	}

	jwtauth.SetIssuedNow(claims)
	jwtauth.SetExpiryIn(claims, 24*time.Hour)

	_, refreshToken, err := refreshTokenAuth.Encode(claims)
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}
