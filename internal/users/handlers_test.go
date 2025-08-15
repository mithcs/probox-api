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
	t.Run("invalid username", func(t *testing.T) {
		userReq := CreateUserRequest{
			Username: "",
			Password: "correct horse battery staple",
		}
		data, err := json.Marshal(userReq)
		if err != nil {
			t.Errorf("got err %v, expected nil", err)
		}
		body := strings.NewReader(string(data))

		req := httptest.NewRequest(http.MethodPost, "/users", body)
		res := httptest.NewRecorder()
		CreateUser(res, req)

		gotCode := res.Code
		wantCode := http.StatusBadRequest

		if gotCode != wantCode {
			t.Errorf("got status code %d, expected %d", gotCode, wantCode)
		}
		// test body here
	})

	t.Run("invalid password", func(t *testing.T) {
		userReq := CreateUserRequest{
			Username: "example",
			Password: "",
		}
		data, err := json.Marshal(userReq)
		if err != nil {
			t.Errorf("got err %v, expected nil", err)
		}
		body := strings.NewReader(string(data))

		req := httptest.NewRequest(http.MethodPost, "/users", body)
		res := httptest.NewRecorder()
		CreateUser(res, req)

		gotCode := res.Code
		wantCode := http.StatusBadRequest

		if gotCode != wantCode {
			t.Errorf("got status code %d, expected %d", gotCode, wantCode)
		}
		// test body here
	})

	t.Run("valid username and password", func(t *testing.T) {
		userReq := CreateUserRequest{
			Username: "example",
			Password: "correct horse battery staple",
		}
		data, err := json.Marshal(userReq)
		if err != nil {
			t.Errorf("got err %v, expected nil", err)
		}
		body := strings.NewReader(string(data))

		req := httptest.NewRequest(http.MethodPost, "/users", body)
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

		assertJWT(t, userRes.AccessToken, "access token")
		assertJWT(t, userRes.RefreshToken, "refresh token")
	})
}

func assertJWT(t *testing.T, jwt string, tokenType string) {
	t.Helper()

	if !strings.HasPrefix(jwt, "eyJhbGciOi") ||
		strings.Count(jwt, ".") != 2 {
		t.Errorf("invalid %s", tokenType)
	}
}
