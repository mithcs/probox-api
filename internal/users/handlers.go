package users

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"

	"github.com/mithcs/probox-api/internal/globals"
	passValidator "github.com/wagslane/go-password-validator"
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
	var error globals.ErrorResponse

	var user CreateUserRequest
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

	err = user.validate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error.Title = "Bad Request."
		error.Details = err.Error()

		response, _ := json.Marshal(error)
		w.Write(response)

		return
	}

	// store username and hash of password
	userId := 1

	accessToken, refreshToken, err := globals.GenerateAccessAndRefreshTokens(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error.Title = "Server Error."
		error.Details = "Could not generate tokens for you."

		response, _ := json.Marshal(error)
		w.Write(response)

		return
	}

	tokens := CreateUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	response, err := json.Marshal(tokens)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error.Title = "Server Error."
		error.Details = "Could not give you tokens."

		response, _ := json.Marshal(error)
		w.Write(response)

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
