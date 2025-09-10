package problems

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/mithcs/probox-api/internal/globals"
)

func getProblemByID(ctx context.Context, problemID int64) (GetProblemResponse, *globals.HTTPError) {
	p, err := getProblemByIDFromDB(ctx, problemID)
	if err != nil {
		if err == sql.ErrNoRows {
			return GetProblemResponse{}, &globals.HTTPError{
				Status: http.StatusNoContent,
				Title:  "No Content.",
				Err:    errors.New("No problem to show."),
			}
		}

		return GetProblemResponse{}, &globals.HTTPError{
			Status: http.StatusInternalServerError,
			Title:  "Database Error.",
			Err:    errors.New("Could not get problem from database."),
		}
	}

	problem := GetProblemResponse{
		ID:                 p.ID,
		Title:              p.Title,
		Description:        p.Description,
		OwnerID:            p.OwnerID,
		OwnerName:          p.OwnerName,
		AcceptedSolutionID: p.AcceptedSolutionID.Int64,
		CreatedAt:          p.CreatedAt,
	}

	return problem, nil
}

func getAllProblems(ctx context.Context) ([]GetProblemResponse, *globals.HTTPError) {
	problemRows, err := getProblemsFromDB(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return []GetProblemResponse{}, &globals.HTTPError{
				Status: http.StatusNoContent,
				Title:  "No Content.",
				Err:    errors.New("No problems to show."),
			}
		}

		return []GetProblemResponse{}, &globals.HTTPError{
			Status: http.StatusInternalServerError,
			Title:  "Database Error.",
			Err:    errors.New("Could not get problems from database."),
		}
	}

	var problems []GetProblemResponse
	for _, problemRow := range problemRows {
		problems = append(problems, GetProblemResponse{
			ID:                 problemRow.ID,
			Title:              problemRow.Title,
			Description:        problemRow.Description,
			OwnerID:            problemRow.OwnerID,
			OwnerName:          problemRow.OwnerName,
			AcceptedSolutionID: problemRow.AcceptedSolutionID.Int64,
			CreatedAt:          problemRow.CreatedAt,
		})
	}

	return problems, nil
}

func storeProblem(ctx context.Context, problem ProblemToStore) (CreateProblemResponse, *globals.HTTPError) {
	p, err := storeProblemInDB(ctx, problem)
	if err != nil {
		return CreateProblemResponse{}, &globals.HTTPError{
			Status: http.StatusInternalServerError,
			Title:  "Database Error.",
			Err:    errors.New("Could not store problem in database."),
		}
	}

	prob := CreateProblemResponse{
		ID:                 p.ID,
		Title:              p.Title,
		Description:        p.Description,
		OwnerID:            p.OwnerID,
		OwnerName:          p.OwnerName,
		AcceptedSolutionID: p.AcceptedSolutionID.Int64,
		CreatedAt:          p.CreatedAt,
	}

	return prob, nil
}

func getFullNameByID(ctx context.Context, id int64) (string, *globals.HTTPError) {
	name, err := getFullNameByIDFromDB(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", &globals.HTTPError{
				Status: http.StatusNoContent,
				Title:  "No Content.",
				Err:    errors.New("No problem to show."),
			}
		}

		return "", &globals.HTTPError{
			Status: http.StatusInternalServerError,
			Title:  "Database Error.",
			Err:    errors.New("Could not get fullname from database."),
		}
	}

	return name, nil
}

func canDeleteProblem(ctx context.Context, problem GetProblemResponse) *globals.HTTPError {
	solutionCount, err := getSolutionCountForProblemFromDB(ctx, problem.ID)
	if err != nil {
		return &globals.HTTPError{
			Status: http.StatusInternalServerError,
			Title:  "Database Error.",
			Err:    errors.New("Could not get solution count from database."),
		}
	}

	if solutionCount != 0 {
		return &globals.HTTPError{
			Status: http.StatusBadRequest,
			Title:  "Deletion Error.",
			Err:    errors.New("Problem does not have 0 solutions."),
		}
	}

	return nil
}

func deleteProblemByID(ctx context.Context, id int64) *globals.HTTPError {
	err := deleteProblemByIDFromDB(ctx, id)
	if err != nil {
		return &globals.HTTPError{
			Status: http.StatusInternalServerError,
			Title:  "Database Error.",
			Err:    errors.New("Could not get problem from database."),
		}
	}

	return nil
}

func setAcceptedSolutionID(ctx context.Context, pid int64, sid int64) *globals.HTTPError {
	err := setAcceptedSolutionIDInDB(ctx, pid, sid)
	if err != nil {
		return &globals.HTTPError{
			Status: http.StatusInternalServerError,
			Title:  "Database Error.",
			Err:    errors.New("Could not set accepted solution id in database."),
		}
	}

	return nil
}

func unsetAcceptedSolutionID(ctx context.Context, pid int64) *globals.HTTPError {
	err := unsetAcceptedSolutionIDInDB(ctx, pid)
	if err != nil {
		return &globals.HTTPError{
			Status: http.StatusInternalServerError,
			Title:  "Database Error.",
			Err:    errors.New("Could not set accepted solution id in database."),
		}
	}

	return nil
}
