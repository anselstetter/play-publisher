package cmd

import (
	"fmt"

	"github.com/anselstetter/play-publisher/internal/application"
	"github.com/anselstetter/play-publisher/internal/logger"
	"github.com/spf13/cobra"
)

const synopsisInfo = `Prints information about the provided app in the following form:

File name:        {file name}
Application type: {aab|apk}
Package name:     {package name}
Version name:     {version name}
Version code:     {version code}`

func NewInfoCommand(analyzer application.Analyzer, logger logger.Logger) *cobra.Command {
	infoCmd := &cobra.Command{
		Use:    "info <file>",
		Short:  "Display the package name and type of the app",
		Long:   synopsisInfo,
		PreRun: ignoreAdditonalArgs(logger.StderrWriter(), 1),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("please provide a file to analyze")
			}
			return nil
		},
		RunE: silenceUsageE(func(cmd *cobra.Command, args []string) error {
			fileName := args[0]
			return Info(fileName, analyzer, logger)
		}),
	}
	return infoCmd
}

func Info(fileName string, analyzer application.Analyzer, logger logger.Logger) error {
	applicationInfo, err := analyzer.Analyze(fileName)
	if err != nil {
		return err
	}
	logger.StdoutTable(
		"Application type", applicationInfo.ApplicationType,
		"Package name", applicationInfo.PackageName,
		"Version name", applicationInfo.VersionName,
		"Version code", applicationInfo.VersionCode,
	)
	return nil
}
