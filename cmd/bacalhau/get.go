package bacalhau

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/filecoin-project/bacalhau/pkg/ipfs"
	"github.com/filecoin-project/bacalhau/pkg/system"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var getCmdFlags = struct {
	timeoutSecs int
	outputDir   string
}{
	timeoutSecs: 60,
	outputDir:   ".",
}

func init() { // nolint:gochecknoinits // Using init in cobra command is idomatic
	getCmd.Flags().IntVar(&getCmdFlags.timeoutSecs, "timeout-secs",
		getCmdFlags.timeoutSecs, "Timeout duration for IPFS downloads.")
	getCmd.Flags().StringVar(&getCmdFlags.outputDir, "output-dir",
		getCmdFlags.outputDir, "Directory to write the output to.")
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the results of a job",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error { // nolintunparam // incorrectly suggesting unused
		cm := system.NewCleanupManager()
		defer cm.Cleanup()

		clientID, err := system.GetClientID()
		if err != nil {
			log.Error().Msgf("Failed to get client ID: %s", err)
			return err
		}

		log.Info().Msgf("Fetching results of job '%s'...", args[0])
		job, ok, err := getAPIClient().Get(context.Background(), clientID, args[0])
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("job not found")
		}

		var resultCIDs []string
		for _, jobState := range job.State {
			if jobState.ResultsID != "" {
				resultCIDs = append(resultCIDs, jobState.ResultsID)
			}
		}
		log.Debug().Msgf("Job has result CIDs: %v", resultCIDs)

		log.Debug().Msg("Spinning up IPFS client...")
		cl, err := ipfs.NewClient(cm)
		if err != nil {
			return err
		}

		for _, cid := range resultCIDs {
			outputDir := filepath.Join(getCmdFlags.outputDir, cid)
			log.Info().Msgf("Downloading result CID '%s' to '%s'...",
				cid, outputDir)

			ctx, cancel := context.WithDeadline(context.Background(),
				time.Now().Add(time.Second*time.Duration(getCmdFlags.timeoutSecs)))
			defer cancel()

			err = cl.Get(ctx, cid, outputDir)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					log.Error().Msg("Timed out while downloading result.")
				}

				return err
			}
		}

		return nil
	},
}
