package users

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/mithcs/probox-api/internal/globals"
	"golang.org/x/crypto/bcrypt"
)

func deleteUserByID(ctx context.Context, id int64) *globals.HTTPError {
	err := deleteUserByIDFromDB(ctx, id)
	if err != nil {
		return &globals.HTTPError{
			Status: http.StatusInternalServerError,
			Title:  "Database Error.",
			Err:    errors.New("Could not get user from database."),
		}
	}

	return nil
}

func getUserByID(ctx context.Context, id int64) (GetUserResponse, *globals.HTTPError) {
	userRows, err := getUserByIDFromDB(ctx, id)
	if err != nil {
		return GetUserResponse{}, &globals.HTTPError{
			Status: http.StatusInternalServerError,
			Title:  "Database Error.",
			Err:    errors.New("Could not get user from database."),
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
			Status: http.StatusInternalServerError,
			Title:  "Database Error.",
			Err:    errors.New("Could not get user from database."),
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
				Status: http.StatusNoContent,
				Title:  "No Content.",
				Err:    errors.New("No users to show."),
			}
		}

		return []GetUserResponse{}, &globals.HTTPError{
			Status: http.StatusInternalServerError,
			Title:  "Database Error.",
			Err:    errors.New("Could not get users from database."),
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
			Status: http.StatusInternalServerError,
			Title:  "Hashing Error.",
			Err:    errors.New("Could not hash password."),
		}
	}

	return string(hashed), nil
}

func storeUser(ctx context.Context, username string, password string, fullname string) (int64, *globals.HTTPError) {
	id, err := storeUserInDB(ctx, username, password, fullname)
	if err != nil {
		return -1, &globals.HTTPError{
			Status: http.StatusInternalServerError,
			Title:  "Database Error.",
			Err:    errors.New("Could not store user in database."),
		}
	}

	return id, nil
}

func generateTokens(userID int64) (CreateUserResponse, *globals.HTTPError) {
	accessToken, refreshToken, err := globals.GenerateTokens(userID)
	if err != nil {
		return CreateUserResponse{}, &globals.HTTPError{
			Status: http.StatusInternalServerError,
			Title:  "Token Generation Error.",
			Err:    errors.New("Could not generate access and refresh tokens."),
		}
	}

	tokens := CreateUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokens, nil
}
