package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"slices"
)

const (
	IdEdit      = "edit"
	IdUploadAab = "uploadAab"
	IdUploadApk = "uploadApk"
	IdApp       = "app"
	IdTrack     = "track"
	IdCommit    = "commit"
)

// -----------------------------------------------------------------------------
// Test client
// -----------------------------------------------------------------------------
type roundTripFunc func(req *http.Request) *http.Response

func (rtf roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return rtf(req), nil
}

func newTestClient(fn roundTripFunc) *http.Client {
	return &http.Client{
		Transport: roundTripFunc(fn),
	}
}

// -----------------------------------------------------------------------------
// Response
// -----------------------------------------------------------------------------
type response struct {
	id         string
	host       string
	path       string
	statusCode int
	fn         func(body io.ReadCloser, statusCode int) *http.Response
}

func newEditResponse(body io.ReadCloser, statusCode int) *http.Response {
	response := struct {
		Id                string
		ExpiryTimeSeconds string
	}{
		Id:                "session",
		ExpiryTimeSeconds: "10000",
	}
	b, _ := json.Marshal(response)

	res := &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewBuffer(b)),
		Header:     make(http.Header),
	}
	return res
}

func newUploadResponse(body io.ReadCloser, statusCode int) *http.Response {
	header := make(http.Header)
	header.Add("Location", "http://localhost")

	res := &http.Response{
		StatusCode: statusCode,
		Header:     header,
	}
	return res
}

func newAppResponse(body io.ReadCloser, statusCode int) *http.Response {
	response := struct {
		VersionCode int64
		Sha1        string
		Sha256      string
	}{
		VersionCode: 1,
		Sha1:        "Sha1",
		Sha256:      "Sha256",
	}
	b, _ := json.Marshal(response)

	res := &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewBuffer(b)),
		Header:     make(http.Header),
	}
	return res
}

func newTrackResponse(body io.ReadCloser, statusCode int) *http.Response {
	res := &http.Response{
		StatusCode: statusCode,
		Body:       body,
		Header:     make(http.Header),
	}
	return res
}

func newCommitResponse(body io.ReadCloser, statusCode int) *http.Response {
	return newEditResponse(body, statusCode)
}

// -----------------------------------------------------------------------------
// Response mapping
// -----------------------------------------------------------------------------
type responseMapping struct {
	calls     []string
	responses []response
}

func (rm *responseMapping) respond(req *http.Request) *http.Response {
	idx := slices.IndexFunc(rm.responses, func(r response) bool {
		regex := regexp.MustCompile(r.path)
		return r.host == req.URL.Host && regex.MatchString(req.URL.Path)
	})
	if idx == -1 {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			//	Body:       body,
			Header: make(http.Header),
		}
	}
	response := rm.responses[idx]

	rm.calls = append(rm.calls, response.id)

	return response.fn(req.Body, response.statusCode)
}

func (rm *responseMapping) WithStatusCode(id string, statusCode int) responseMapping {
	idx := slices.IndexFunc(rm.responses, func(r response) bool {
		return r.id == id
	})
	rm.responses[idx].statusCode = statusCode
	return *rm
}

func newResponseMapping() responseMapping {
	return responseMapping{
		calls: []string{},
		responses: []response{
			{
				id:         IdEdit,
				host:       "androidpublisher.googleapis.com",
				path:       "^/androidpublisher/v3/applications/.+/edits$",
				statusCode: http.StatusOK,
				fn:         newEditResponse,
			},
			{
				id:         IdUploadAab,
				host:       "androidpublisher.googleapis.com",
				path:       "^/upload/androidpublisher/v3/applications/.+/edits/.+/bundles$",
				statusCode: http.StatusOK,
				fn:         newUploadResponse,
			},
			{
				id:         IdUploadApk,
				host:       "androidpublisher.googleapis.com",
				path:       "^/upload/androidpublisher/v3/applications/.+/edits/.+/apks$",
				statusCode: http.StatusOK,
				fn:         newUploadResponse,
			},
			{
				id:         IdApp,
				host:       "localhost",
				path:       "",
				statusCode: http.StatusOK,
				fn:         newAppResponse,
			},
			{
				id:         IdTrack,
				host:       "androidpublisher.googleapis.com",
				path:       "^/androidpublisher/v3/applications/.+/edits/.+/tracks/.+$",
				statusCode: http.StatusOK,
				fn:         newTrackResponse,
			},
			{
				id:         IdCommit,
				host:       "androidpublisher.googleapis.com",
				path:       "^/androidpublisher/v3/applications/.+/edits/.+:commit$",
				statusCode: http.StatusOK,
				fn:         newCommitResponse,
			},
		},
	}
}

// -----------------------------------------------------------------------------
// Id and status code
// -----------------------------------------------------------------------------
type IdStatusCode struct {
	id           string
	idStatusCode int
}

func NewIdStatusCode(id string, statusCode int) IdStatusCode {
	return IdStatusCode{
		id:           id,
		idStatusCode: statusCode,
	}
}

// -----------------------------------------------------------------------------
// Publisher client
// -----------------------------------------------------------------------------
func NewStubClient(idStatusCodes ...IdStatusCode) (*http.Client, *[]string) {
	mapping := newResponseMapping()

	for _, idSidStatusCode := range idStatusCodes {
		mapping = mapping.WithStatusCode(idSidStatusCode.id, idSidStatusCode.idStatusCode)
	}
	c := newTestClient(func(req *http.Request) *http.Response {
		return mapping.respond(req)
	})
	return c, &mapping.calls
}
