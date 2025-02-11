package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/anselstetter/play-publisher/internal/application"
	"github.com/anselstetter/play-publisher/internal/cmd"
	"github.com/anselstetter/play-publisher/internal/logger"
	"github.com/anselstetter/play-publisher/internal/publisher"
	"github.com/anselstetter/play-publisher/internal/version"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr, debug.ReadBuildInfo, nil))
}

func run(args []string, stdout io.Writer, stderr io.Writer, buildInfoFunc version.BuildInfoFunc, client *http.Client) int {
	var (
		ctx       = context.Background()
		version   = version.New(buildInfoFunc, "Dev")
		analyzer  = application.New()
		logger    = logger.New(logger.WithStdout(stdout), logger.WithStderr(stderr))
		publisher = publisher.New(analyzer, logger, publisher.WithHttpClient(client))
		root      = cmd.NewRootCommand()
	)
	root.AddCommand(
		cmd.NewVersionCommand(version, logger),
		cmd.NewInfoCommand(analyzer, logger),
		cmd.NewUploadCommand(publisher, logger),
		cmd.NewGenerateDocsCommand(),
	)
	root.SetArgs(args)
	root.SetOut(stdout)
	root.SetErr(stderr)

	if err := root.ExecuteContext(ctx); err != nil {
		return fail(logger, err, 1)
	}
	return 0
}

func fail(logger logger.Logger, err error, exitCode int) int {
	logger.Stderrf("%s\n", err.Error())
	return exitCode
}
