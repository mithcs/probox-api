package solutions

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/mithcs/probox-api/internal/globals"
)

type CreateSolutionRequest struct {
	ProblemId int    `json:"problemId"`
	Solution  string `json:"solution"`
}
type CreateSolutionResponse struct {
	Id        int    `json:"id"`
	ProblemId int    `json:"problemId"`
	Solution  string `json:"solution"`
}

func CreateSolution(w http.ResponseWriter, r *http.Request) {
	var solution CreateSolutionRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := globals.ReturnErrorResponse(
			"Server Error.",
			"Could not read the body",
		)
		w.Write(response)

		return
	}

	err = json.Unmarshal(body, &solution)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := globals.ReturnErrorResponse(
			"Bad Request.",
			"Could not parse body.",
		)
		w.Write(response)

		return
	}

	err = solution.validate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := globals.ReturnErrorResponse(
			"Bad Request.",
			err.Error(),
		)
		w.Write(response)

		return
	}

	// store solution in db
	// and get solutionId
	solutionId := 1

	solutionRes := CreateSolutionResponse{
		Id:        solutionId,
		ProblemId: solution.ProblemId,
		Solution:  solution.Solution,
	}

	response, err := json.Marshal(solutionRes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := globals.ReturnErrorResponse(
			"Server Error.",
			"Could not return created problem.",
		)
		w.Write(response)

		return
	}

	w.Write(response)
}

func (solution *CreateSolutionRequest) validate() error {
	// TODO

	if solution.ProblemId == 0 {
		return errors.New("Invalid problemId.")
	}
	if solution.Solution == "solution" {
		return errors.New("Invalid solution.")
	}

	return nil
}
