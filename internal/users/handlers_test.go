package users

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

func TestValidateUsername(t *testing.T) {
	t.Run("valid username", func(t *testing.T) {
		username := "valid"
		err := validateUsername(username)
		am.AssertErrNil(t, err)
	})

	t.Run("invalid username", func(t *testing.T) {
		username := "1nval1d"
		err := validateUsername(username)
		am.AssertErrNotNil(t, err)
	})
}

func TestValidatePassword(t *testing.T) {
	t.Run("valid password", func(t *testing.T) {
		password := "correct horse battery staple"
		err := validatePassword(password)
		am.AssertErrNil(t, err)
	})

	t.Run("invalid password", func(t *testing.T) {
		password := "insecure"
		err := validatePassword(password)
		am.AssertErrNotNil(t, err)
	})
}
