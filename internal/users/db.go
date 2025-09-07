package users

import (
	"context"
	"crypto/rand"

	"github.com/mithcs/probox-api/internal/db"
	"github.com/mithcs/probox-api/internal/globals"
)

func deleteUserByIDFromDB(ctx context.Context, id int64) error {
	random := rand.Text() + rand.Text()
	err := globals.Queries.DeleteUserByID(ctx, db.DeleteUserByIDParams{
		ID:       id,
		Password: random,
	})

	return err
}

func getUserByIDFromDB(ctx context.Context, id int64) (db.GetUserByIDRow, error) {
	userRows, err := globals.Queries.GetUserByID(ctx, id)

	return userRows, err
}

func getUserByUsernameFromDB(ctx context.Context, username string) (db.GetUserByUsernameRow, error) {
	userRows, err := globals.Queries.GetUserByUsername(ctx, username)

	return userRows, err
}

func getUsersFromDB(ctx context.Context) ([]db.GetUsersRow, error) {
	userRows, err := globals.Queries.GetUsers(ctx)

	return userRows, err
}

func storeUserInDB(ctx context.Context, username string, password string, fullname string) (int64, error) {
	userID, err := globals.Queries.CreateUser(ctx, db.CreateUserParams{
		Username: username,
		Password: password,
		FullName: fullname,
	})

	return userID, err
}
