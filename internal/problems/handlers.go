package problems

import (
	"net/http"

	"github.com/mithcs/probox-api/internal/globals"
)

func GetProblem(w http.ResponseWriter, r *http.Request) {
	id, err := globals.GetIDFromRequest(r)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	// TODO: problem might not exist
	problem, err := getProblemByID(r.Context(), id)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	response, err := globals.EncodeJson(problem)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	w.Write(response)
}

func GetProblems(w http.ResponseWriter, r *http.Request) {
	problems, err := getAllProblems(r.Context())
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	response, err := globals.EncodeJson(problems)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	w.Write(response)
}

func CreateProblem(w http.ResponseWriter, r *http.Request) {
	problem, err := globals.ParseBody[CreateProblemRequest](r.Body)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	err = problem.Validate()
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

	p, err := storeProblem(r.Context(), ProblemToStore{
		Title:       problem.Title,
		Description: problem.Description,
		OwnerID:     uid,
		OwnerName:   name,
	})
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	problemRes := CreateProblemResponse{
		ID:                 p.ID,
		Title:              p.Title,
		Description:        p.Description,
		OwnerID:            p.OwnerID,
		OwnerName:          p.OwnerName,
		AcceptedSolutionID: p.AcceptedSolutionID.Int64,
	}

	response, err := globals.EncodeJson(problemRes)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	w.Write(response)
}

func DeleteProblem(w http.ResponseWriter, r *http.Request) {
	id, err := globals.GetIDFromRequest(r)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	problem, err := getProblemByID(r.Context(), id)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	err = canDeleteProblem(r.Context(), problem)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	err = globals.VerifyUIDWithToken(r.Context(), problem.OwnerID)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	err = deleteProblemByID(r.Context(), id)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func AcceptSolution(w http.ResponseWriter, r *http.Request) {
	pid, err := globals.GetIDFromRequest(r)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	problem, err := getProblemByID(r.Context(), pid)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	err = globals.VerifyUIDWithToken(r.Context(), problem.OwnerID)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	solution, err := globals.ParseBody[AcceptSolutionRequest](r.Body)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	err = setAcceptedSolutionID(r.Context(), pid, solution.SolutionID)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UnacceptSolution(w http.ResponseWriter, r *http.Request) {
	pid, err := globals.GetIDFromRequest(r)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	problem, err := getProblemByID(r.Context(), pid)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	err = globals.VerifyUIDWithToken(r.Context(), problem.OwnerID)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	err = unsetAcceptedSolutionID(r.Context(), pid)
	if err != nil {
		globals.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
