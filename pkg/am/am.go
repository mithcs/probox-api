package am

import (
	"strings"
	"testing"
)

func AssertInt(t *testing.T, got int, want int) {
	t.Helper()

	if got != want {
		t.Errorf("got %d, expcted %d", got, want)
	}
}

func AssertString(t *testing.T, got string, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %q, expected %q", got, want)
	}
}

func AssertErrNil(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("got err %v, expected nil", err)
	}
}

func AssertErrNotNil(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Errorf("got err nil, expected not nil")
	}
}

func AssertAny(t *testing.T, got any, want any) {
	t.Helper()

	if got != want {
		t.Errorf("got %v, expected %v", got, want)
	}
}

func AssertJWT(t *testing.T, jwt string, tokenType string) {
	t.Helper()

	if !strings.HasPrefix(jwt, "eyJhbGciOi") ||
		strings.Count(jwt, ".") != 2 {
		t.Errorf("invalid %s", tokenType)
	}
}
