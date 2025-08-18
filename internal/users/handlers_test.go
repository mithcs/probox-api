package users

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mithcs/probox-api/internal/globals"
	"github.com/mithcs/probox-api/pkg/am"
)

func TestCreateUser(t *testing.T) {
	t.Run("invalid username", func(t *testing.T) {
		userReq := CreateUserRequest{
			Username: "",
			Password: "correct horse battery staple",
		}

		reqData, err := json.Marshal(userReq)
		am.AssertErrNil(t, err)

		reqBody := strings.NewReader(string(reqData))
		req := httptest.NewRequest(http.MethodPost, "/users", reqBody)
		res := httptest.NewRecorder()
		CreateUser(res, req)

		gotCode := res.Code
		wantCode := http.StatusBadRequest
		am.AssertInt(t, gotCode, wantCode)

		resBody, err := io.ReadAll(res.Body)
		am.AssertErrNil(t, err)

		var error globals.ErrorResponse
		err = json.Unmarshal(resBody, &error)
		am.AssertErrNil(t, err)

		gotTitle := error.Title
		wantTitle := "Bad Request."
		am.AssertString(t, gotTitle, wantTitle)

		gotDetails := error.Details
		wantDetails := "Invalid username."
		am.AssertString(t, gotDetails, wantDetails)
	})

	t.Run("invalid password", func(t *testing.T) {
		userReq := CreateUserRequest{
			Username: "example",
			Password: "",
		}

		data, err := json.Marshal(userReq)
		am.AssertErrNil(t, err)

		body := strings.NewReader(string(data))
		req := httptest.NewRequest(http.MethodPost, "/users", body)
		res := httptest.NewRecorder()
		CreateUser(res, req)

		gotCode := res.Code
		wantCode := http.StatusBadRequest
		am.AssertInt(t, gotCode, wantCode)

		resBody, err := io.ReadAll(res.Body)
		am.AssertErrNil(t, err)

		var error globals.ErrorResponse
		err = json.Unmarshal(resBody, &error)
		am.AssertErrNil(t, err)

		gotTitle := error.Title
		wantTitle := "Bad Request."
		am.AssertString(t, gotTitle, wantTitle)

		gotDetails := error.Details
		wantDetails := "Password is insecure."
		am.AssertString(t, gotDetails, wantDetails)
	})

	t.Run("valid username and password", func(t *testing.T) {
		userReq := CreateUserRequest{
			Username: "example",
			Password: "correct horse battery staple",
		}

		data, err := json.Marshal(userReq)
		am.AssertErrNil(t, err)

		body := strings.NewReader(string(data))
		req := httptest.NewRequest(http.MethodPost, "/users", body)
		res := httptest.NewRecorder()
		CreateUser(res, req)

		gotCode := res.Code
		wantCode := http.StatusOK
		am.AssertInt(t, gotCode, wantCode)

		gotBody, err := io.ReadAll(res.Body)
		am.AssertErrNil(t, err)

		var userRes CreateUserResponse
		err = json.Unmarshal(gotBody, &userRes)
		am.AssertErrNil(t, err)

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
