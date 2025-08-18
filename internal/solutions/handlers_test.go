package solutions

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

func TestCreateSolution(t *testing.T) {
	t.Run("invalid problem id", func(t *testing.T) {
		solutionReq := CreateSolutionRequest{
			ProblemId: 0,
			Solution:  "example solution",
		}

		reqData, err := json.Marshal(solutionReq)
		am.AssertErrNil(t, err)

		reqBody := strings.NewReader(string(reqData))

		req := httptest.NewRequest(http.MethodPost, "/solutions", reqBody)
		res := httptest.NewRecorder()
		CreateSolution(res, req)

		gotCode := res.Code
		wantCode := http.StatusBadRequest

		if gotCode != wantCode {
			t.Errorf("got status code %d, expected %d", gotCode, wantCode)
		}

		resBody, err := io.ReadAll(res.Body)
		am.AssertErrNil(t, err)

		var error globals.ErrorResponse
		err = json.Unmarshal(resBody, &error)
		am.AssertErrNil(t, err)

		gotTitle := error.Title
		wantTitle := "Bad Request."
		am.AssertString(t, gotTitle, wantTitle)

		gotDetails := error.Details
		wantDetails := "Invalid problemId."
		am.AssertString(t, gotDetails, wantDetails)
	})

	t.Run("invalid solution", func(t *testing.T) {
		solutionReq := CreateSolutionRequest{
			ProblemId: 1,
			Solution:  "solution",
		}

		reqData, err := json.Marshal(solutionReq)
		am.AssertErrNil(t, err)

		reqBody := strings.NewReader(string(reqData))
		req := httptest.NewRequest(http.MethodPost, "/solutions", reqBody)
		res := httptest.NewRecorder()
		CreateSolution(res, req)

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
		wantDetails := "Invalid solution."
		am.AssertString(t, gotDetails, wantDetails)
	})

	t.Run("valid problem id and solution", func(t *testing.T) {
		solutionReq := CreateSolutionRequest{
			ProblemId: 1,
			Solution:  "example solution",
		}

		reqData, err := json.Marshal(solutionReq)
		am.AssertErrNil(t, err)

		reqBody := strings.NewReader(string(reqData))
		req := httptest.NewRequest(http.MethodPost, "/solutions", reqBody)
		res := httptest.NewRecorder()
		CreateSolution(res, req)

		gotCode := res.Code
		wantCode := http.StatusOK
		am.AssertInt(t, gotCode, wantCode)

		gotBody, err := io.ReadAll(res.Body)
		am.AssertErrNil(t, err)

		var solutionRes CreateSolutionResponse
		err = json.Unmarshal(gotBody, &solutionRes)
		am.AssertErrNil(t, err)

		gotSolutionId := solutionRes.Id
		wantSolutionId := 1
		am.AssertInt(t, gotSolutionId, wantSolutionId)

		gotProblemId := solutionRes.ProblemId
		wantProblemId := 1
		am.AssertInt(t, gotProblemId, wantProblemId)

		gotSolution := solutionRes.Solution
		wantSolution := "example solution"
		am.AssertString(t, gotSolution, wantSolution)
	})
}
