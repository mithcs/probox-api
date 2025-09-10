package users

import (
	"context"
	"errors"
	"net/http"
	"regexp"

	"github.com/mithcs/probox-api/internal/globals"
	validator "github.com/wagslane/go-password-validator"
)

type GetUserResponse struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FullName  string `json:"full_name"`
	CreatedAt string `json:"created_at"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

type CreateUserResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (user *CreateUserRequest) Validate(ctx context.Context) *globals.HTTPError {
	if err := validateUsername(user.Username); err != nil {
		return &globals.HTTPError{
			Status: http.StatusBadRequest,
			Title:  "Invalid Username.",
			Err:    errors.New("Username does not meet requirements."),
		}
	}

	if err := validatePassword(user.Password); err != nil {
		return &globals.HTTPError{
			Status: http.StatusBadRequest,
			Title:  "Invalid Password.",
			Err:    errors.New("Password is not secure enough."),
		}
	}

	if err := validateFullName(user.FullName); err != nil {
		return &globals.HTTPError{
			Status: http.StatusBadRequest,
			Title:  "Invalid Fullname.",
			Err:    errors.New("Fullname does not meet requirements."),
		}
	}

	if exists := usernameExists(ctx, user.Username); exists {
		return &globals.HTTPError{
			Status: http.StatusBadRequest,
			Title:  "Invalid Username",
			Err:    errors.New("Username is already occupied."),
		}
	}

	return nil
}

func validateUsername(username string) error {
	// must start with alphabetic character
	// may include underscore and/or number
	// min 3 chars
	// max 24 chars
	r := regexp.MustCompile(`^[A-Za-z]\w{1,22}\w$`)
	if !r.MatchString(username) {
		return errors.New("Invalid Username.")
	}

	return nil
}

func validatePassword(password string) error {
	err := validator.Validate(password, 60)
	if err != nil {
		return err
	}

	return nil
}

func validateFullName(fullname string) error {
	// must start and end with alphabet
	// may include space
	// min 3 chars
	// max 24 chars
	r := regexp.MustCompile(`^[A-Za-z][A-Za-z ]{1,22}[A-Za-z]$`)
	if !r.MatchString(fullname) {
		return errors.New("Invalid Fullname.")
	}

	return nil
}

func usernameExists(ctx context.Context, username string) bool {
	user, _ := globals.Queries.GetUserByUsername(ctx, username)
	if user.ID == 0 || user.Username == "" || user.FullName == "" {
		return false
	}

	return true
}
