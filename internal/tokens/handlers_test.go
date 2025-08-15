package tokens

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateTokens(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/tokens", nil)
	res := httptest.NewRecorder()
	CreateTokens(res, req)

	gotCode := res.Code
	wantCode := http.StatusOK

	if gotCode != wantCode {
		t.Errorf("got status code %d, expected %d", gotCode, wantCode)
	}

	gotBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("got err %v, expected nil", err)
	}

	var tokens CreateTokensResponse
	err = json.Unmarshal(gotBody, &tokens)
	if err != nil {
		t.Errorf("got err %v, expected nil", err)
	}

	// checks only `{ "alg":` part to verify whether its JWT
	if !strings.HasPrefix(tokens.AccessToken, "eyJhbGciOi") {
		t.Error("invalid access token")
	}

	if !strings.HasPrefix(tokens.RefreshToken, "eyJhbGciOi") {
		t.Error("invalid refresh token")
	}
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
