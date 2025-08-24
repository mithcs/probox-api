package users

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"github.com/mithcs/probox-api/internal/db"
	"github.com/mithcs/probox-api/internal/globals"
	passValidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	user, err := globals.ParseBody[CreateUserRequest](r.Body)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusBadRequest,
			Title:   "Bad Request.",
			Details: "Could not parse body.",
		})

		return
	}

	err = user.validate()
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusBadRequest,
			Title:   "Bad Request.",
			Details: err.Error(),
		})

		return
	}

	hashedPass, err := hashPass(user.Password)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not hash your password.",
		})

		return
	}

	userId, err := storeUser(r.Context(), user.Username, hashedPass)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not store credentials.",
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

func (user *CreateUserRequest) validate() error {
	// must start with alphabetic character
	// may include underscore and/or number
	// min 3 chars
	// max 16 chars
	r := regexp.MustCompile(`^[A-Za-z]\w{2,15}$`)
	if !r.MatchString(user.Username) {
		return errors.New("Invalid username.")
	}

	// if the password has enough entropy, err is nil
	err := passValidator.Validate(user.Password, 60)
	if err != nil {
		return errors.New("Password is insecure.")
	}

	return nil
}

func hashPass(pass string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)

	return string(hashed), err
}

func storeUser(ctx context.Context, username string, password string) (int64, error) {
	userId, err := globals.Queries.CreateUser(ctx, db.CreateUserParams{
		Username: username,
		Password: password,
	})

	return userId, err
}

func generateTokens(userId int64) (CreateUserResponse, error) {
	accessToken, refreshToken, err := globals.GenerateAccessAndRefreshTokens(userId)

	tokens := CreateUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokens, err
}
