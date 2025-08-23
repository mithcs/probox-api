package users

import (
	"encoding/json"
	"io"
	"strings"
	"testing"

	"github.com/mithcs/probox-api/pkg/am"
)

func TestParseBody(t *testing.T) {
	userReq := CreateUserRequest{
		Username: "username",
		Password: "correct horse battery staple",
	}
	userReqJson, err := json.Marshal(userReq)
	am.AssertErrNil(t, err)

	bodyReader := strings.NewReader(string(userReqJson))
	body := io.NopCloser(bodyReader)

	got, err := parseBody(body)
	am.AssertErrNil(t, err)

	want := userReq

	am.AssertAny(t, got, want)
}

func TestGenerateTokens(t *testing.T) {
	userId := int64(1)

	tokens, err := generateTokens(userId)
	am.AssertErrNil(t, err)

	am.AssertJWT(t, tokens.AccessToken, "access token")
	am.AssertJWT(t, tokens.RefreshToken, "refresh token")
}

func TestValidateUser(t *testing.T) {
	t.Run("valid username and password", func(t *testing.T) {
		user := CreateUserRequest{
			Username: "username",
			Password: "correct horse battery staple",
		}

		err := user.validate()
		am.AssertErrNil(t, err)
	})

	t.Run("invalid username", func(t *testing.T) {
		user := CreateUserRequest{
			Username: "jj",
			Password: "correct horse battery staple",
		}

		err := user.validate()
		am.AssertErrNotNil(t, err)
	})

	t.Run("invalid password", func(t *testing.T) {
		user := CreateUserRequest{
			Username: "username",
			Password: "aaaaaaaa",
		}

		err := user.validate()
		am.AssertErrNotNil(t, err)
	})
}
