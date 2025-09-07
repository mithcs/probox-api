package solutions

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mithcs/probox-api/internal/globals"
)

func GetSolution(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusBadRequest,
			Title:   "Bad Request.",
			Details: "Invalid Solution ID.",
		})

		return
	}

	// TODO: solution might not exist
	solution, err := getSolutionByID(r.Context(), id)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not get solution.",
		})

		return
	}

	response, err := json.Marshal(solution)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not list solution.",
		})

		return
	}

	w.Write(response)
}

func GetSolutions(w http.ResponseWriter, r *http.Request) {
	if problemID := r.URL.Query().Get("problem_id"); problemID != "" {
		GetSolutionsForProblem(w, r, problemID)
		return
	}

	solutions, err := getAllSolutions(r.Context())
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not get all solutions.",
		})

		return
	}

	response, err := json.Marshal(solutions)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not list all solutions.",
		})

		return
	}

	w.Write(response)
}

func GetSolutionsForProblem(w http.ResponseWriter, r *http.Request, problemID string) {
	pid, err := strconv.ParseInt(problemID, 10, 64)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusBadRequest,
			Title:   "Bad Request.",
			Details: "Invalid Problem ID.",
		})

		return
	}

	solutions, err := getSolutionsByProblemID(r.Context(), pid)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusBadRequest,
			Title:   "Bad Request.",
			Details: err.Error(),
		})

		return
	}

	if len(solutions) == 0 {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusNoContent,
			Title:   "No Content.",
			Details: "There are no solutions to list.",
		})

		return
	}

	response, err := json.Marshal(solutions)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not list all solutions.",
		})

		return
	}

	w.Write(response)
}

func CreateSolution(w http.ResponseWriter, r *http.Request) {
	solution, err := globals.ParseBody[CreateSolutionRequest](r.Body)
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusBadRequest,
			Title:   "Bad Request.",
			Details: "Could not parse body.",
		})

		return
	}

	err = solution.Validate(r.Context())
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

	s, err := storeSolution(r.Context(), SolutionToStore{
		ProblemID:   solution.ProblemID,
		Title:       solution.Title,
		Description: solution.Description,
		OwnerID:     uid,
		OwnerName:   name,
	})
	if err != nil {
		globals.WriteErrorResponse(w, globals.Error{
			Status:  http.StatusInternalServerError,
			Title:   "Server Error.",
			Details: "Could not store solution.",
		})

		return
	}

	solutionRes := CreateSolutionResponse{
		ID:          s.ID,
		ProblemID:   s.ProblemID,
		Title:       s.Title,
		Description: s.Description,
		OwnerID:     s.OwnerID,
		OwnerName:   s.OwnerName,
	}

	response, err := json.Marshal(solutionRes)
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
