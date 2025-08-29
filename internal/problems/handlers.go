package problems

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/mithcs/probox-api/internal/db"
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

func GetProblem(w http.ResponseWriter, r *http.Request) {
	problemID, err := strconv.ParseInt(chi.URLParam(r, "problemID"), 10, 64)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusBadRequest,
			Title:   "Bad Request.",
			Details: "Invalid problem id.",
		})

		return
	}

	problem, err := getProblemByID(r.Context(), problemID)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not get problem.",
		})

		return
	}

	response, err := json.Marshal(problem)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not list problem.",
		})

		return
	}

	w.Write(response)
}

func getProblemByID(ctx context.Context, problemID int64) (GetProblemResponse, error) {
	p, err := getProblemByIDFromDb(ctx, problemID)

	problem := GetProblemResponse{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		OwnerID:     p.OwnerID,
		OwnerName:   p.OwnerName,
	}

	return problem, err
}

func getProblemByIDFromDb(ctx context.Context, problemID int64) (db.Problem, error) {
	problem, err := globals.Queries.GetProblemByID(ctx, problemID)

	return problem, err
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

func getAllProblems(ctx context.Context) ([]GetProblemResponse, error) {
	problemRows, err := getProblemsFromDb(ctx)
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

	uid, err := getUserIDFromContext(r.Context())
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not get User ID from token.",
		})

		return
	}

	name, err := getFullNameByID(r.Context(), uid)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not get Full Name from User ID.",
		})

		return
	}

	p, err := storeProblem(r.Context(), ProblemToStore{
		Title:       problem.Title,
		Description: problem.Description,
		OwnerID:     uid,
		OwnerName:   name,
	})
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not store problem.",
		})

		return
	}

	problemRes := CreateProblemResponse{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		OwnerID:     p.OwnerID,
		OwnerName:   p.OwnerName,
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

func storeProblem(ctx context.Context, problem ProblemToStore) (db.Problem, error) {
	p, err := globals.Queries.CreateProblem(ctx, db.CreateProblemParams{
		Title:       problem.Title,
		Description: problem.Description,
		OwnerID:     problem.OwnerID,
		OwnerName:   problem.OwnerName,
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

func getUserIDFromContext(ctx context.Context) (int64, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return -1, err
	}

	uid, err := strconv.ParseInt(fmt.Sprintf("%v", claims["uid"]), 10, 64)
	if err != nil {
		return -1, err
	}

	return uid, nil
}

func getFullNameByID(ctx context.Context, id int64) (string, error) {
	return getFullNameByIDFromDb(ctx, id)
}

func getFullNameByIDFromDb(ctx context.Context, id int64) (string, error) {
	return globals.Queries.GetFullNameByID(ctx, id)
}
