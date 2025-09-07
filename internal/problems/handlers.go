package problems

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mithcs/probox-api/internal/globals"
)

func GetProblem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusBadRequest,
			Title:   "Bad Request.",
			Details: "Invalid problem id.",
		})

		return
	}

	problem, err := getProblemByID(r.Context(), id)
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

	err = problem.Validate()
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusBadRequest,
			Title:   "Bad Request.",
			Details: err.Error(),
		})

		return
	}

	uid, err := globals.GetUserIDFromContext(r.Context())
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
