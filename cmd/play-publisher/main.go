package main

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/anselstetter/play-publisher/internal/application"
	"github.com/anselstetter/play-publisher/internal/cmd"
	"github.com/anselstetter/play-publisher/internal/logger"
	"github.com/anselstetter/play-publisher/internal/publisher"
)

var Version string = "Dev"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr, nil))
}

func run(args []string, stdout io.Writer, stderr io.Writer, client *http.Client) int {
	var (
		ctx       = context.Background()
		analyzer  = application.New()
		logger    = logger.New(logger.WithStdout(stdout), logger.WithStderr(stderr))
		publisher = publisher.New(analyzer, logger, publisher.WithHttpClient(client))
		root      = cmd.NewRootCommand()
	)
	root.AddCommand(
		cmd.NewVersionCommand(Version, logger),
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
