package tokens

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/mithcs/probox-api/internal/globals"
	"golang.org/x/crypto/bcrypt"
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
	creds, err := globals.ParseBody[CreateTokensRequest](r.Body)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusBadRequest,
			Title:   "Bad Request.",
			Details: "Could not parse body.",
		})

		return
	}

	userId, err := creds.verify(r.Context())
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

func RefreshTokens(w http.ResponseWriter, r *http.Request) {
	userId, err := retrieveUserId(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusBadRequest,
			Title:   "Bad Request.",
			Details: "You are not authorized",
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

func generateTokens(userId int64) (CreateTokensResponse, error) {
	accessToken, refreshToken, err := globals.GenerateAccessAndRefreshTokens(userId)

	tokens := CreateTokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokens, err
}

func retrieveUserId(ctx context.Context) (int64, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return 0, err
	}

	userId, err := strconv.ParseInt(fmt.Sprintf("%v", claims["uid"]), 10, 64)

	return userId, err
}

func (req *CreateTokensRequest) verify(ctx context.Context) (int64, error) {
	userId, err := verifyCredentials(ctx, req.Username, req.Password)

	return userId, err
}

func verifyCredentials(ctx context.Context, username string, password string) (int64, error) {
	user, err := globals.Queries.GetUserByUsername(ctx, username)
	if err == sql.ErrNoRows {
		return 0, errors.New("Incorrect username.")
	}
	if err != nil {
		return 0, errors.New("Invalid credentials.")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return 0, errors.New("Incorrect password.")
	}

	return user.ID, nil
}
