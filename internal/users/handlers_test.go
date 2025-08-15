package users

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateUser(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/users", nil)
	res := httptest.NewRecorder()
	CreateUser(res, req)

	gotCode := res.Code
	wantCode := http.StatusOK

	if gotCode != wantCode {
		t.Errorf("got status code %d, expected %d", gotCode, wantCode)
	}

	gotBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("got err %v, expected nil", err)
	}

	var userRes CreateUserResponse
	err = json.Unmarshal(gotBody, &userRes)
	if err != nil {
		t.Errorf("got err %v, expected nil", err)
	}

	// checks only `{ "alg":` part to verify whether its JWT
	if !strings.HasPrefix(userRes.AccessToken, "eyJhbGciOi") {
		t.Error("invalid access token")
	}

	if !strings.HasPrefix(userRes.RefreshToken, "eyJhbGciOi") {
		t.Error("invalid refresh token")
	}
}
