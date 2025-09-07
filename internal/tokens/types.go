package tokens

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/mithcs/probox-api/internal/globals"
	"golang.org/x/crypto/bcrypt"
)

type CreateTokensRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateTokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (req *CreateTokensRequest) Verify(ctx context.Context) (int64, *globals.HTTPError) {
	userID, err := verifyCredentials(ctx, req.Username, req.Password)
	if err != nil {
		return -1, err
	}

	return userID, nil
}

func verifyCredentials(ctx context.Context, username string, password string) (int64, *globals.HTTPError) {
	user, err := globals.Queries.GetCredentialsByUsername(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, &globals.HTTPError{
				Status:      http.StatusBadRequest,
				Title:       "Incorrect Username.",
				Description: "No user exists with username.",
			}
		}

		return -1, &globals.HTTPError{
			Status:      http.StatusInternalServerError,
			Title:       "Database Error.",
			Description: "Could not get user information from database.",
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return -1, &globals.HTTPError{
			Status:      http.StatusBadRequest,
			Title:       "Incorrect Password.",
			Description: "Password does not match.",
		}
	}

	return user.ID, nil
}
