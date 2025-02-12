package client_test

import (
	"io"
	"strings"
	"testing"

	"github.com/anselstetter/play-publisher/internal/assert"
	"github.com/anselstetter/play-publisher/internal/publisher/client"
)

func TestJsonDecode(t *testing.T) {
	t.Parallel()

	type test struct{ Key string }

	readCloser := io.NopCloser(strings.NewReader(`{"key":"value"}`))
	got, _ := client.Decode[test](readCloser)
	want := test{Key: "value"}

	assert.Equals(t, got, want)
}

func TestJsonDecodeError(t *testing.T) {
	t.Parallel()

	type test struct{ Key string }

	readCloser := io.NopCloser(strings.NewReader(`malformed`))
	_, err := client.Decode[test](readCloser)

	assert.AnyError(t, err)
}
