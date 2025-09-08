package solutions

import (
	"net/http"

	"github.com/mithcs/probox-api/internal/globals"
)

func GetSolution(w http.ResponseWriter, r *http.Request) {
	id, err := globals.GetIDFromRequest(r)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	// TODO: solution might not exist
	solution, err := getSolutionByID(r.Context(), id)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	response, err := globals.EncodeJson(solution)
	if err != nil {
		globals.WriteErrorResponse(w, err)
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
		globals.WriteErrorResponse(w, err)

		return
	}

	response, err := globals.EncodeJson(solutions)
	if err != nil {
		globals.WriteErrorResponse(w, err)

		return
	}

	w.Write(response)
}

func GetSolutionsForProblem(w http.ResponseWriter, r *http.Request, problemID string) {
	pid, err := globals.GetIDFromRequest(r)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	solutions, err := getSolutionsByProblemID(r.Context(), pid)
	if err != nil {
		globals.WriteErrorResponse(w, err)

		return
	}

	response, err := globals.EncodeJson(solutions)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	w.Write(response)
}

func CreateSolution(w http.ResponseWriter, r *http.Request) {
	solution, err := globals.ParseBody[CreateSolutionRequest](r.Body)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	err = solution.Validate(r.Context())
	if err != nil {
		globals.WriteErrorResponse(w, err)

		return
	}

	uid, err := globals.GetUserIDFromContext(r.Context())
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	name, err := getFullNameByID(r.Context(), uid)
	if err != nil {
		globals.WriteErrorResponse(w, err)
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
		globals.WriteErrorResponse(w, err)
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

	response, err := globals.EncodeJson(solutionRes)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	w.Write(response)
}
