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
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
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

	userID, err := creds.verify(r.Context())
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusBadRequest,
			Title:   "Bad Request.",
			Details: err.Error(),
		})

		return
	}

	tokens, err := generateTokens(userID)
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
	userID, err := retrieveUserID(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusBadRequest,
			Title:   "Bad Request.",
			Details: "You are not authorized",
		})

		return
	}

	tokens, err := generateTokens(userID)
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

func generateTokens(userID int64) (CreateTokensResponse, error) {
	accessToken, refreshToken, err := globals.GenerateAccessAndRefreshTokens(userID)

	tokens := CreateTokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokens, err
}

func retrieveUserID(ctx context.Context) (int64, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return 0, err
	}

	userID, err := strconv.ParseInt(fmt.Sprintf("%v", claims["uid"]), 10, 64)

	return userID, err
}

func (req *CreateTokensRequest) verify(ctx context.Context) (int64, error) {
	userID, err := verifyCredentials(ctx, req.Username, req.Password)

	return userID, err
}

func verifyCredentials(ctx context.Context, username string, password string) (int64, error) {
	user, err := globals.Queries.GetUserByUsername(ctx, username)
	if err == sql.ErrNoRows {
		return 0, errors.New("Incorrect Username.")
	}
	if err != nil {
		return 0, errors.New("Invalid Credentials.")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return 0, errors.New("Incorrect Password.")
	}

	return user.ID, nil
}
