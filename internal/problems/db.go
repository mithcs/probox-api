package problems

import (
	"context"
	"database/sql"

	"github.com/mithcs/probox-api/internal/db"
	"github.com/mithcs/probox-api/internal/globals"
)

func getProblemByIDFromDB(ctx context.Context, problemID int64) (db.Problem, error) {
	problem, err := globals.Queries.GetProblemByID(ctx, problemID)

	return problem, err
}

func getProblemsFromDB(ctx context.Context) ([]db.Problem, error) {
	problemRows, err := globals.Queries.GetProblems(ctx)

	return problemRows, err
}

func getFullNameByIDFromDB(ctx context.Context, id int64) (string, error) {
	return globals.Queries.GetFullNameByID(ctx, id)
}

func getSolutionCountForProblemFromDB(ctx context.Context, id int64) (int64, error) {
	return globals.Queries.GetSolutionCountByProblemID(ctx, id)
}

func deleteProblemByIDFromDB(ctx context.Context, id int64) error {
	return globals.Queries.DeleteProblemByID(ctx, id)
}

func storeProblemInDB(ctx context.Context, problem ProblemToStore) (db.Problem, error) {
	p, err := globals.Queries.CreateProblem(ctx, db.CreateProblemParams{
		Title:       problem.Title,
		Description: problem.Description,
		OwnerID:     problem.OwnerID,
		OwnerName:   problem.OwnerName,
	})

	return p, err
}

func setAcceptedSolutionIDInDB(ctx context.Context, pid int64, sid int64) error {
	solution := sql.NullInt64{Int64: sid, Valid: true}

	return globals.Queries.SetAcceptedSolutionID(ctx, db.SetAcceptedSolutionIDParams{
		ID:                 pid,
		AcceptedSolutionID: solution,
	})
}

func unsetAcceptedSolutionIDInDB(ctx context.Context, pid int64) error {
	return globals.Queries.UnsetAcceptedSolutionID(ctx, pid)
}
