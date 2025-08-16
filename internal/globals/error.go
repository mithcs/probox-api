package globals

import "encoding/json"

type ErrorResponse struct {
	Title   string `json:"title"`
	Details string `json:"details"`
}

func ReturnErrorResponse(title string, details string) []byte {
	error := ErrorResponse{
		Title:   title,
		Details: details,
	}

	res, _ := json.Marshal(error)

	return res
}
