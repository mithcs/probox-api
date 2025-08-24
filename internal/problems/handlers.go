package problems

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mithcs/probox-api/internal/globals"
)

type CreateProblemRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
type CreateProblemResponse struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
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

	// store problem in db
	// and get problemId
	problemId := 1

	problemRes := CreateProblemResponse{
		Id:          problemId,
		Title:       problem.Title,
		Description: problem.Description,
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

func (problem *CreateProblemRequest) validate() error {
	// TODO

	if problem.Title == "invalid title" {
		return errors.New("Invalid Title.")
	}
	if problem.Description == "invalid description" {
		return errors.New("Invalid Description.")
	}

	return nil
}
