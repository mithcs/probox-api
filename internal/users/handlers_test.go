package users

import (
	"io"
	"net/http"
	"net/http/httptest"
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

	wantBody := "create new user and generate access/refresh tokens"

	if string(gotBody) != wantBody {
		t.Errorf("got %q, expected %q", gotBody, wantBody)
	}
}
