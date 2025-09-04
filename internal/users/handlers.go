package users

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/mithcs/probox-api/internal/db"
	"github.com/mithcs/probox-api/internal/globals"
	"golang.org/x/crypto/bcrypt"
)

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusBadRequest,
			Title:   "Bad Request.",
			Details: "Invalid User ID.",
		})

		return
	}

	err = compareUserIDWithToken(r.Context(), id)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusUnauthorized,
			Title:   "Unauthorized Action.",
			Details: err.Error(),
		})

		return
	}

	err = deleteUserByID(r.Context(), id)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not delete user.",
		})

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func compareUserIDWithToken(ctx context.Context, id int64) error {
	uid, err := getUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	if uid != id {
		return errors.New("Not authorized to delete user.")
	}

	return nil
}

func getUserIDFromContext(ctx context.Context) (int64, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return -1, err
	}

	uid_s := fmt.Sprintf("%v", claims["uid"])
	uid, err := strconv.ParseInt(uid_s, 10, 64)

	return uid, err
}

func deleteUserByID(ctx context.Context, id int64) error {
	return deleteUserByIDFromDB(ctx, id)
}

func deleteUserByIDFromDB(ctx context.Context, id int64) error {
	random := rand.Text() + rand.Text()
	err := globals.Queries.DeleteUserByID(ctx, db.DeleteUserByIDParams{
		ID:       id,
		Password: random,
	})

	return err
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	identifier := chi.URLParam(r, "identifier")
	if id, err := strconv.ParseInt(identifier, 10, 64); err == nil {
		GetUserByID(w, r, id)
	} else {
		GetUserByUsername(w, r, identifier)
	}
}

func GetUserByID(w http.ResponseWriter, r *http.Request, id int64) {
	user, err := getUserByID(r.Context(), id)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusBadRequest,
			Title:   "Bad Request.",
			Details: "Could not get user.",
		})

		return
	}

	response, err := json.Marshal(user)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not list user.",
		})

		return
	}

	w.Write(response)
}

func getUserByID(ctx context.Context, id int64) (GetUserResponse, error) {
	userRows, err := getUserByIDFromDB(ctx, id)
	if err != nil {
		return GetUserResponse{}, err
	}

	user := GetUserResponse{
		ID:       userRows.ID,
		Username: userRows.Username,
		FullName: userRows.FullName,
	}

	return user, nil
}

func getUserByIDFromDB(ctx context.Context, id int64) (db.GetUserByIDRow, error) {
	userRows, err := globals.Queries.GetUserByID(ctx, id)

	return userRows, err
}

func GetUserByUsername(w http.ResponseWriter, r *http.Request, username string) {
	user, err := getUserByUsername(r.Context(), username)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusBadRequest,
			Title:   "Bad Request.",
			Details: "Could not get user.",
		})

		return
	}

	response, err := json.Marshal(user)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not list user.",
		})

		return
	}

	w.Write(response)
}

func getUserByUsername(ctx context.Context, username string) (GetUserResponse, error) {
	userRows, err := getUserByUsernameFromDB(ctx, username)
	if err != nil {
		return GetUserResponse{}, err
	}

	user := GetUserResponse{
		ID:       userRows.ID,
		Username: userRows.Username,
		FullName: userRows.FullName,
	}

	return user, nil
}

func getUserByUsernameFromDB(ctx context.Context, username string) (db.GetUserByUsernameRow, error) {
	userRows, err := globals.Queries.GetUserByUsername(ctx, username)

	return userRows, err
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

func getAllUsers(ctx context.Context) ([]GetUserResponse, error) {
	userRows, err := getUsersFromDB(ctx)
	if err != nil {
		return []GetUserResponse{}, err
	}

	var users []GetUserResponse
	for _, userRow := range userRows {
		users = append(users, GetUserResponse{
			ID:       userRow.ID,
			Username: userRow.Username,
			FullName: userRow.FullName,
		})
	}

	return users, nil
}

func getUsersFromDB(ctx context.Context) ([]db.GetUsersRow, error) {
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

	err = user.Validate(r.Context())
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

	userID, err := storeUser(r.Context(), user.Username, hashedPass, user.FullName)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not store credentials.",
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

func hashPass(pass string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)

	return string(hashed), err
}

func storeUser(ctx context.Context, username string, password string, fullname string) (int64, error) {
	userID, err := globals.Queries.CreateUser(ctx, db.CreateUserParams{
		Username: username,
		Password: password,
		FullName: fullname,
	})

	return userID, err
}

func generateTokens(userID int64) (CreateUserResponse, error) {
	accessToken, refreshToken, err := globals.GenerateAccessAndRefreshTokens(userID)

	tokens := CreateUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokens, err
}
