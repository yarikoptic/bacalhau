package bacalhau

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/filecoin-project/bacalhau/pkg/executor"
	pjob "github.com/filecoin-project/bacalhau/pkg/job"
	"github.com/filecoin-project/bacalhau/pkg/storage"
	"github.com/filecoin-project/bacalhau/pkg/system"
	"github.com/filecoin-project/bacalhau/pkg/verifier"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var jobspec *executor.JobSpec
var filename string
var jobfConcurrency int
var jobfInputUrls []string
var jobfInputVolumes []string
var jobfOutputVolumes []string
var jobTags []string

func init() { // nolint:gochecknoinits
	applyCmd.PersistentFlags().StringVarP(
		&filename, "filename", "f", "",
		`Path to the job file`,
	)

	applyCmd.PersistentFlags().IntVarP(
		&jobfConcurrency, "concurrency", "c", 1,
		`How many nodes should run the job in parallel`,
	)

	applyCmd.PersistentFlags().BoolVarP(
		&waitForJobToFinishAndPrintOutput, "wait", "w", false,
		`Wait For Job To Finish And Print Output`,
	)

	applyCmd.PersistentFlags().StringSliceVarP(&jobTags,
		"labels", "l", []string{},
		`List of jobTags for the job. In the format 'a,b,c,1'. All characters not matching /a-zA-Z0-9_:|-/ and all emojis will be stripped.`,
	)
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Submit a job.json or job.yaml file and run it on the network",
	Args:  cobra.MinimumNArgs(0),
	PostRun: func(cmd *cobra.Command, args []string) {
		// Can't think of any reason we'd want these to persist.
		// The below is to clean out for testing purposes. (Kinda ugly to put it in here,
		// but potentially cleaner than making things public, which would
		// be the other way to attack this.)
		jobInputs = []string{}
		jobInputUrls = []string{}
		jobInputVolumes = []string{}
		jobOutputVolumes = []string{}
		jobEnv = []string{}
		jobLabels = []string{}

		jobEngine = "docker"
		jobVerifier = "ipfs"
		jobConcurrency = 1
		jobCPU = ""
		jobMemory = ""
		jobGPU = ""
		skipSyntaxChecking = false
		waitForJobToFinishAndPrintOutput = false
		jobIpfsGetTimeOut = 10
	},
	RunE: func(cmd *cobra.Command, cmdArgs []string) error { // nolintunparam // incorrect that cmd is unused.
		ctx := context.Background()
		fileextension := filepath.Ext(filename)

		fileContent, err := os.Open(filename)

		if err != nil {
			return err
		}

		defer fileContent.Close()

		byteResult, err := io.ReadAll(fileContent)

		if err != nil {
			return err
		}

		if fileextension == ".json" {
			err = json.Unmarshal(byteResult, &jobspec)
			if err != nil {
				return err
			}
		}

		if fileextension == ".yaml" || fileextension == ".yml" {
			err = yaml.Unmarshal(byteResult, &jobspec)
			if err != nil {
				return err
			}
		}

		jobImage := jobspec.Docker.Image

		jobEntrypoint := jobspec.Docker.Entrypoint

		if len(jobspec.Inputs) != 0 {
			for _, jobspecInput := range jobspec.Inputs {
				var storageSpecEngineType storage.StorageSourceType
				storageSpecEngineType, err = storage.ParseStorageSourceType(jobspecInput.EngineName)
				if err != nil {
					return err
				}
				if jobspecInput.Path == "" {
					return fmt.Errorf("empty volume mount point %+v", jobspecInput)
				}
				if storageSpecEngineType == storage.StorageSourceIPFS {
					if jobspecInput.Cid == "" {
						return fmt.Errorf("empty ipfs volume cid %+v", jobspecInput)
					}
					is := jobspecInput.Cid + ":" + jobspecInput.Path
					jobfInputVolumes = append(jobfInputVolumes, is)
				} else if storageSpecEngineType == storage.StorageSourceURLDownload {
					if jobspecInput.URL == "" {
						return fmt.Errorf("empty url volume url %+v", jobspecInput)
					}
					is := jobspecInput.URL + ":" + jobspecInput.Path
					jobfInputUrls = append(jobfInputUrls, is)
				} else {
					return fmt.Errorf("unknown storage source type %s", jobspecInput.EngineName)
				}
			}
		}

		if len(jobspec.Outputs) != 0 {
			for _, jobspecsOutputs := range jobspec.Outputs {
				is := jobspecsOutputs.Name + ":" + jobspecsOutputs.Path
				jobfOutputVolumes = append(jobfOutputVolumes, is)
			}
		}

		engineType, err := executor.ParseEngineType(jobspec.EngineName)
		if err != nil {
			cmd.Printf("Error parsing engine type: %s", err)
			return err
		}

		verifierType, err := verifier.ParseVerifierType(jobspec.VerifierName)
		if err != nil {
			cmd.Printf("Error parsing engine type: %s", err)
			return err
		}

		spec, deal, err := pjob.ConstructDockerJob(
			engineType,
			verifierType,
			jobspec.Resources.CPU,
			jobspec.Resources.Memory,
			jobspec.Resources.GPU,
			jobfInputUrls,
			jobfInputVolumes,
			jobfOutputVolumes,
			jobspec.Docker.Env,
			jobEntrypoint,
			jobImage,
			jobfConcurrency,
			jobTags,
		)
		if err != nil {
			return err
		}

		if !skipSyntaxChecking {
			err = system.CheckBashSyntax(jobEntrypoint)
			if err != nil {
				return err
			}
		}

		job, err := getAPIClient().Submit(ctx, spec, deal, nil)
		if err != nil {
			return err
		}
		states, err := getAPIClient().GetExecutionStates(ctx, job.ID)
		if err != nil {
			return err
		}
		currentNodeID, _ := pjob.GetCurrentJobState(states)
		nodeIds := []string{currentNodeID}
		if waitForJobToFinishAndPrintOutput {
			err = WaitForJob(ctx, job.ID, job,
				WaitForJobThrowErrors(job, []executor.JobStateType{
					executor.JobStateCancelled,
					executor.JobStateError,
				}),
				WaitForJobAllHaveState(nodeIds, executor.JobStateComplete),
			)
			if err != nil {
				return err
			}

			cidl := Get(job.ID, jobIpfsGetTimeOut)
			var cidv string
			for cid := range cidl {
				cidv = cid
			}
			body, err := os.ReadFile(cidv + "/stdout")
			if err != nil {
				return err
			}
			fmt.Println()
			fmt.Println(string(body))
		}
		cmd.Printf("%s\n", job.ID)
		return nil

	},
}
