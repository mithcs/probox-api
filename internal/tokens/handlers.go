package tokens

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/mithcs/probox-api/internal/globals"
)

type CreateTokensRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateTokensResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokensResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func CreateTokens(w http.ResponseWriter, r *http.Request) {
	var error globals.ErrorResponse

	var user CreateTokensRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error.Title = "Server Error."
		error.Details = "Could not read the body."

		response, _ := json.Marshal(error)
		w.Write(response)

		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error.Title = "Bad Request."
		error.Details = "Could not parse body."

		response, _ := json.Marshal(error)
		w.Write(response)

		return
	}

	userId, err := user.verify()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error.Title = "Bad Request."
		error.Details = err.Error()

		response, _ := json.Marshal(error)
		w.Write(response)

		return
	}

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

func (req *CreateTokensRequest) verify() (int, error) {
	userId := 1

	// verify username and password here
	if req.Username == "username" &&
		req.Password == "password" {
		return -1, errors.New("Invalid Credentials.")
	}

	return userId, nil
}
