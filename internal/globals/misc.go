package globals

import (
	"encoding/json"
	"io"
	"net/http"
)

func ParseBody[T any](body io.ReadCloser) (T, *HTTPError) {
	var v T

	readBody, err := io.ReadAll(body)
	if err != nil {
		return v, &HTTPError{
			Status:      http.StatusInternalServerError,
			Title:       "Read Error.",
			Description: "Could not read body.",
		}
	}

	err = json.Unmarshal(readBody, &v)
	if err != nil {
		return v, &HTTPError{
			Status:      http.StatusInternalServerError,
			Title:       "Decode Error.",
			Description: "Could not unmarshal json.",
		}
	}

	return v, nil
}

func GetJson(v any) ([]byte, *HTTPError) {
	json, err := json.Marshal(v)
	if err != nil {
		return []byte(""), &HTTPError{
			Status:      http.StatusInternalServerError,
			Title:       "Encode Error.",
			Description: "Could not marshal json.",
		}
	}

	return json, nil
}
