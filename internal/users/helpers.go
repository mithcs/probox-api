package users

import (
	"context"
	"errors"

	"github.com/mithcs/probox-api/internal/globals"
	"golang.org/x/crypto/bcrypt"
)

func compareUserIDWithToken(ctx context.Context, id int64) error {
	uid, err := globals.GetUserIDFromContext(ctx)
	if err != nil {
		return errors.New("Could not get User ID from context.")
	}

	if uid != id {
		return errors.New("Not authorized to delete user.")
	}

	return nil
}

func deleteUserByID(ctx context.Context, id int64) error {
	return deleteUserByIDFromDB(ctx, id)
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

func hashPass(pass string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)

	return string(hashed), err
}

func storeUser(ctx context.Context, username string, password string, fullname string) (int64, error) {
	return storeUserInDB(ctx, username, password, fullname)
}

func generateTokens(userID int64) (CreateUserResponse, error) {
	accessToken, refreshToken, err := globals.GenerateAccessAndRefreshTokens(userID)

	tokens := CreateUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokens, err
}
