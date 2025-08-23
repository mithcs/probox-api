package globals

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Status  int
	Title   string
	Details string
}

type errorResponse struct {
	Title   string `json:"title"`
	Details string `json:"details"`
}

func WriteErrorResponse(w http.ResponseWriter, error Error) {
	w.WriteHeader(error.Status)

	e := errorResponse{
		Title:   error.Title,
		Details: error.Details,
	}

	response, _ := json.Marshal(e)
	w.Write(response)
}
