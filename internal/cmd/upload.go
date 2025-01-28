package cmd

import (
	"fmt"

	"github.com/anselstetter/play-publisher/internal/logger"
	"github.com/anselstetter/play-publisher/internal/publisher"
	"github.com/spf13/cobra"
)

const synopsisUpload = `Upload an APK or AAB to the Play Store.

The service account, which has been linked to the Play Store, is needed.
The track is optional and defaults to "internal", if omitted`

func NewUploadCommand(publisher publisher.Publisher, logger logger.Logger) *cobra.Command {
	var (
		serviceAccount string
		track          string
	)
	uploadCmd := &cobra.Command{
		Use:    "upload <file>",
		Short:  "Upload an APK or AAB to the Play Store",
		Long:   synopsisUpload,
		PreRun: ignoreAdditonalArgs(logger.StderrWriter(), 1),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("please provide a file to upload")
			}
			return nil
		},
		RunE: silenceUsageE(func(cmd *cobra.Command, args []string) error {
			fileName := args[0]

			return publisher.Upload(cmd.Context(), fileName, track, serviceAccount)
		}),
	}
	uploadCmd.Flags().StringVarP(&serviceAccount, "service-account", "s", "", "the service account (required)")
	uploadCmd.Flags().StringVarP(&track, "track", "t", "internal", "the track")

	if err := uploadCmd.MarkFlagRequired("service-account"); err != nil {
		logger.Stderrf("Could not mark service-account as required: %s", err)
	}
	return uploadCmd
}
