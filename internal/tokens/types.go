package tokens

import "context"

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

func (req *CreateTokensRequest) Verify(ctx context.Context) (int64, error) {
	userID, err := verifyCredentials(ctx, req.Username, req.Password)

	return userID, err
}
