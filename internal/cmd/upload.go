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
		status         string
	)
	cmd := &cobra.Command{
		Use:    "upload <file>",
		Short:  "Upload an APK or AAB to the Play Store",
		Long:   synopsisUpload,
		PreRun: ignoreAdditonalArgs(logger.StderrWriter(), 1),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("please provide a file to upload")
			}
			if _, err := toStatus(status); err != nil {
				return err
			}
			return nil
		},
		RunE: silenceUsageE(func(cmd *cobra.Command, args []string) error {
			fileName := args[0]
			status, _ := toStatus(status)

			return publisher.Upload(cmd.Context(), fileName, track, *status, serviceAccount)
		}),
	}
	cmd.Flags().StringVarP(&serviceAccount, "service-account", "s", "", "the service account (required)")
	cmd.Flags().StringVarP(&track, "track", "t", "internal", "the track")
	cmd.Flags().StringVarP(&status, "status", "S", "completed", "status (completed, inProgress, draft, halted)")

	if err := cmd.MarkFlagRequired("service-account"); err != nil {
		logger.Stderrf("Could not mark service-account as required: %s", err)
	}
	return cmd
}

func toStatus(status string) (*publisher.Status, error) {
	return publisher.ToStatus(status)
}
