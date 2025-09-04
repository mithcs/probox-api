package problems

import "errors"

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

func (problem *CreateProblemRequest) Validate() error {
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
