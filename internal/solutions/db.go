package solutions

import (
	"context"

	"github.com/mithcs/probox-api/internal/db"
	"github.com/mithcs/probox-api/internal/globals"
)

func getSolutionByIDFromDB(ctx context.Context, solutionID int64) (db.Solution, error) {
	solution, err := globals.Queries.GetSolutionByID(ctx, solutionID)

	return solution, err
}

func getSolutionsFromDB(ctx context.Context) ([]db.Solution, error) {
	solutionRows, err := globals.Queries.GetSolutions(ctx)

	return solutionRows, err
}

func getSolutionsByProblemIDFromDB(ctx context.Context, problemID int64) ([]db.Solution, error) {
	solutionRows, err := globals.Queries.GetSolutionsByProblemID(ctx, problemID)

	return solutionRows, err
}

func getFullNameByIDFromDB(ctx context.Context, id int64) (string, error) {
	return globals.Queries.GetFullNameByID(ctx, id)
}

func storeSolutionInDB(ctx context.Context, solution SolutionToStore) (db.Solution, error) {
	s, err := globals.Queries.CreateSolution(ctx, db.CreateSolutionParams{
		ProblemID:   solution.ProblemID,
		Title:       solution.Title,
		Description: solution.Description,
		OwnerID:     solution.OwnerID,
		OwnerName:   solution.OwnerName,
	})

	return s, err
}
