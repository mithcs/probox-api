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

// // TODO: refactor TestRefreshTokens()
//
// func TestRefreshTokens(t *testing.T) {
// 	t.Run("valid refresh token", func(t *testing.T) {
// 		userId := 1

// 		jwtString, err := globals.GenerateRefreshToken(userId)
// 		am.AssertErrNil(t, err)

// 		jwt, err := globals.RefreshTokenAuth.Decode(jwtString)
// 		am.AssertErrNil(t, err)

// 		ctx := context.WithValue(t.Context(), jwtauth.TokenCtxKey, jwt)
// 		ctx = context.WithValue(ctx, jwtauth.ErrorCtxKey, nil)

// 		req := httptest.NewRequestWithContext(ctx, http.MethodPut, "/tokens", nil)
// 		res := httptest.NewRecorder()
// 		RefreshTokens(res, req)

// 		gotCode := res.Code
// 		wantCode := http.StatusOK
// 		am.AssertInt(t, gotCode, wantCode)

// 		resBody, err := io.ReadAll(res.Body)
// 		am.AssertErrNil(t, err)

// 		var rTokenRes RefreshTokensResponse
// 		err = json.Unmarshal(resBody, &rTokenRes)
// 		am.AssertErrNil(t, err)

// 		assertJWT(t, rTokenRes.AccessToken, "access token")
// 		assertJWT(t, rTokenRes.RefreshToken, "refresh token")
// 	})
// }
