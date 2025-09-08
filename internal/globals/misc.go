package globals

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func ParseBody[T any](body io.ReadCloser) (T, *HTTPError) {
	var v T

	readBody, err := io.ReadAll(body)
	if err != nil {
		return v, &HTTPError{
			Status: http.StatusInternalServerError,
			Title:  "Read Error.",
			Err:    errors.New("Could not read body."),
		}
	}

	err = json.Unmarshal(readBody, &v)
	if err != nil {
		return v, &HTTPError{
			Status: http.StatusInternalServerError,
			Title:  "Decode Error.",
			Err:    errors.New("Could not unmarshal json."),
		}
	}

	return v, nil
}

func EncodeJson(v any) ([]byte, *HTTPError) {
	json, err := json.Marshal(v)
	if err != nil {
		return []byte(""), &HTTPError{
			Status: http.StatusInternalServerError,
			Title:  "Encode Error.",
			Err:    errors.New("Could not marshal json."),
		}
	}

	return json, nil
}
