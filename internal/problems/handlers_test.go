package problems

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mithcs/probox-api/internal/globals"
)

func TestCreateProblem(t *testing.T) {
	t.Run("invalid title", func(t *testing.T) {
		problemReq := CreateProblemRequest{
			Title:       "invalid title",
			Description: "description of problem",
		}

		reqData, err := json.Marshal(problemReq)
		if err != nil {
			t.Errorf("got err %v, expected nil", err)
		}

		reqBody := strings.NewReader(string(reqData))

		req := httptest.NewRequest(http.MethodPost, "/problems", reqBody)
		res := httptest.NewRecorder()
		CreateProblem(res, req)

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
		wantDetails := "Invalid Title."

		if gotDetails != wantDetails {
			t.Errorf("got %q, expected %q", gotDetails, wantDetails)
		}
	})

	t.Run("invalid description", func(t *testing.T) {
		problemReq := CreateProblemRequest{
			Title:       "example problem",
			Description: "invalid description",
		}

		reqData, err := json.Marshal(problemReq)
		if err != nil {
			t.Errorf("got err %v, expected nil", err)
		}

		reqBody := strings.NewReader(string(reqData))

		req := httptest.NewRequest(http.MethodPost, "/problems", reqBody)
		res := httptest.NewRecorder()
		CreateProblem(res, req)

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
		wantDetails := "Invalid Description."

		if gotDetails != wantDetails {
			t.Errorf("got %q, expected %q", gotDetails, wantDetails)
		}
	})

	t.Run("valid problem", func(t *testing.T) {
		problemReq := CreateProblemRequest{
			Title:       "example problem",
			Description: "description of problem",
		}

		reqData, err := json.Marshal(problemReq)
		if err != nil {
			t.Errorf("got err %v, expected nil", err)
		}

		reqBody := strings.NewReader(string(reqData))

		req := httptest.NewRequest(http.MethodPost, "/problems", reqBody)
		res := httptest.NewRecorder()
		CreateProblem(res, req)

		gotCode := res.Code
		wantCode := http.StatusOK

		if gotCode != wantCode {
			t.Errorf("got status code %d, expected %d", gotCode, wantCode)
		}

		gotBody, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("got err %v, expected nil", err)
		}

		var problemRes CreateProblemResponse
		err = json.Unmarshal(gotBody, &problemRes)
		if err != nil {
			t.Errorf("got err %v, expected nil", err)
		}

		gotId := problemRes.Id
		wantId := 1

		if gotId != wantId {
			t.Errorf("got id %d, expected %d", gotId, wantId)
		}

		gotTitle := problemRes.Title
		wantTitle := "example problem"

		if gotTitle != wantTitle {
			t.Errorf("got title %q, expected %q", gotTitle, wantTitle)
		}

		gotDescription := problemRes.Description
		wantDescription := "description of problem"

		if gotDescription != wantDescription {
			t.Errorf("got description %q, expected %q", gotDescription, wantDescription)
		}
	})
}
