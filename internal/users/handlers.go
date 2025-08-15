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
	var user CreateUserRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = user.validate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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
	// otherwise, a formatted error message is provided
	err := passValidator.Validate(user.Password, 60)
	if err != nil {
		return err
	}

	return nil
}
