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
			Upload(context.Background(), app, track, publisher.StatusCompleted, serviceAccount)

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

	t.Run("should return publisher.ErrUpload", func(t *testing.T) {
		t.Parallel()

		var (
			app            = "invalid"
			track          = "internal"
			serviceAccount = "missing"
		)
		calls, err := upload(app, track, serviceAccount)

		assert.IsError(t, err, publisher.ErrUpload)
		assert.Equals(t, *calls, []string{})
	})

	t.Run("should return publisher.ErrorUpload", func(t *testing.T) {
		t.Parallel()

		var (
			app            = "../testdata/app-release.aab"
			track          = "internal"
			serviceAccount = "missing"
		)
		calls, err := upload(app, track, serviceAccount)

		assert.IsError(t, err, publisher.ErrUpload)
		assert.Equals(t, *calls, []string{})
	})

	t.Run("should return publisher.ErrStartEdit", func(t *testing.T) {
		t.Parallel()

		var (
			app            = "../testdata/app-release.aab"
			track          = "internal"
			serviceAccount = "../testdata/google-service-account.json"
		)
		calls, err := upload(app, track, serviceAccount, client.NewIdStatusCode(client.IdEdit, http.StatusInternalServerError))

		assert.IsError(t, err, publisher.ErrStartEdit)
		assert.Equals(t, *calls, []string{client.IdEdit})
	})

	t.Run("should return publisher.ErrUploadAab", func(t *testing.T) {
		t.Parallel()

		var (
			app            = "../testdata/app-release.aab"
			track          = "internal"
			serviceAccount = "../testdata/google-service-account.json"
		)
		calls, err := upload(app, track, serviceAccount, client.NewIdStatusCode(client.IdUploadAab, http.StatusInternalServerError))

		assert.IsError(t, err, publisher.ErrUploadAab)
		assert.Equals(t, *calls, []string{client.IdEdit, client.IdUploadAab})
	})

	t.Run("should return publisher.ErrUploadApk", func(t *testing.T) {
		t.Parallel()

		var (
			app            = "../testdata/app-release.apk"
			track          = "internal"
			serviceAccount = "../testdata/google-service-account.json"
		)
		calls, err := upload(app, track, serviceAccount, client.NewIdStatusCode(client.IdUploadApk, http.StatusInternalServerError))

		assert.IsError(t, err, publisher.ErrUploadApk)
		assert.Equals(t, *calls, []string{client.IdEdit, client.IdUploadApk})
	})

	t.Run("should return publisher.ErrUpdateTrack", func(t *testing.T) {
		t.Parallel()

		var (
			app            = "../testdata/app-release.aab"
			track          = "internal"
			serviceAccount = "../testdata/google-service-account.json"
		)
		calls, err := upload(app, track, serviceAccount, client.NewIdStatusCode(client.IdTrack, http.StatusInternalServerError))

		assert.IsError(t, err, publisher.ErrUpdateTrack)
		assert.Equals(t, *calls, []string{client.IdEdit, client.IdUploadAab, client.IdApp, client.IdTrack})
	})

	t.Run("should return publisher.ErrCommitEdit", func(t *testing.T) {
		t.Parallel()

		var (
			track          = "internal"
			serviceAccount = "../testdata/google-service-account.json"
			app            = "../testdata/app-release.aab"
		)
		calls, err := upload(app, track, serviceAccount, client.NewIdStatusCode(client.IdCommit, http.StatusInternalServerError))

		assert.IsError(t, err, publisher.ErrCommitEdit)
		assert.Equals(t, *calls, []string{client.IdEdit, client.IdUploadAab, client.IdApp, client.IdTrack, client.IdCommit})
	})
}

func TestToStatus(t *testing.T) {
	t.Run("should return publisher.StatusCompleted", func(t *testing.T) {
		t.Parallel()

		status, err := publisher.ToStatus("completed")

		assert.NoError(t, err)
		assert.Equals(t, *status, publisher.StatusCompleted)
	})

	t.Run("should return publisher.StatusInProgress", func(t *testing.T) {
		t.Parallel()

		status, err := publisher.ToStatus("inProgress")

		assert.NoError(t, err)
		assert.Equals(t, *status, publisher.StatusInProgress)
	})

	t.Run("should return publisher.StatusDraft", func(t *testing.T) {
		t.Parallel()

		status, err := publisher.ToStatus("draft")

		assert.NoError(t, err)
		assert.Equals(t, *status, publisher.StatusDraft)
	})

	t.Run("should return publisher.StatusHalted", func(t *testing.T) {
		t.Parallel()

		status, err := publisher.ToStatus("halted")

		assert.NoError(t, err)
		assert.Equals(t, *status, publisher.StatusHalted)
	})

	t.Run("should return publisher.ErrConvertStatus", func(t *testing.T) {
		t.Parallel()

		status, err := publisher.ToStatus("invalid")

		assert.IsError(t, err, publisher.ErrConvertStatus)
		assert.Equals(t, status, (*publisher.Status)(nil))
	})
}
