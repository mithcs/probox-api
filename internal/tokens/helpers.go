package tokens

import (
	"errors"
	"net/http"

	"github.com/mithcs/probox-api/internal/globals"
)

func generateTokens(userID int64) (CreateTokensResponse, *globals.HTTPError) {
	accessToken, refreshToken, err := globals.GenerateTokens(userID)
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
