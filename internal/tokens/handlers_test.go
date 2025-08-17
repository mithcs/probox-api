package tokens

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/jwtauth/v5"
	"github.com/mithcs/probox-api/internal/globals"
	"github.com/mithcs/probox-api/internal/users"
)

func TestCreateTokens(t *testing.T) {
	t.Run("wrong credentials", func(t *testing.T) {
		userReq := users.CreateUserRequest{
			Username: "username",
			Password: "password",
		}

		reqData, err := json.Marshal(userReq)
		if err != nil {
			t.Errorf("got err %v, expected nil", err)
		}

		reqBody := strings.NewReader(string(reqData))

		req := httptest.NewRequest(http.MethodPost, "/tokens", reqBody)
		res := httptest.NewRecorder()
		CreateTokens(res, req)

		gotCode := res.Code
		wantCode := http.StatusBadRequest

		if gotCode != wantCode {
			t.Errorf("got status code %d, expected %d", gotCode, wantCode)
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("got err %v, expected nil", err)
		}

		var error globals.ErrorResponse

		err = json.Unmarshal(resBody, &error)
		if err != nil {
			t.Errorf("got err %v, expected nil", err)
		}

		gotTitle := error.Title
		wantTitle := "Bad Request."

		if gotTitle != wantTitle {
			t.Errorf("got %q, expected %q", gotTitle, wantTitle)
		}

		gotDetails := error.Details
		wantDetails := "Invalid Credentials."

		if gotDetails != wantDetails {
			t.Errorf("got %q, expected %q", gotDetails, wantDetails)
		}
	})

	t.Run("valid credentials", func(t *testing.T) {
		userReq := users.CreateUserRequest{
			Username: "example",
			Password: "correct horse battery staple",
		}

		reqData, err := json.Marshal(userReq)
		if err != nil {
			t.Errorf("got err %v, expected nil", err)
		}

		reqBody := strings.NewReader(string(reqData))

		req := httptest.NewRequest(http.MethodPost, "/tokens", reqBody)
		res := httptest.NewRecorder()
		CreateTokens(res, req)

		gotCode := res.Code
		wantCode := http.StatusOK

		if gotCode != wantCode {
			t.Errorf("got status code %d, expected %d", gotCode, wantCode)
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("got err %v, expected nil", err)
		}

		var tokens CreateTokensResponse

		err = json.Unmarshal(resBody, &tokens)
		if err != nil {
			t.Errorf("got err %v, expected nil", err)
		}

		assertJWT(t, tokens.AccessToken, "access token")
		assertJWT(t, tokens.RefreshToken, "refresh token")
	})
}

// actually, the verification of tokens takes place at router level (via middleware)
func TestRefreshTokens(t *testing.T) {
	t.Run("valid refresh token", func(t *testing.T) {
		userId := 1

		jwtString, err := globals.GenerateRefreshToken(userId)
		if err != nil {
			t.Errorf("got err %v, expected nil", err)
		}

		jwt, err := globals.RefreshTokenAuth.Decode(jwtString)
		if err != nil {
			t.Errorf("got err %v, expected nil", err)
		}

		ctx := context.WithValue(t.Context(), jwtauth.TokenCtxKey, jwt)
		ctx = context.WithValue(ctx, jwtauth.ErrorCtxKey, nil)

		req := httptest.NewRequestWithContext(ctx, http.MethodPut, "/tokens", nil)
		res := httptest.NewRecorder()
		RefreshTokens(res, req)

		gotCode := res.Code
		wantCode := http.StatusOK

		if gotCode != wantCode {
			t.Errorf("got status code %d, expected %d", gotCode, wantCode)
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("got err %v, expected nil", err)
		}

		var rTokenRes RefreshTokensResponse

		err = json.Unmarshal(resBody, &rTokenRes)
		if err != nil {
			t.Errorf("got err %v, expected nil", err)
		}

		assertJWT(t, rTokenRes.AccessToken, "access token")
		assertJWT(t, rTokenRes.RefreshToken, "refresh token")
	})
}

func assertJWT(t *testing.T, jwt string, tokenType string) {
	t.Helper()

	if !strings.HasPrefix(jwt, "eyJhbGciOi") ||
		strings.Count(jwt, ".") != 2 {
		t.Errorf("invalid %s", tokenType)
	}
}
