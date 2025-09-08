package tokens

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/jwtauth/v5"
	"github.com/mithcs/probox-api/internal/globals"
)

func generateTokens(userID int64) (CreateTokensResponse, *globals.HTTPError) {
	accessToken, refreshToken, err := globals.GenerateAccessAndRefreshTokens(userID)
	if err != nil {
		return CreateTokensResponse{}, &globals.HTTPError{
			Status: http.StatusInternalServerError,
			Title:  "Token Generation Error.",
			Err:    errors.New("Could not generate access and refresh tokens."),
		}
	}

	tokens := CreateTokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokens, nil
}

func retrieveUserID(ctx context.Context) (int64, *globals.HTTPError) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return -1, &globals.HTTPError{
			Status: http.StatusBadRequest,
			Title:  "Invalid Token.",
			Err:    errors.New("Could not get user id from token."),
		}
	}

	userID, err := strconv.ParseInt(fmt.Sprintf("%v", claims["uid"]), 10, 64)
	if err != nil {
		return -1, &globals.HTTPError{
			Status: http.StatusBadRequest,
			Title:  "Invalid Token.",
			Err:    errors.New("Could not parse user id."),
		}
	}

	return userID, nil
}
