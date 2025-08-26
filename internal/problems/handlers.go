package problems

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mithcs/probox-api/internal/db"
	"github.com/mithcs/probox-api/internal/globals"
)

type GetProblemsResponse struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateProblemRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateProblemResponse struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func GetProblems(w http.ResponseWriter, r *http.Request) {
	problems, err := getAllProblems(r.Context())
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not get all problems.",
		})

		return
	}

	response, err := json.Marshal(problems)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not list all problems.",
		})

		return
	}

	w.Write(response)
}

func getAllProblems(ctx context.Context) ([]GetProblemsResponse, error) {
	problemRows, err := getProblemsFromDb(ctx)
	if err != nil {
		return []GetProblemsResponse{}, err
	}

	var problems []GetProblemsResponse
	for _, problemRow := range problemRows {
		problems = append(problems, GetProblemsResponse{
			Id:          problemRow.ID,
			Title:       problemRow.Title,
			Description: problemRow.Description,
		})
	}

	return problems, nil
}

func getProblemsFromDb(ctx context.Context) ([]db.Problem, error) {
	problemRows, err := globals.Queries.GetProblems(ctx)

	return problemRows, err
}

func CreateProblem(w http.ResponseWriter, r *http.Request) {
	problem, err := globals.ParseBody[CreateProblemRequest](r.Body)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusBadRequest,
			Title:   "Bad Request",
			Details: "Could not parse the body.",
		})

		return
	}

	err = problem.validate()
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusBadRequest,
			Title:   "Bad Request.",
			Details: err.Error(),
		})

		return
	}

	p, err := storeProblem(r.Context(), problem)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not store problem.",
		})

		return
	}

	problemRes := CreateProblemResponse{
		Id:          p.ID,
		Title:       p.Title,
		Description: p.Description,
	}

	response, err := json.Marshal(problemRes)
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

func storeProblem(ctx context.Context, problem CreateProblemRequest) (db.Problem, error) {
	p, err := globals.Queries.CreateProblem(ctx, db.CreateProblemParams{
		Title:       problem.Title,
		Description: problem.Description,
	})

	return p, err
}

func (problem *CreateProblemRequest) validate() error {
	if err := validateTitle(problem.Title); err != nil {
		return err
	}

	if err := validateDescription(problem.Description); err != nil {
		return err
	}

	return nil
}

func validateTitle(title string) error {
	minLen := 5
	maxLen := 120

	if len(title) > minLen && len(title) < maxLen {
		return nil
	}

	return errors.New("Invalid Title.")
}

func validateDescription(description string) error {
	minLen := 8
	maxLen := 20000

	if len(description) > minLen && len(description) < maxLen {
		return nil
	}

	return errors.New("Invalid Description.")
}
