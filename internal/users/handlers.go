package users

import (
	"encoding/json"
	"net/http"

	"github.com/mithcs/probox-api/internal/globals"
)

type CreateUserResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	// get username and password from request body
	// validate username and password
	// store username and hash of password
	userId := 1

	accessToken, refreshToken, err := globals.GenerateAccessAndRefreshTokens(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tokens := CreateUserResponse{
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
