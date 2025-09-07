package solutions

import (
	"context"

	"github.com/mithcs/probox-api/internal/db"
)

func getSolutionByID(ctx context.Context, solutionID int64) (GetSolutionResponse, error) {
	p, err := getSolutionByIDFromDB(ctx, solutionID)

	solution := GetSolutionResponse{
		ID:          p.ID,
		ProblemID:   p.ProblemID,
		Title:       p.Title,
		Description: p.Description,
		OwnerID:     p.OwnerID,
		OwnerName:   p.OwnerName,
	}

	return solution, err
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
			OwnerID:     solutionRow.OwnerID,
			OwnerName:   solutionRow.OwnerName,
		})
	}

	return solutions, nil
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
			OwnerID:     solutionRow.OwnerID,
			OwnerName:   solutionRow.OwnerName,
		})
	}

	return solutions, nil
}

func storeSolution(ctx context.Context, solution SolutionToStore) (db.Solution, error) {
	return storeSolutionInDB(ctx, solution)
}

func getFullNameByID(ctx context.Context, id int64) (string, error) {
	return getFullNameByIDFromDB(ctx, id)
}
