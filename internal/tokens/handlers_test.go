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

	t.Run("test status code", func(t *testing.T) {
		got := res.Code
		want := http.StatusOK

		if got != want {
			t.Errorf("got status code %d, expected %d", got, want)
		}
	})

	t.Run("test body", func(t *testing.T) {
		got, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected err to be nil, got %v", err)
		}

		want := "create new pair of access and refresh tokens from username/password"

		if string(got) != want {
			t.Errorf("got %q, expected %q", got, want)
		}
	})
}

func TestRefreshTokens(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/tokens", nil)
	res := httptest.NewRecorder()
	RefreshTokens(res, req)

	t.Run("test status code", func(t *testing.T) {
		got := res.Code
		want := http.StatusOK

		if got != want {
			t.Errorf("got status code %d, expected %d", got, want)
		}
	})

	t.Run("test body", func(t *testing.T) {
		got, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected err to be nil, got %v", err)
		}

		want := "create new pair of access and refresh tokens from refresh token"

		if string(got) != want {
			t.Errorf("got %q, expected %q", got, want)
		}
	})
}
