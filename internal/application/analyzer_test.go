package application_test

import (
	"testing"

	"github.com/anselstetter/play-publisher/internal/application"
	"github.com/anselstetter/play-publisher/internal/assert"
)

func TestApplication(t *testing.T) {
	t.Run("should return application.Info for aab", func(t *testing.T) {
		t.Parallel()

		app := "../testdata/app-release.aab"
		want := application.ApplicationInfo{
			ApplicationType: application.ApplicationTypeAab,
			PackageName:     "tld.domain.test",
			VersionName:     "1.0",
			VersionCode:     1,
		}
		got, err := application.New().Analyze(app)
		assert.NoError(t, err)
		assert.Equals(t, got, want)
	})

	t.Run("should return application.Info for apk", func(t *testing.T) {
		t.Parallel()

		app := "../testdata/app-release.apk"
		want := application.ApplicationInfo{
			ApplicationType: application.ApplicationTypeApk,
			PackageName:     "tld.domain.test",
			VersionName:     "1.0",
			VersionCode:     1,
		}
		got, err := application.New().Analyze(app)
		assert.NoError(t, err)
		assert.Equals(t, got, want)
	})

	t.Run("should return application.ErrAnalyze", func(t *testing.T) {
		t.Parallel()

		app := "not-existing"
		_, err := application.New().Analyze(app)
		assert.IsError(t, err, application.ErrAnalyze)
	})
}

func TestType(t *testing.T) {
	t.Run("should return AAB for application.TypeAab", func(t *testing.T) {
		t.Parallel()

		assert.Equals(t, application.ApplicationTypeAab.String(), "AAB")
	})

	t.Run("should return APK for application.TypeApk", func(t *testing.T) {
		t.Parallel()

		assert.Equals(t, application.ApplicationTypeApk.String(), "APK")
	})

	t.Run("should return Unknown for application.TypeUnknown", func(t *testing.T) {
		t.Parallel()

		assert.Equals(t, application.ApplicationTypeUnknown.String(), "Unknown")
	})
}
