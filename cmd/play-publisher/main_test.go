package main

import (
	"bytes"
	"flag"
	"os"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/anselstetter/play-publisher/internal/assert"
	"github.com/anselstetter/play-publisher/internal/publisher/client"
)

// go test -run ^TestMain$ --generate-golden-files
var generate = flag.Bool("generate-golden-files", false, "Generate golden files")

func TestMain(t *testing.T) {
	aab := "../../internal/testdata/app-release.aab"
	apk := "../../internal/testdata/app-release.apk"
	serviceAccount := "../../internal/testdata/google-service-account.json"
	tests := []struct {
		name      string
		args      []string
		exitCode  int
		substring bool
		golden    string
	}{
		{
			name:      "should display help with exit code 0",
			args:      []string{"help"},
			exitCode:  0,
			substring: false,
			golden:    "help",
		},
		{
			name:      "should display help for upload command with exit code 0",
			args:      []string{"help", "upload"},
			exitCode:  0,
			substring: false,
			golden:    "help-upload",
		},
		{
			name:      "should display help for info command with exit code 0",
			args:      []string{"help", "info"},
			exitCode:  0,
			substring: false,
			golden:    "help-info",
		},
		{
			name:      "should display help for version command with exit code 0",
			args:      []string{"help", "version"},
			exitCode:  0,
			substring: false,
			golden:    "help-version",
		},
		{
			name:      "should display version with exit code 0",
			args:      []string{"version"},
			exitCode:  0,
			substring: false,
			golden:    "version",
		},
		{
			name:      "should display infos for provided AAB with exit code 0",
			args:      []string{"info", aab},
			exitCode:  0,
			substring: false,
			golden:    "info-aab",
		},
		{
			name:      "should display infos with additional args for provided AAB with exit code 0",
			args:      []string{"info", aab, "extra"},
			exitCode:  0,
			substring: false,
			golden:    "info-aab-additional",
		},
		{
			name:      "should display infos for provided APK with exit code 0",
			args:      []string{"info", apk},
			exitCode:  0,
			substring: false,
			golden:    "info-apk",
		},
		{
			name:      "should display infos with additional args for provided APK with exit code 0",
			args:      []string{"info", apk, "extra"},
			exitCode:  0,
			substring: false,
			golden:    "info-apk-additional",
		},
		{
			name:      "should display missing filename error for info command with exit code 1",
			args:      []string{"info"},
			exitCode:  1,
			substring: false,
			golden:    "info-missing-filename",
		},
		{
			name:      "should display invalid app error for info command with exit code 1",
			args:      []string{"info", "invalid"},
			exitCode:  1,
			substring: false,
			golden:    "info-invalid-app",
		},
		{
			name:      "should succeed upload with exit code 0",
			args:      []string{"upload", aab, "--service-account", serviceAccount},
			exitCode:  0,
			substring: true,
			golden:    "upload",
		},
		{
			name:      "should succeed upload with additional args with exit code 0",
			args:      []string{"upload", aab, "--service-account", serviceAccount, "extra"},
			exitCode:  0,
			substring: true,
			golden:    "upload-additional",
		},
		{
			name:      "should display missing service account error for upload command with exit code 1",
			args:      []string{"upload", aab},
			exitCode:  1,
			substring: false,
			golden:    "upload-missing-service_account",
		},
		{
			name:      "should display read service account error for upload command with exit code 1",
			args:      []string{"upload", aab, "--service-account", "invalid"},
			exitCode:  1,
			substring: false,
			golden:    "upload-invalid-service-account",
		},
		{
			name:      "should display invalid app error for upload command with exit code 1",
			args:      []string{"upload", "invalid", "--service-account", serviceAccount},
			exitCode:  1,
			substring: false,
			golden:    "upload-invalid-app",
		},
	}

	stubClient, _ := client.NewStubClient()

	if *generate {
		for _, tc := range tests {
			buffer := bytes.NewBuffer([]byte{})
			_ = run(tc.args, buffer, buffer, buildInfo, stubClient)

			err := os.WriteFile("testdata/"+tc.golden+".golden", buffer.Bytes(), 0666)
			assert.NoError(t, err)
		}
		return
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			buffer := bytes.NewBuffer([]byte{})
			exitCode := run(tc.args, buffer, buffer, buildInfo, stubClient)
			golden, _ := os.ReadFile("testdata/" + tc.golden + ".golden")

			if tc.substring {
				assert.True(t, strings.Contains(buffer.String(), string(golden)))
			} else {
				assert.Equals(t, buffer.String(), string(golden))
			}
			assert.Equals(t, exitCode, tc.exitCode)
		})
	}
}

func buildInfo() (info *debug.BuildInfo, ok bool) {
	buildInfo := &debug.BuildInfo{
		Settings: []debug.BuildSetting{
			{Key: "-ldflags", Value: "-s -w -X main.Version=Test -s"},
		},
	}
	return buildInfo, true
}
