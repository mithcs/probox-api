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

	assertJWT(t, tokens.AccessToken, "access token")
	assertJWT(t, tokens.RefreshToken, "refresh token")
}

func TestValidateUser(t *testing.T) {
	user := CreateUserRequest{
		Username: "username",
		Password: "correct horse battery staple",
	}

	err := user.validate()
	am.AssertErrNil(t, err)
}

func assertJWT(t *testing.T, jwt string, tokenType string) {
	t.Helper()

	if !strings.HasPrefix(jwt, "eyJhbGciOi") ||
		strings.Count(jwt, ".") != 2 {
		t.Errorf("invalid %s", tokenType)
	}
}
