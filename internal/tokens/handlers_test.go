package tokens

import (
	"encoding/json"
	"io"
	"strings"
	"testing"

	"github.com/mithcs/probox-api/pkg/am"
)

func TestParseBody(t *testing.T) {
	tokenReq := CreateTokensRequest{
		Username: "username",
		Password: "password",
	}

	tokenReqJson, err := json.Marshal(tokenReq)
	am.AssertErrNil(t, err)

	bodyReader := strings.NewReader(string(tokenReqJson))
	body := io.NopCloser(bodyReader)

	got, err := parseBody(body)
	am.AssertErrNil(t, err)

	want := tokenReq

	am.AssertAny(t, got, want)
}

func TestGenerateTokens(t *testing.T) {
	userId := int64(1)

	tokens, err := generateTokens(userId)
	am.AssertErrNil(t, err)

	am.AssertJWT(t, tokens.AccessToken, "access token")
	am.AssertJWT(t, tokens.RefreshToken, "refresh token")
}
