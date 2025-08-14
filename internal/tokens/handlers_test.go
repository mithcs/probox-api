package tokens

import (
	"io"
	"net/http"
	"net/http/httptest"
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

	wantBody := "generate new pair of access and refresh tokens from refresh token"

	if string(gotBody) != wantBody {
		t.Errorf("got %q, expected %q", gotBody, wantBody)
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

	wantBody := "generate new access token using refresh token"
	if string(gotBody) != wantBody {
		t.Errorf("got %q, expected %q", gotBody, wantBody)
	}
}
