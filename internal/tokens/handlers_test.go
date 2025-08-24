package tokens

import (
	"testing"

	"github.com/mithcs/probox-api/pkg/am"
)

func TestGenerateTokens(t *testing.T) {
	userId := int64(1)

	tokens, err := generateTokens(userId)
	am.AssertErrNil(t, err)

	am.AssertJWT(t, tokens.AccessToken, "access token")
	am.AssertJWT(t, tokens.RefreshToken, "refresh token")
}
