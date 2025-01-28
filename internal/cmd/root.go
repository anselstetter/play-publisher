package cmd

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "play-publisher",
		Short: "Upload an app to the Play Store",
		Long:  "A Play Store uploader for your convenience",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		SilenceErrors: true,
	}
	rootCmd.Root().CompletionOptions.DisableDefaultCmd = true

	return rootCmd
}
