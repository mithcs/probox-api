package solutions

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/mithcs/probox-api/internal/globals"
)

func getSolutionByID(ctx context.Context, solutionID int64) (GetSolutionResponse, *globals.HTTPError) {
	p, err := getSolutionByIDFromDB(ctx, solutionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return GetSolutionResponse{}, &globals.HTTPError{
				Status: http.StatusNoContent,
				Title:  "No Content.",
				Err:    errors.New("No solution to show."),
			}
		}

		return GetSolutionResponse{}, &globals.HTTPError{
			Status: http.StatusInternalServerError,
			Title:  "Database Error.",
			Err:    errors.New("Could not get solution from database."),
		}
	}

	solution := GetSolutionResponse{
		ID:          p.ID,
		ProblemID:   p.ProblemID,
		Title:       p.Title,
		Description: p.Description,
		OwnerID:     p.OwnerID,
		OwnerName:   p.OwnerName,
		CreatedAt:   p.CreatedAt.String(),
	}

	return solution, nil
}

func getAllSolutions(ctx context.Context) ([]GetSolutionResponse, *globals.HTTPError) {
	solutionRows, err := getSolutionsFromDB(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return []GetSolutionResponse{}, &globals.HTTPError{
				Status: http.StatusNoContent,
				Title:  "No Content.",
				Err:    errors.New("No solution to show."),
			}
		}

		return []GetSolutionResponse{}, &globals.HTTPError{
			Status: http.StatusInternalServerError,
			Title:  "Database Error.",
			Err:    errors.New("Could not get solutions from database."),
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
			CreatedAt:   solutionRow.CreatedAt.String(),
		})
	}

	return solutions, nil
}

func getSolutionsByProblemID(ctx context.Context, problemID int64) ([]GetSolutionResponse, *globals.HTTPError) {
	solutionRows, err := getSolutionsByProblemIDFromDB(ctx, problemID)
	if err != nil {
		if err == sql.ErrNoRows {
			return []GetSolutionResponse{}, &globals.HTTPError{
				Status: http.StatusNoContent,
				Title:  "No Content.",
				Err:    errors.New("No solution to show."),
			}
		}

		return []GetSolutionResponse{}, &globals.HTTPError{
			Status: http.StatusInternalServerError,
			Title:  "Database Error.",
			Err:    errors.New("Could not get solutions for problem from database."),
		}
	}

	var solutions []GetSolutionResponse
	for _, solutionRow := range solutionRows {
		solutions = append(solutions, GetSolutionResponse{
			ID:          solutionRow.ID,
			Title:       solutionRow.Title,
			Description: solutionRow.Description,
			OwnerID:     solutionRow.OwnerID,
			OwnerName:   solutionRow.OwnerName,
			CreatedAt:   solutionRow.CreatedAt.String(),
		})
	}

	return solutions, nil
}

func storeSolution(ctx context.Context, solution SolutionToStore) (CreateSolutionResponse, *globals.HTTPError) {
	s, err := storeSolutionInDB(ctx, solution)
	if err != nil {
		return CreateSolutionResponse{}, &globals.HTTPError{
			Status: http.StatusInternalServerError,
			Title:  "Database Error.",
			Err:    errors.New("Could not store solution in database."),
		}
	}

	solutionRes := CreateSolutionResponse{
		ID:          s.ID,
		ProblemID:   s.ProblemID,
		Title:       s.Title,
		Description: s.Description,
		OwnerID:     s.OwnerID,
		OwnerName:   s.OwnerName,
		CreatedAt:   s.CreatedAt.String(),
	}

	return solutionRes, nil
}

func getFullNameByID(ctx context.Context, id int64) (string, *globals.HTTPError) {
	name, err := getFullNameByIDFromDB(ctx, id)
	if err != nil {
		return "", &globals.HTTPError{
			Status: http.StatusInternalServerError,
			Title:  "Database Error.",
			Err:    errors.New("Could not get full name from database."),
		}
	}

	return name, nil
}
