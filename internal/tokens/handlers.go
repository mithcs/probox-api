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
	creds, err := parseBody(r.Body)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusBadRequest,
			Title:   "Bad Request.",
			Details: "Could not parse body.",
		})

		return
	}

	userId, err := creds.verify()
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusBadRequest,
			Title:   "Bad Request.",
			Details: err.Error(),
		})

		return
	}

	tokens, err := generateTokens(userId)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not generate tokens for you.",
		})

		return
	}

	response, err := json.Marshal(tokens)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not give you tokens.",
		})

		return
	}

	w.Write(response)
}

// TODO: refactor RefreshTokens()

func RefreshTokens(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusBadRequest,
			Title:   "Bad Request.",
			Details: "You are not authorized.",
		})

		return
	}

	userId, err := strconv.ParseInt(fmt.Sprintf("%v", claims["uid"]), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusBadRequest,
			Title:   "Bad Request.",
			Details: "You are not authorized",
		})

		return
	}

	accessToken, refreshToken, err := globals.GenerateAccessAndRefreshTokens(userId)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not generate tokens for you.",
		})

		return
	}

	tokens := RefreshTokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	response, err := json.Marshal(tokens)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not give you tokens.",
		})

		return
	}

	w.Write(response)
}

func (req *CreateTokensRequest) verify() (int64, error) {
	userId := int64(1)

	// verify username and password here
	if req.Username == "username" &&
		req.Password == "password" {
		return -1, errors.New("Invalid Credentials.")
	}

	return userId, nil
}

func parseBody(body io.ReadCloser) (CreateTokensRequest, error) {
	var creds CreateTokensRequest

	readBody, err := io.ReadAll(body)
	if err != nil {
		return creds, err
	}

	err = json.Unmarshal(readBody, &creds)
	if err != nil {
		return creds, err
	}

	return creds, nil
}

func generateTokens(userId int64) (CreateTokensResponse, error) {
	accessToken, refreshToken, err := globals.GenerateAccessAndRefreshTokens(userId)

	tokens := CreateTokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokens, err
}
