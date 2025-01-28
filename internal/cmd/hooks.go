package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
)

func ignoreAdditonalArgs(stderr io.Writer, n int) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		additionalArgs := args[n:]

		if len(additionalArgs) >= n {
			if _, err := fmt.Fprintf(stderr, "Ignoring additional args: %s\n\n", strings.Join(additionalArgs, ", ")); err != nil {
				// Ignore
			}
		}
	}
}

func silenceUsageE(f func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		return f(cmd, args)
	}
}
