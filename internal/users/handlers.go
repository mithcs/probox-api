package users

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"github.com/mithcs/probox-api/internal/db"
	"github.com/mithcs/probox-api/internal/globals"
	validator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
)

type GetUsersResponse struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := getAllUsers(r.Context())
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not get all users.",
		})

		return
	}

	response, err := json.Marshal(users)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not list all users.",
		})

		return
	}

	w.Write(response)
}

func getAllUsers(ctx context.Context) ([]GetUsersResponse, error) {
	userRows, err := getUsersFromDb(ctx)
	if err != nil {
		return []GetUsersResponse{}, err
	}

	var users []GetUsersResponse
	for _, userRow := range userRows {
		users = append(users, GetUsersResponse{
			Id:       userRow.ID,
			Username: userRow.Username,
		})
	}

	return users, nil
}

func getUsersFromDb(ctx context.Context) ([]db.GetUsersRow, error) {
	userRows, err := globals.Queries.GetUsers(ctx)

	return userRows, err
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

	err = user.validate(r.Context())
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

func (user *CreateUserRequest) validate(ctx context.Context) error {
	if err := validateUsername(user.Username); err != nil {
		return err
	}

	if err := validatePassword(user.Password); err != nil {
		return errors.New("Password is insecure.")
	}

	if exists := usernameExists(ctx, user.Username); exists {
		return errors.New("Username already exists.")
	}

	return nil
}

func validateUsername(username string) error {
	// must start with alphabetic character
	// may include underscore and/or number
	// min 3 chars
	// max 16 chars
	r := regexp.MustCompile(`^[A-Za-z]\w{2,15}$`)
	if !r.MatchString(username) {
		return errors.New("Invalid username.")
	}

	return nil
}

func validatePassword(password string) error {
	err := validator.Validate(password, 60)
	if err != nil {
		return errors.New("Password is insecure.")
	}

	return nil
}

func usernameExists(ctx context.Context, username string) bool {
	user, _ := globals.Queries.GetUserByUsername(ctx, username)
	if user.ID == 0 || user.Username == "" || user.Password == "" {
		return false
	}

	return true
}
