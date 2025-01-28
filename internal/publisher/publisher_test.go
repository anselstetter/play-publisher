package publisher_test

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/anselstetter/play-publisher/internal/application"
	"github.com/anselstetter/play-publisher/internal/assert"
	"github.com/anselstetter/play-publisher/internal/logger"
	"github.com/anselstetter/play-publisher/internal/publisher"
	"github.com/anselstetter/play-publisher/internal/publisher/client"
)

func TestPublisher(t *testing.T) {
	upload := func(app string, track string, serviceAccount string, idStatusCodes ...client.IdStatusCode) (*[]string, error) {
		c, calls := client.NewStubClient(idStatusCodes...)
		buffer := bytes.NewBuffer([]byte{})
		analyzer := application.New()
		logger := logger.New(logger.WithStdout(buffer), logger.WithStderr(buffer))

		err := publisher.New(analyzer, logger, publisher.WithHttpClient(c)).
			Upload(context.Background(), app, track, serviceAccount)

		return calls, err
	}

	t.Run("should succeed uploading AAB without error", func(t *testing.T) {
		t.Parallel()

		var (
			app            = "../testdata/app-release.aab"
			track          = "internal"
			serviceAccount = "../testdata/google-service-account.json"
		)
		calls, err := upload(app, track, serviceAccount)

		assert.NoError(t, err)
		assert.Equals(t, *calls, []string{client.IdEdit, client.IdUploadAab, client.IdApp, client.IdTrack, client.IdCommit})
	})

	t.Run("should succeed uploading APK without error", func(t *testing.T) {
		t.Parallel()

		var (
			app            = "../testdata/app-release.apk"
			track          = "internal"
			serviceAccount = "../testdata/google-service-account.json"
		)
		calls, err := upload(app, track, serviceAccount)

		assert.NoError(t, err)
		assert.Equals(t, *calls, []string{client.IdEdit, client.IdUploadApk, client.IdApp, client.IdTrack, client.IdCommit})
	})

	t.Run("should return publisher.ErrorAnalyzeApp", func(t *testing.T) {
		t.Parallel()

		var (
			app            = "invalid"
			track          = "internal"
			serviceAccount = "missing"
		)
		calls, err := upload(app, track, serviceAccount)

		assert.IsError(t, err, publisher.ErrorAnalyzeApp)
		assert.Equals(t, *calls, []string{})
	})

	t.Run("should return publisher.ErrorReadServiceAccount", func(t *testing.T) {
		t.Parallel()

		var (
			app            = "../testdata/app-release.aab"
			track          = "internal"
			serviceAccount = "missing"
		)
		calls, err := upload(app, track, serviceAccount)

		assert.IsError(t, err, publisher.ErrorReadServiceAccount)
		assert.Equals(t, *calls, []string{})
	})

	t.Run("should return publisher.ErrorCreateEdit", func(t *testing.T) {
		t.Parallel()

		var (
			app            = "../testdata/app-release.aab"
			track          = "internal"
			serviceAccount = "../testdata/google-service-account.json"
		)
		calls, err := upload(app, track, serviceAccount, client.NewIdStatusCode(client.IdEdit, http.StatusInternalServerError))

		assert.IsError(t, err, publisher.ErrorCreateEdit)
		assert.Equals(t, *calls, []string{client.IdEdit})
	})

	t.Run("should return publisher.ErrorUploadAab", func(t *testing.T) {
		t.Parallel()

		var (
			app            = "../testdata/app-release.aab"
			track          = "internal"
			serviceAccount = "../testdata/google-service-account.json"
		)
		calls, err := upload(app, track, serviceAccount, client.NewIdStatusCode(client.IdUploadAab, http.StatusInternalServerError))

		assert.IsError(t, err, publisher.ErrorUploadAab)
		assert.Equals(t, *calls, []string{client.IdEdit, client.IdUploadAab})
	})

	t.Run("should return publisher.ErrorUploadApk", func(t *testing.T) {
		t.Parallel()

		var (
			app            = "../testdata/app-release.apk"
			track          = "internal"
			serviceAccount = "../testdata/google-service-account.json"
		)
		calls, err := upload(app, track, serviceAccount, client.NewIdStatusCode(client.IdUploadApk, http.StatusInternalServerError))

		assert.IsError(t, err, publisher.ErrorUploadApk)
		assert.Equals(t, *calls, []string{client.IdEdit, client.IdUploadApk})
	})

	t.Run("should return publisher.ErrorUpdateTrack", func(t *testing.T) {
		t.Parallel()

		var (
			app            = "../testdata/app-release.aab"
			track          = "internal"
			serviceAccount = "../testdata/google-service-account.json"
		)
		calls, err := upload(app, track, serviceAccount, client.NewIdStatusCode(client.IdTrack, http.StatusInternalServerError))

		assert.IsError(t, err, publisher.ErrorUpdateTrack)
		assert.Equals(t, *calls, []string{client.IdEdit, client.IdUploadAab, client.IdApp, client.IdTrack})
	})

	t.Run("should return publisher.ErrorCommitEdit", func(t *testing.T) {
		t.Parallel()

		var (
			track          = "internal"
			serviceAccount = "../testdata/google-service-account.json"
			app            = "../testdata/app-release.aab"
		)
		calls, err := upload(app, track, serviceAccount, client.NewIdStatusCode(client.IdCommit, http.StatusInternalServerError))

		assert.IsError(t, err, publisher.ErrorCommitEdit)
		assert.Equals(t, *calls, []string{client.IdEdit, client.IdUploadAab, client.IdApp, client.IdTrack, client.IdCommit})
	})
}
