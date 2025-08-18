package solutions

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mithcs/probox-api/internal/globals"
)

func TestCreateSolution(t *testing.T) {
	t.Run("invalid problem id", func(t *testing.T) {
		solutionReq := CreateSolutionRequest{
			ProblemId: 0,
			Solution:  "example solution",
		}

		reqData, err := json.Marshal(solutionReq)
		if err != nil {
			t.Errorf("got err %v, expected nil", err)
		}

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
		wantDetails := "Invalid problemId."

		if gotDetails != wantDetails {
			t.Errorf("got %q, expected %q", gotDetails, wantDetails)
		}
	})

	t.Run("invalid solution", func(t *testing.T) {
		solutionReq := CreateSolutionRequest{
			ProblemId: 1,
			Solution:  "solution",
		}

		reqData, err := json.Marshal(solutionReq)
		if err != nil {
			t.Errorf("got err %v, expected nil", err)
		}

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
		wantDetails := "Invalid solution."

		if gotDetails != wantDetails {
			t.Errorf("got %q, expected %q", gotDetails, wantDetails)
		}
	})

	t.Run("valid problem id and solution", func(t *testing.T) {
		solutionReq := CreateSolutionRequest{
			ProblemId: 1,
			Solution:  "example solution",
		}

		reqData, err := json.Marshal(solutionReq)
		if err != nil {
			t.Errorf("got err %v, expected nil", err)
		}

		reqBody := strings.NewReader(string(reqData))

		req := httptest.NewRequest(http.MethodPost, "/solutions", reqBody)
		res := httptest.NewRecorder()
		CreateSolution(res, req)

		gotCode := res.Code
		wantCode := http.StatusOK

		if gotCode != wantCode {
			t.Errorf("got status code %d, expected %d", gotCode, wantCode)
		}

		gotBody, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("got err %v, expected nil", err)
		}

		var solutionRes CreateSolutionResponse
		err = json.Unmarshal(gotBody, &solutionRes)
		if err != nil {
			t.Errorf("got err %v, expected nil", err)
		}

		gotSolutionId := solutionRes.Id
		wantSolutionId := 1

		if gotSolutionId != wantSolutionId {
			t.Errorf("got solution id %d, expected %d", gotSolutionId, wantSolutionId)
		}

		gotProblemId := solutionRes.ProblemId
		wantProblemId := 1

		if gotProblemId != wantProblemId {
			t.Errorf("got problem id %d, expected %d", gotProblemId, wantProblemId)
		}

		gotSolution := solutionRes.Solution
		wantSolution := "example solution"

		if gotSolution != wantSolution {
			t.Errorf("got solution %q, expected %q", gotSolution, wantSolution)
		}
	})
}
