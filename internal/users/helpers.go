package users

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-chi/jwtauth/v5"
	"github.com/mithcs/probox-api/internal/db"
	"github.com/mithcs/probox-api/internal/globals"
	"golang.org/x/crypto/bcrypt"
)

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
