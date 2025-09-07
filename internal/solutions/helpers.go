package solutions

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mithcs/probox-api/internal/db"
	"github.com/mithcs/probox-api/internal/globals"
)

func getSolutionByID(ctx context.Context, solutionID int64) (GetSolutionResponse, *globals.HTTPError) {
	p, err := getSolutionByIDFromDB(ctx, solutionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return GetSolutionResponse{}, &globals.HTTPError{
				Status:      http.StatusNoContent,
				Title:       "No Content.",
				Description: "No solution to show.",
			}
		}

		return GetSolutionResponse{}, &globals.HTTPError{
			Status:      http.StatusInternalServerError,
			Title:       "Database Error.",
			Description: "Could not get solution from database.",
		}
	}

	solution := GetSolutionResponse{
		ID:          p.ID,
		ProblemID:   p.ProblemID,
		Title:       p.Title,
		Description: p.Description,
		OwnerID:     p.OwnerID,
		OwnerName:   p.OwnerName,
	}

	return solution, nil
}

func getAllSolutions(ctx context.Context) ([]GetSolutionResponse, *globals.HTTPError) {
	solutionRows, err := getSolutionsFromDB(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return []GetSolutionResponse{}, &globals.HTTPError{
				Status:      http.StatusNoContent,
				Title:       "No Content.",
				Description: "No solution to show.",
			}
		}

		return []GetSolutionResponse{}, &globals.HTTPError{
			Status:      http.StatusInternalServerError,
			Title:       "Database Error.",
			Description: "Could not get solutions from database.",
		}
	}

	var solutions []GetSolutionResponse
	for _, solutionRow := range solutionRows {
		solutions = append(solutions, GetSolutionResponse{
			ID:          solutionRow.ID,
			ProblemID:   solutionRow.ProblemID,
			Title:       solutionRow.Title,
			Description: solutionRow.Description,
			OwnerID:     solutionRow.OwnerID,
			OwnerName:   solutionRow.OwnerName,
		})
	}

	return solutions, nil
}

func getSolutionsByProblemID(ctx context.Context, problemID int64) ([]GetSolutionsForProblemResponse, *globals.HTTPError) {
	solutionRows, err := getSolutionsByProblemIDFromDB(ctx, problemID)
	if err != nil {
		if err == sql.ErrNoRows {
			return []GetSolutionsForProblemResponse{}, &globals.HTTPError{
				Status:      http.StatusNoContent,
				Title:       "No Content.",
				Description: "No solution to show.",
			}
		}

		return []GetSolutionsForProblemResponse{}, &globals.HTTPError{
			Status:      http.StatusInternalServerError,
			Title:       "Database Error.",
			Description: "Could not get solutions for problem from database.",
		}
	}

	var solutions []GetSolutionsForProblemResponse
	for _, solutionRow := range solutionRows {
		solutions = append(solutions, GetSolutionsForProblemResponse{
			ID:          solutionRow.ID,
			Title:       solutionRow.Title,
			Description: solutionRow.Description,
			OwnerID:     solutionRow.OwnerID,
			OwnerName:   solutionRow.OwnerName,
		})
	}

	return solutions, nil
}

func storeSolution(ctx context.Context, solution SolutionToStore) (db.Solution, *globals.HTTPError) {
	s, err := storeSolutionInDB(ctx, solution)
	if err != nil {
		return db.Solution{}, &globals.HTTPError{
			Status:      http.StatusInternalServerError,
			Title:       "Database Error.",
			Description: "Could not store solution in database.",
		}
	}

	return s, nil
}

func getFullNameByID(ctx context.Context, id int64) (string, *globals.HTTPError) {
	name, err := getFullNameByIDFromDB(ctx, id)
	if err != nil {
		return "", &globals.HTTPError{
			Status:      http.StatusInternalServerError,
			Title:       "Database Error.",
			Description: "Could not full name from database.",
		}
	}

	return name, nil
}

func getIDFromRequest(r *http.Request) (int64, *globals.HTTPError) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		return -1, &globals.HTTPError{
			Status:      http.StatusBadRequest,
			Title:       "Invalid ID.",
			Description: "Could not parse id.",
		}
	}

	return id, nil
}
