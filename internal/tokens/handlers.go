package tokens

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/jwtauth/v5"
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
	var user CreateTokensRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := globals.ReturnErrorResponse(
			"Server Error.",
			"Could not read the body.",
		)
		w.Write(response)

		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := globals.ReturnErrorResponse(
			"Bad Request.",
			"Could not parse body.",
		)
		w.Write(response)

		return
	}

	userId, err := user.verify()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := globals.ReturnErrorResponse(
			"Bad Request.",
			err.Error(),
		)
		w.Write(response)

		return
	}

	accessToken, refreshToken, err := globals.GenerateAccessAndRefreshTokens(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := globals.ReturnErrorResponse(
			"Server Error.",
			"Could not generate tokens for you.",
		)
		w.Write(response)

		return
	}

	tokens := CreateTokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	response, err := json.Marshal(tokens)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := globals.ReturnErrorResponse(
			"Server Error.",
			"Could not give you tokens.",
		)
		w.Write(response)

		return
	}

	w.Write(response)
}

func RefreshTokens(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := globals.ReturnErrorResponse(
			"Bad Request",
			"You are not authorized",
		)
		w.Write(response)

		return
	}

	userId, err := strconv.Atoi(fmt.Sprintf("%v", claims["user_id"]))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := globals.ReturnErrorResponse(
			"Bad Request",
			"Invalid user_id",
		)
		w.Write(response)

		return
	}

	accessToken, refreshToken, err := globals.GenerateAccessAndRefreshTokens(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := globals.ReturnErrorResponse(
			"Server Error.",
			"Could not generate tokens for you.",
		)
		w.Write(response)

		return
	}

	tokens := RefreshTokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	response, err := json.Marshal(tokens)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := globals.ReturnErrorResponse(
			"Server Error.",
			"Could not give you tokens.",
		)
		w.Write(res)

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
