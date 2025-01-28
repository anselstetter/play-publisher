package publisher

import (
	"context"
	"errors"
	"io"
	"os"
	"time"

	"github.com/anselstetter/play-publisher/internal/application"
	"github.com/anselstetter/play-publisher/internal/logger"
	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/googleapi"
	apiOptions "google.golang.org/api/option"
)

type uploaderFunc func(service *androidpublisher.Service, applicationInfo application.ApplicationInfo, editId string, reader io.Reader, size int64) (int64, error)

const clearLine = "\033[1A\033[K"

const (
	trackReleaseStatusCompleted  = "completed"
	trackReleaseStatusInProgress = "inProgress"
	trackReleaseStatusDraft      = "draft"
	trackReleaseStatusHalted     = "halted"
)

var (
	ErrorUnhandled          = errors.New("unhandled. This is a bug")
	ErrorAnalyzeApp         = errors.New("could not analyze app")
	ErrorReadServiceAccount = errors.New("could not open service account")
	ErrorCreateService      = errors.New("could not instantiate service")
	ErrorCreateEdit         = errors.New("could create new edit")
	ErrorUploadAab          = errors.New("could not upload aab")
	ErrorUploadApk          = errors.New("could not upload apk")
	ErrorUpdateTrack        = errors.New("could not update track")
	ErrorCommitEdit         = errors.New("could not commit edit")
)

type Publisher struct {
	analyzer application.Analyzer
	logger   logger.Logger
	options  options
}

// The uploading process is more complicated, than it should be.
//
// After creating the androidpublisher.Service, Google wants some kind of session,
// which is represented by an edit id, and has to be used throughout every operation.
//
// After the app has been uploaded to a specific track, the track has to be updated,
// otherwise the app would only be visible in the bundle explorer.
//
// Eventually the edit has to be committed.
// Note: When committing the edit, a review could be triggered, which is not done here.
func (p Publisher) Upload(ctx context.Context, fileName string, track string, serviceAccount string) error {
	applicationInfo, err := p.analyzer.Analyze(fileName)
	if err != nil {
		return errors.Join(ErrorAnalyzeApp, err)
	}

	p.printInfo(applicationInfo, track, serviceAccount)

	json, err := os.ReadFile(serviceAccount)
	if err != nil {
		return errors.Join(ErrorReadServiceAccount, err)
	}

	service, err := p.createService(ctx, json)
	if err != nil {
		return errors.Join(ErrorCreateService, err)
	}

	edit, err := p.startEdit(service, applicationInfo)
	if err != nil {
		return errors.Join(ErrorCreateEdit, err)
	}
	versionCode := int64(0)

	switch applicationInfo.ApplicationType {
	case application.ApplicationTypeAab:
		versionCode, err = p.upload(service, applicationInfo, fileName, edit.Id, p.uploadAab)
		if err != nil {
			return errors.Join(ErrorUploadAab, err)
		}
	case application.ApplicationTypeApk:
		versionCode, err = p.upload(service, applicationInfo, fileName, edit.Id, p.uploadApk)
		if err != nil {
			return errors.Join(ErrorUploadApk, err)
		}
	case application.ApplicationTypeUnknown:
		return ErrorUnhandled
	}

	_, err = p.updateTrack(service, applicationInfo, track, edit, versionCode)
	if err != nil {
		return errors.Join(ErrorUpdateTrack, err)
	}

	_, err = p.commit(service, applicationInfo, edit)
	if err != nil {
		return errors.Join(ErrorCommitEdit, err)
	}

	return nil
}

func (p Publisher) printInfo(applicationInfo application.ApplicationInfo, track string, serviceAccount string) {
	p.logger.StdoutTable(
		"Application type:", applicationInfo.ApplicationType,
		"Package name:", applicationInfo.PackageName,
		"Version name:", applicationInfo.VersionName,
		"Version code:", applicationInfo.VersionCode,
		"Service account:", serviceAccount,
		"Track:", track,
	)
	p.logger.Stdoutln()
}

func (p Publisher) printProgress(current, total int64) {
	percentage := int(float32(current) / float32(total) * 100)

	p.logger.Stdoutf("\n%s", clearLine)
	p.logger.Stdoutf("%d %%", percentage)
}

func (p Publisher) createService(ctx context.Context, json []byte) (*androidpublisher.Service, error) {
	timer := p.duration("Creating service")
	defer timer(" | Done")

	return androidpublisher.NewService(ctx, apiOptions.WithCredentialsJSON(json), apiOptions.WithHTTPClient(p.options.httpClient))
}

func (p Publisher) startEdit(service *androidpublisher.Service, applicationInfo application.ApplicationInfo) (*androidpublisher.AppEdit, error) {
	timer := p.duration("Starting edit")
	defer timer(" | Done")

	return service.Edits.Insert(applicationInfo.PackageName, nil).Do()
}

func (p Publisher) updateTrack(service *androidpublisher.Service, applicationInfo application.ApplicationInfo, track string, edit *androidpublisher.AppEdit, versionCode int64) (*androidpublisher.Track, error) {
	timer := p.duration("Updating track")
	defer timer(" | Done")

	releaseTrack := androidpublisher.Track{
		Track: track,
		Releases: []*androidpublisher.TrackRelease{{
			VersionCodes: []int64{versionCode},
			Status:       trackReleaseStatusCompleted,
		}},
	}
	return service.Edits.Tracks.
		Update(applicationInfo.PackageName, edit.Id, track, &releaseTrack).
		Do()
}

func (p Publisher) commit(service *androidpublisher.Service, applicationInfo application.ApplicationInfo, edit *androidpublisher.AppEdit) (*androidpublisher.AppEdit, error) {
	timer := p.duration("Committing edit")
	defer timer(" | Done")

	return service.Edits.Commit(applicationInfo.PackageName, edit.Id).
		ChangesNotSentForReview(true).
		Do()
}

// This function takes a higher order function to do the actual work.
// In this case, one for uploading an apk, and another one to upload an aab.
//
// Note: There is a progress updater which takes another function,
// with the current uploaded bytes and the total bytes.
// Unfortunately, the total bytes are never determined in the androidpublisher lib,
// and is always 0, which renders the updater function useless.
// As a workaround, the size is determined here and passed to the uploader function,
// which can then be used to dispatch to the outside progress function,
// to render the correct percentage.
func (p Publisher) upload(service *androidpublisher.Service, applicationInfo application.ApplicationInfo, fileName string, editId string, uploader uploaderFunc) (int64, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer func(fs *os.File) {
		if err := fs.Close(); err != nil {
			p.logger.Stderrf("Warning: %s", err)
		}
	}(file)

	info, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return uploader(service, applicationInfo, editId, file, info.Size())
}

func (p Publisher) uploadAab(service *androidpublisher.Service, applicationInfo application.ApplicationInfo, editId string, reader io.Reader, size int64) (int64, error) {
	timer := p.duration("Uploading AAB\n")
	defer timer(" | Done")

	bundle, err := service.Edits.Bundles.
		Upload(applicationInfo.PackageName, editId).
		Media(reader, googleapi.ContentType("application/octet-stream"), googleapi.ChunkSize(googleapi.MinUploadChunkSize)).
		ProgressUpdater(func(current, _ int64) { p.printProgress(current, size) }).
		Do()

	if err != nil {
		return 0, err
	}
	return bundle.VersionCode, nil
}

// This has never been tested, since aab is kinda standard now.
func (p Publisher) uploadApk(service *androidpublisher.Service, applicationInfo application.ApplicationInfo, editId string, reader io.Reader, size int64) (int64, error) {
	timer := p.duration("Uploading APK\n")
	defer timer(" | Done")

	apk, err := service.Edits.Apks.
		Upload(applicationInfo.PackageName, editId).
		Media(reader, googleapi.ContentType("application/octet-stream"), googleapi.ChunkSize(googleapi.MinUploadChunkSize)).
		ProgressUpdater(func(current, _ int64) { p.printProgress(current, size) }).
		Do()
	if err != nil {
		return 0, err
	}
	return apk.VersionCode, nil
}

func (p Publisher) duration(info string) func(name string) {
	start := time.Now()
	p.logger.Stdoutf("%s", info)

	return func(name string) {
		elapsed := time.Since(start)

		p.logger.Stdoutf("%s (took %s)\n", name, elapsed)
	}
}

func New(analyzer application.Analyzer, logger logger.Logger, opts ...option) Publisher {
	options := newOptions(opts...)

	return Publisher{
		analyzer: analyzer,
		logger:   logger,
		options:  options,
	}
}
