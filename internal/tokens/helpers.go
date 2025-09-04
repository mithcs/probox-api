package tokens

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-chi/jwtauth/v5"
	"github.com/mithcs/probox-api/internal/globals"
	"golang.org/x/crypto/bcrypt"
)

func generateTokens(userID int64) (CreateTokensResponse, error) {
	accessToken, refreshToken, err := globals.GenerateAccessAndRefreshTokens(userID)

	tokens := CreateTokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokens, err
}

func retrieveUserID(ctx context.Context) (int64, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return 0, err
	}

	userID, err := strconv.ParseInt(fmt.Sprintf("%v", claims["uid"]), 10, 64)

	return userID, err
}

func verifyCredentials(ctx context.Context, username string, password string) (int64, error) {
	user, err := globals.Queries.GetCredentialsByUsername(ctx, username)
	if err == sql.ErrNoRows {
		return 0, errors.New("Incorrect Username.")
	}
	if err != nil {
		return 0, errors.New("Invalid Credentials.")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return 0, errors.New("Incorrect Password.")
	}

	return user.ID, nil
}
