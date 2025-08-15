package tokens

import (
	"encoding/json"
	"net/http"

	"github.com/mithcs/probox-api/internal/globals"
)

type CreateTokensResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokensResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func CreateTokens(w http.ResponseWriter, r *http.Request) {
	// get username and password from request body
	// verify username and password
	userId := 1

	accessToken, refreshToken, err := globals.GenerateAccessAndRefreshTokens(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tokens := CreateTokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	response, err := json.Marshal(tokens)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(response)
}

func RefreshTokens(w http.ResponseWriter, r *http.Request) {
	// get refresh token from body
	// verify refresh token
	userId := 1

	accessToken, refreshToken, err := globals.GenerateAccessAndRefreshTokens(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tokens := RefreshTokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	response, err := json.Marshal(tokens)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(response)
}
