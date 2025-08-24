package globals

import (
	"encoding/json"
	"io"
)

func ParseBody[T any](body io.ReadCloser) (T, error) {
	var v T

	readBody, err := io.ReadAll(body)
	if err != nil {
		return v, err
	}

	err = json.Unmarshal(readBody, &v)
	return v, err
}
