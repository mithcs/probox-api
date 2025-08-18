package problems

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

func TestCreateProblem(t *testing.T) {
	t.Run("invalid title", func(t *testing.T) {
		problemReq := CreateProblemRequest{
			Title:       "invalid title",
			Description: "description of problem",
		}

		reqData, err := json.Marshal(problemReq)
		am.AssertErrNil(t, err)

		reqBody := strings.NewReader(string(reqData))
		req := httptest.NewRequest(http.MethodPost, "/problems", reqBody)
		res := httptest.NewRecorder()
		CreateProblem(res, req)

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
		wantDetails := "Invalid Title."
		am.AssertString(t, gotDetails, wantDetails)
	})

	t.Run("invalid description", func(t *testing.T) {
		problemReq := CreateProblemRequest{
			Title:       "example problem",
			Description: "invalid description",
		}

		reqData, err := json.Marshal(problemReq)
		am.AssertErrNil(t, err)

		reqBody := strings.NewReader(string(reqData))
		req := httptest.NewRequest(http.MethodPost, "/problems", reqBody)
		res := httptest.NewRecorder()
		CreateProblem(res, req)

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
		wantDetails := "Invalid Description."
		am.AssertString(t, gotDetails, wantDetails)
	})

	t.Run("valid problem", func(t *testing.T) {
		problemReq := CreateProblemRequest{
			Title:       "example problem",
			Description: "description of problem",
		}

		reqData, err := json.Marshal(problemReq)
		am.AssertErrNil(t, err)

		reqBody := strings.NewReader(string(reqData))
		req := httptest.NewRequest(http.MethodPost, "/problems", reqBody)
		res := httptest.NewRecorder()
		CreateProblem(res, req)

		gotCode := res.Code
		wantCode := http.StatusOK
		am.AssertInt(t, gotCode, wantCode)

		gotBody, err := io.ReadAll(res.Body)
		am.AssertErrNil(t, err)

		var problemRes CreateProblemResponse
		err = json.Unmarshal(gotBody, &problemRes)
		am.AssertErrNil(t, err)

		gotId := problemRes.Id
		wantId := 1
		am.AssertInt(t, gotId, wantId)

		gotTitle := problemRes.Title
		wantTitle := "example problem"
		am.AssertString(t, gotTitle, wantTitle)

		gotDescription := problemRes.Description
		wantDescription := "description of problem"
		am.AssertString(t, gotDescription, wantDescription)
	})
}
