package problems

import (
	"encoding/json"
	"errors"
	"io"
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
	var problem CreateProblemRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := globals.ReturnErrorResponse(
			"Server Error.",
			"Could not read the body",
		)
		w.Write(response)

		return
	}

	err = json.Unmarshal(body, &problem)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := globals.ReturnErrorResponse(
			"Bad Request.",
			"Could not parse body.",
		)
		w.Write(response)

		return
	}

	err = problem.validate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := globals.ReturnErrorResponse(
			"Bad Request.",
			err.Error(),
		)
		w.Write(response)

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
		w.WriteHeader(http.StatusInternalServerError)
		response := globals.ReturnErrorResponse(
			"Server Error.",
			"Could not return created problem.",
		)
		w.Write(response)

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
