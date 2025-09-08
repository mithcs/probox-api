package globals

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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

func GetIDFromRequest(r *http.Request) (int64, *HTTPError) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		return -1, &HTTPError{
			Status: http.StatusBadRequest,
			Title:  "Invalid ID.",
			Err:    errors.New("Could not parse id."),
		}
	}

	return id, nil
}

func VerifyUIDWithToken(ctx context.Context, id int64) *HTTPError {
	uid, err := GetUserIDFromContext(ctx)
	if err != nil {
		return &HTTPError{
			Status: http.StatusBadRequest,
			Title:  "Invalid User ID.",
			Err:    errors.New("Could not get user id from request."),
		}
	}

	if uid != id {
		return &HTTPError{
			Status: http.StatusBadRequest,
			Title:  "Unauthorized Action.",
			Err:    errors.New("Not authorized to perform this action to this user."),
		}
	}

	return nil
}
