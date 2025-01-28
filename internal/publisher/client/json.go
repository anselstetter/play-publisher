package client

import (
	"encoding/json"
	"fmt"
	"io"
)

func Decode[T any](r io.ReadCloser) (T, error) {
	var v T
	if err := json.NewDecoder(r).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}
