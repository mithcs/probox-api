package users

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mithcs/probox-api/internal/globals"
	"golang.org/x/crypto/bcrypt"
)

func compareUserIDWithToken(ctx context.Context, id int64) *globals.HTTPError {
	uid, err := globals.GetUserIDFromContext(ctx)
	if err != nil {
		return &globals.HTTPError{
			Status:      http.StatusInternalServerError,
			Title:       "Invalid User ID.",
			Description: "Could not get user id from token.",
		}
	}

	if uid != id {
		return &globals.HTTPError{
			Status:      http.StatusBadRequest,
			Title:       "Unauthorized Action.",
			Description: "Not authorized to delete this user.",
		}
	}

	return nil
}

func deleteUserByID(ctx context.Context, id int64) *globals.HTTPError {
	err := deleteUserByIDFromDB(ctx, id)
	if err != nil {
		return &globals.HTTPError{
			Status:      http.StatusInternalServerError,
			Title:       "Database Error.",
			Description: "Could not get user from database.",
		}
	}

	return nil
}

func getUserByID(ctx context.Context, id int64) (GetUserResponse, *globals.HTTPError) {
	userRows, err := getUserByIDFromDB(ctx, id)
	if err != nil {
		return GetUserResponse{}, &globals.HTTPError{
			Status:      http.StatusInternalServerError,
			Title:       "Database Error.",
			Description: "Could not get user from database.",
		}
	}

	user := GetUserResponse{
		ID:       userRows.ID,
		Username: userRows.Username,
		FullName: userRows.FullName,
	}

	return user, nil
}

func getUserByUsername(ctx context.Context, username string) (GetUserResponse, *globals.HTTPError) {
	userRows, err := getUserByUsernameFromDB(ctx, username)
	if err != nil {
		return GetUserResponse{}, &globals.HTTPError{
			Status:      http.StatusInternalServerError,
			Title:       "Database Error.",
			Description: "Could not get user from database.",
		}
	}

	user := GetUserResponse{
		ID:       userRows.ID,
		Username: userRows.Username,
		FullName: userRows.FullName,
	}

	return user, nil
}

func getAllUsers(ctx context.Context) ([]GetUserResponse, *globals.HTTPError) {
	userRows, err := getUsersFromDB(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return []GetUserResponse{}, &globals.HTTPError{
				Status:      http.StatusNoContent,
				Title:       "No Content.",
				Description: "No users to show.",
			}
		}

		return []GetUserResponse{}, &globals.HTTPError{
			Status:      http.StatusInternalServerError,
			Title:       "Database Error.",
			Description: "Could not get users from database.",
		}
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

func hashPass(pass string) (string, *globals.HTTPError) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", &globals.HTTPError{
			Status:      http.StatusInternalServerError,
			Title:       "Hashing Error.",
			Description: "Could not hash password.",
		}
	}

	return string(hashed), nil
}

func storeUser(ctx context.Context, username string, password string, fullname string) (int64, *globals.HTTPError) {
	id, err := storeUserInDB(ctx, username, password, fullname)
	if err != nil {
		return -1, &globals.HTTPError{
			Status:      http.StatusInternalServerError,
			Title:       "Database Error.",
			Description: "Could not store user in database.",
		}
	}

	return id, nil
}

func generateTokens(userID int64) (CreateUserResponse, *globals.HTTPError) {
	accessToken, refreshToken, err := globals.GenerateAccessAndRefreshTokens(userID)
	if err != nil {
		return CreateUserResponse{}, &globals.HTTPError{
			Status:      http.StatusInternalServerError,
			Title:       "Token Generation Error.",
			Description: "Could not generate access and refresh tokens.",
		}
	}

	tokens := CreateUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokens, nil
}

func getIDFromRequest(r *http.Request) (int64, *globals.HTTPError) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		return -1, &globals.HTTPError{
			Status:      http.StatusBadRequest,
			Title:       "Invalid ID.",
			Description: "Could not parse id.",
		}
	}

	return id, nil
}
