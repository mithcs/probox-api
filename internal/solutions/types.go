package solutions

import (
	"context"
	"errors"

	"github.com/mithcs/probox-api/internal/globals"
)

type GetSolutionResponse struct {
	ID          int64  `json:"id"`
	ProblemID   int64  `json:"problem_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type GetSolutionsForProblemResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateSolutionRequest struct {
	ProblemID   int64  `json:"problem_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateSolutionResponse struct {
	ID          int64  `json:"id"`
	ProblemID   int64  `json:"problem_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (solution *CreateSolutionRequest) Validate(ctx context.Context) error {
	if err := validateProblemID(ctx, solution.ProblemID); err != nil {
		return err
	}

	if err := validateTitle(solution.Title); err != nil {
		return err
	}

	if err := validateDescription(solution.Description); err != nil {
		return err
	}

	return nil
}

func validateProblemID(ctx context.Context, problemID int64) error {
	problem, _ := globals.Queries.GetProblemByID(ctx, problemID)
	if problem.ID == 0 || problem.Title == "" || problem.Description == "" {
		return errors.New("Invalid Problem ID.")
	}

	return nil
}

func validateTitle(title string) error {
	minLen := 3
	maxLen := 120

	if len(title) > minLen && len(title) < maxLen {
		return nil
	}

	return errors.New("Invalid Title.")
}

func validateDescription(description string) error {
	minLen := 5
	maxLen := 20000

	if len(description) > minLen && len(description) < maxLen {
		return nil
	}

	return errors.New("Invalid Description.")
}
