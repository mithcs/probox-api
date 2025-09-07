package problems

import (
	"context"
	"errors"

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

func canDeleteProblem(ctx context.Context, problem GetProblemResponse) error {
	solutionCount, err := getSolutionCountForProblem(ctx, problem.ID)
	if err != nil {
		return errors.New("Cannot get solution count for problem.")
	}

	if solutionCount != 0 {
		return errors.New("Problem does not have 0 Solutions. Cannot Delete.")
	}

	return nil
}

func getSolutionCountForProblem(ctx context.Context, id int64) (int64, error) {
	return getSolutionCountForProblemFromDB(ctx, id)
}

func getSolutionCountForProblemFromDB(ctx context.Context, id int64) (int64, error) {
	return globals.Queries.GetSolutionCountByProblemID(ctx, id)
}

func compareUserIDWithToken(ctx context.Context, id int64) error {
	uid, err := globals.GetUserIDFromContext(ctx)
	if err != nil {
		return errors.New("Could not get User ID from context.")
	}

	if uid != id {
		return errors.New("Not authorized to delete user.")
	}

	return nil
}

func deleteProblemByID(ctx context.Context, id int64) error {
	return deleteProblemByIDFromDB(ctx, id)
}

func deleteProblemByIDFromDB(ctx context.Context, id int64) error {
	return globals.Queries.DeleteProblemByID(ctx, id)
}
