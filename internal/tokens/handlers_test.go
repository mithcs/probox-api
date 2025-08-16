package tokens

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

func TestRefreshTokens(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/tokens", nil)
	res := httptest.NewRecorder()
	RefreshTokens(res, req)

	gotCode := res.Code
	wantCode := http.StatusOK

	if gotCode != wantCode {
		t.Errorf("got status code %d, expected %d", gotCode, wantCode)
	}

	gotBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("got err %v, expected nil", err)
	}

	var rTokenRes RefreshTokensResponse
	err = json.Unmarshal(gotBody, &rTokenRes)
	if err != nil {
		t.Errorf("got err %v, expected nil", err)
	}

	// checks only `{ "alg":` part to verify whether its JWT
	if !strings.HasPrefix(rTokenRes.AccessToken, "eyJhbGciOi") {
		t.Error("invalid access token")
	}

	if !strings.HasPrefix(rTokenRes.RefreshToken, "eyJhbGciOi") {
		t.Error("invalid refresh token")
	}
}

func assertJWT(t *testing.T, jwt string, tokenType string) {
	t.Helper()

	if !strings.HasPrefix(jwt, "eyJhbGciOi") ||
		strings.Count(jwt, ".") != 2 {
		t.Errorf("invalid %s", tokenType)
	}
}
