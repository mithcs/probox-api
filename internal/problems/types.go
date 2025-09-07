package problems

import (
	"errors"
	"net/http"

	"github.com/mithcs/probox-api/internal/globals"
)

type GetProblemResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	OwnerID     int64  `json:"owner_id"`
	OwnerName   string `json:"owner_name"`
}

type CreateProblemRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateProblemResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	OwnerID     int64  `json:"owner_id"`
	OwnerName   string `json:"owner_name"`
}

type ProblemToStore struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	OwnerID     int64  `json:"owner_id"`
	OwnerName   string `json:"owner_name"`
}

func (problem *CreateProblemRequest) Validate() *globals.HTTPError {
	if err := validateTitle(problem.Title); err != nil {
		return &globals.HTTPError{
			Status:      http.StatusBadRequest,
			Title:       "Invalid Title.",
			Description: "Title does not meet requirements.",
		}
	}

	if err := validateDescription(problem.Description); err != nil {
		return &globals.HTTPError{
			Status:      http.StatusBadRequest,
			Title:       "Invalid Description.",
			Description: "Description does not meet requirements.",
		}
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
