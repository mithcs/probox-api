package problems

import (
	"context"

	"github.com/mithcs/probox-api/internal/db"
	"github.com/mithcs/probox-api/internal/globals"
)

func getProblemByID(ctx context.Context, problemID int64) (GetProblemResponse, error) {
	p, err := getProblemByIDFromDB(ctx, problemID)

	problem := GetProblemResponse{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		OwnerID:     p.OwnerID,
		OwnerName:   p.OwnerName,
	}

	return problem, err
}

func getProblemByIDFromDB(ctx context.Context, problemID int64) (db.Problem, error) {
	problem, err := globals.Queries.GetProblemByID(ctx, problemID)

	return problem, err
}

func getAllProblems(ctx context.Context) ([]GetProblemResponse, error) {
	problemRows, err := getProblemsFromDB(ctx)
	if err != nil {
		return []GetProblemResponse{}, err
	}

	var problems []GetProblemResponse
	for _, problemRow := range problemRows {
		problems = append(problems, GetProblemResponse{
			ID:          problemRow.ID,
			Title:       problemRow.Title,
			Description: problemRow.Description,
			OwnerID:     problemRow.OwnerID,
			OwnerName:   problemRow.OwnerName,
		})
	}

	return problems, nil
}

func getProblemsFromDB(ctx context.Context) ([]db.Problem, error) {
	problemRows, err := globals.Queries.GetProblems(ctx)

	return problemRows, err
}

func storeProblem(ctx context.Context, problem ProblemToStore) (db.Problem, error) {
	p, err := globals.Queries.CreateProblem(ctx, db.CreateProblemParams{
		Title:       problem.Title,
		Description: problem.Description,
		OwnerID:     problem.OwnerID,
		OwnerName:   problem.OwnerName,
	})

	return p, err
}

func getFullNameByID(ctx context.Context, id int64) (string, error) {
	return getFullNameByIDFromDB(ctx, id)
}

func getFullNameByIDFromDB(ctx context.Context, id int64) (string, error) {
	return globals.Queries.GetFullNameByID(ctx, id)
}
