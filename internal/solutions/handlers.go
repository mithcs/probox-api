package solutions

import (
	"encoding/json"
	"errors"
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
	solution, err := globals.ParseBody[CreateSolutionRequest](r.Body)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusBadRequest,
			Title:   "Bad Request.",
			Details: "Could not parse body.",
		})

		return
	}

	err = solution.validate()
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusBadRequest,
			Title:   "Bad Request.",
			Details: err.Error(),
		})

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
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not return created problem.",
		})

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
