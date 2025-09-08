package globals

import (
	"encoding/json"
	"net/http"
)

type HTTPError struct {
	Status int
	Title  string
	Err    error
}

func (e *HTTPError) Error() string {
	return e.Err.Error()
}

type _ErrorResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func WriteErrorResponse(w http.ResponseWriter, err *HTTPError) {
	w.WriteHeader(err.Status)

	e := _ErrorResponse{
		Title:       err.Title,
		Description: err.Error(),
	}

	response, _ := json.Marshal(e)
	w.Write(response)
}
