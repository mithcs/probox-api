package solutions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mithcs/probox-api/internal/db"
	"github.com/mithcs/probox-api/internal/globals"
)

type CreateSolutionRequest struct {
	ProblemId int64  `json:"problemId"`
	Solution  string `json:"solution"`
}
type CreateSolutionResponse struct {
	Id        int64  `json:"id"`
	ProblemId int64  `json:"problemId"`
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

	err = solution.validate(r.Context())
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusBadRequest,
			Title:   "Bad Request.",
			Details: err.Error(),
		})

		return
	}

	s, err := storeSolution(r.Context(), solution)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not store solution.",
		})

		return
	}

	solutionRes := CreateSolutionResponse{
		Id:        s.ID,
		ProblemId: s.Problemid,
		Solution:  s.Solution,
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

func (solution *CreateSolutionRequest) validate(ctx context.Context) error {
	if err := validateProblemId(ctx, solution.ProblemId); err != nil {
		return err
	}

	if err := validateSolution(solution.Solution); err != nil {
		return err
	}

	return nil
}

func validateProblemId(ctx context.Context, problemId int64) error {
	problem, _ := globals.Queries.GetProblemById(ctx, problemId)
	if problem.ID == 0 || problem.Title == "" || problem.Description == "" {
		return errors.New("Invalid Problem Id.")
	}

	return nil
}

func validateSolution(solution string) error {
	minLen := 5
	maxLen := 20000

	if len(solution) > minLen && len(solution) < maxLen {
		return nil
	}

	fmt.Print(solution)

	return errors.New("Invalid Solution.")
}

func storeSolution(ctx context.Context, solution CreateSolutionRequest) (db.Solution, error) {
	s, err := globals.Queries.CreateSolution(ctx, db.CreateSolutionParams{
		Problemid: solution.ProblemId,
		Solution:  solution.Solution,
	})

	return s, err
}
