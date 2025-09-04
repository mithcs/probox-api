package solutions

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mithcs/probox-api/internal/db"
	"github.com/mithcs/probox-api/internal/globals"
)

func GetSolution(w http.ResponseWriter, r *http.Request) {
	solutionID, err := strconv.ParseInt(chi.URLParam(r, "solution_id"), 10, 64)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusBadRequest,
			Title:   "Bad Request.",
			Details: "Invalid Solution ID.",
		})

		return
	}

	solution, err := getSolutionByID(r.Context(), solutionID)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not get solution.",
		})

		return
	}

	response, err := json.Marshal(solution)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not list solution.",
		})

		return
	}

	w.Write(response)
}

func getSolutionByID(ctx context.Context, solutionID int64) (GetSolutionResponse, error) {
	p, err := getSolutionByIDFromDB(ctx, solutionID)

	solution := GetSolutionResponse{
		ID:          p.ID,
		ProblemID:   p.ProblemID,
		Title:       p.Title,
		Description: p.Description,
	}

	return solution, err
}

func getSolutionByIDFromDB(ctx context.Context, solutionID int64) (db.Solution, error) {
	solution, err := globals.Queries.GetSolutionByID(ctx, solutionID)

	return solution, err
}

func GetSolutions(w http.ResponseWriter, r *http.Request) {
	if problemID := r.URL.Query().Get("problem_id"); problemID != "" {
		GetSolutionsForProblem(w, r, problemID)
		return
	}

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

func getAllSolutions(ctx context.Context) ([]GetSolutionResponse, error) {
	solutionRows, err := getSolutionsFromDB(ctx)
	if err != nil {
		return []GetSolutionResponse{}, err
	}

	var solutions []GetSolutionResponse
	for _, solutionRow := range solutionRows {
		solutions = append(solutions, GetSolutionResponse{
			ID:          solutionRow.ID,
			ProblemID:   solutionRow.ProblemID,
			Title:       solutionRow.Title,
			Description: solutionRow.Description,
		})
	}

	return solutions, nil
}

func getSolutionsFromDB(ctx context.Context) ([]db.Solution, error) {
	solutionRows, err := globals.Queries.GetSolutions(ctx)

	return solutionRows, err
}

func GetSolutionsForProblem(w http.ResponseWriter, r *http.Request, problemID string) {
	pid, err := strconv.ParseInt(problemID, 10, 64)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusBadRequest,
			Title:   "Bad Request.",
			Details: "Invalid Problem ID.",
		})

		return
	}

	solutions, err := getSolutionsByProblemID(r.Context(), pid)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusBadRequest,
			Title:   "Bad Request.",
			Details: err.Error(),
		})

		return
	}

	if len(solutions) == 0 {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusNoContent,
			Title:   "No Content.",
			Details: "There are no solutions to list.",
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

func getSolutionsByProblemID(ctx context.Context, problemID int64) ([]GetSolutionsForProblemResponse, error) {
	solutionRows, err := getSolutionsByProblemIDFromDB(ctx, problemID)
	if err != nil {
		return []GetSolutionsForProblemResponse{}, err
	}

	var solutions []GetSolutionsForProblemResponse
	for _, solutionRow := range solutionRows {
		solutions = append(solutions, GetSolutionsForProblemResponse{
			ID:          solutionRow.ID,
			Title:       solutionRow.Title,
			Description: solutionRow.Description,
		})
	}

	return solutions, nil
}

func getSolutionsByProblemIDFromDB(ctx context.Context, problemID int64) ([]db.Solution, error) {
	solutionRows, err := globals.Queries.GetSolutionsByProblemID(ctx, problemID)

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

	err = solution.Validate(r.Context())
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
		ID:          s.ID,
		ProblemID:   s.ProblemID,
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

func storeSolution(ctx context.Context, solution CreateSolutionRequest) (db.Solution, error) {
	s, err := globals.Queries.CreateSolution(ctx, db.CreateSolutionParams{
		ProblemID:   solution.ProblemID,
		Title:       solution.Title,
		Description: solution.Description,
	})

	return s, err
}
