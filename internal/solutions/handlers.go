package solutions

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mithcs/probox-api/internal/db"
	"github.com/mithcs/probox-api/internal/globals"
)

type GetSolutionsResponse struct {
	Id          int64  `json:"id"`
	ProblemId   int64  `json:"problemId"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateSolutionRequest struct {
	ProblemId   int64  `json:"problemId"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
type CreateSolutionResponse struct {
	Id          int64  `json:"id"`
	ProblemId   int64  `json:"problemId"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func GetSolutions(w http.ResponseWriter, r *http.Request) {
	solutions, err := getAllSolutions(r.Context())
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not get all solutions.",
		})

		return
	}

	response, err := json.Marshal(solutions)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not list all solutions.",
		})

		return
	}

	w.Write(response)
}

func getAllSolutions(ctx context.Context) ([]GetSolutionsResponse, error) {
	solutionRows, err := getSolutionsFromDb(ctx)
	if err != nil {
		return []GetSolutionsResponse{}, err
	}

	var solutions []GetSolutionsResponse
	for _, solutionRow := range solutionRows {
		solutions = append(solutions, GetSolutionsResponse{
			Id:          solutionRow.ID,
			ProblemId:   solutionRow.Problemid,
			Title:       solutionRow.Title,
			Description: solutionRow.Description,
		})
	}

	return solutions, nil
}

func getSolutionsFromDb(ctx context.Context) ([]db.Solution, error) {
	solutionRows, err := globals.Queries.GetSolutions(ctx)

	return solutionRows, err
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
		Id:          s.ID,
		ProblemId:   s.Problemid,
		Title:       s.Title,
		Description: s.Description,
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

	if err := validateTitle(solution.Title); err != nil {
		return err
	}

	if err := validateDescription(solution.Description); err != nil {
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

func validateTitle(title string) error {
	minLen := 3
	maxLen := 120

	if len(title) > minLen && len(title) < maxLen {
		return nil
	}

	return errors.New("Invalid title.")
}

func validateDescription(description string) error {
	minLen := 5
	maxLen := 20000

	if len(description) > minLen && len(description) < maxLen {
		return nil
	}

	return errors.New("Invalid description.")
}

func storeSolution(ctx context.Context, solution CreateSolutionRequest) (db.Solution, error) {
	s, err := globals.Queries.CreateSolution(ctx, db.CreateSolutionParams{
		Problemid:   solution.ProblemId,
		Title:       solution.Title,
		Description: solution.Description,
	})

	return s, err
}
