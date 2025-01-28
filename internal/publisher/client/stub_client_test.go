package client_test

import (
	"net/http"
	"testing"

	"github.com/anselstetter/play-publisher/internal/assert"
	"github.com/anselstetter/play-publisher/internal/publisher/client"
)

func TestStubClient(t *testing.T) {
	tests := []struct {
		name          string
		url           string
		statusCode    int
		fields        *map[string]any
		headers       *http.Header
		idStatusCodes []client.IdStatusCode
	}{
		{
			name:       "should return 200 on edits endpoint",
			url:        "https://androidpublisher.googleapis.com/androidpublisher/v3/applications/application-id/edits",
			statusCode: http.StatusOK,
			fields: &map[string]any{
				"Id":                "session",
				"ExpiryTimeSeconds": "10000",
			},
		},
		{
			name:       "should return 200 on bundles endpoint",
			url:        "https://androidpublisher.googleapis.com/upload/androidpublisher/v3/applications/application-id/edits/session/bundles",
			statusCode: http.StatusOK,
			headers: &http.Header{
				"Location": {"http://localhost"},
			},
		},
		{
			name:       "should return 200 on apks endpoint",
			url:        "https://androidpublisher.googleapis.com/upload/androidpublisher/v3/applications/application-id/edits/session/apks",
			statusCode: http.StatusOK,
			headers: &http.Header{
				"Location": {"http://localhost"},
			},
		},
		{
			name:       "should return 200 on localhost endpoint",
			url:        "http://localhost",
			statusCode: http.StatusOK,
			fields: &map[string]any{
				"VersionCode": float64(1), // When you are JSON, every number is float(64)
				"Sha1":        "Sha1",
				"Sha256":      "Sha256",
			},
		},
		{
			name:       "should return 200 on tracks endpoint",
			url:        "https://androidpublisher.googleapis.com/androidpublisher/v3/applications/application-id/edits/session/tracks/internal",
			statusCode: http.StatusOK,
		},
		{
			name:       "should return 200 on edits commit endpoint",
			url:        "https://androidpublisher.googleapis.com/androidpublisher/v3/applications/application-id/edits/session:commit",
			statusCode: http.StatusOK,
			fields: &map[string]any{
				"Id":                "session",
				"ExpiryTimeSeconds": "10000",
			},
		},
		{
			name:       "should return 500 on edits commit endpoint",
			url:        "https://androidpublisher.googleapis.com/androidpublisher/v3/applications/application-id/edits/session:commit",
			statusCode: http.StatusInternalServerError,
			fields: &map[string]any{
				"Id":                "session",
				"ExpiryTimeSeconds": "10000",
			},
			idStatusCodes: []client.IdStatusCode{client.NewIdStatusCode(client.IdCommit, http.StatusInternalServerError)},
		},
		{
			name:       "should return 404 on missing stub",
			url:        "https://not-included",
			statusCode: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			stubClient, _ := client.NewStubClient(tc.idStatusCodes...)

			res, err := stubClient.Get(tc.url)
			assert.NoError(t, err)

			if tc.fields != nil {
				decoded, err := client.Decode[map[string]any](res.Body)

				assert.NoError(t, err)
				assert.Equals(t, decoded, *tc.fields)
			}
			if tc.headers != nil {
				assert.Equals(t, res.Header, *tc.headers)
			}

			assert.True(t, res.StatusCode == tc.statusCode)
		})
	}
}
