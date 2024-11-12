/*
Copyright © 2024 Jesse Maitland jesse@pytoolbelt.com
*/
package cmd

import (
	"os"
	"sync"
	"time"

	"github.com/pytoolbelt/psenv/internal/config"
	"github.com/pytoolbelt/psenv/internal/paramstore"
	"github.com/spf13/cobra"
)

var mergeFlag bool
var overwriteFlag bool

func putEntrypoint(cmd *cobra.Command, args []string) {
	var err error
	var wg sync.WaitGroup

	cfg, err := config.InitAndLoad()
	HandelError(err)

	envsToProcess := cfg.GetEnvironments() // Assuming this method returns a list of environments
	envChan := make(chan config.Environment, len(envsToProcess))
	resultChan := make(chan config.Environment, len(envsToProcess))

	// Start worker Go routines
	for i := 0; i <= 3; i++ { // Number of workers
		wg.Add(1)
		go putWorker(envChan, resultChan, &wg)
	}

	// Send environments to the channel
	for _, env := range envsToProcess {
		envChan <- env
	}
	close(envChan)

	// Wait for all workers to finish
	wg.Wait()
	close(resultChan)

	// Collect results
	for env := range resultChan {
		cfg.SetEnvironment(&env)
	}

	err = cfg.Save()
	HandelError(err)
	os.Exit(0)
}

func putWorker(envChan <-chan config.Environment, resultChan chan<- config.Environment, wg *sync.WaitGroup) {
	defer wg.Done()

	for env := range envChan {
		psPath := env.GetParamStorePath()
		ps, err := paramstore.NewParamStore(psPath)
		if err != nil {
			HandelError(err)
			continue
		}

		remoteParams, err := ps.GetParameters(false)
		if err != nil {
			HandelError(err)
			continue
		}

		paramsToAdd := env.GetParamsToAdd(remoteParams)
		err = ps.PutParameters(paramsToAdd, overwriteFlag)
		if err != nil {
			HandelError(err)
			continue
		}

		// if we are merging, see if we have any parameters to delete
		if mergeFlag {
			paramsToDelete := env.GetParamsToDelete(remoteParams)

			// if we have parameters to delete, delete them. This means that we have a parameter in the remote store that is not in the local store
			if len(paramsToDelete) > 0 {
				err = ps.DeleteParameters(paramsToDelete)
				if err != nil {
					HandelError(err)
					continue
				}
			}
		}

		time.Sleep(5 * time.Second)
		remoteParams, err = ps.GetParameters(false)
		if err != nil {
			HandelError(err)
			continue
		}

		env.Params = remoteParams
		resultChan <- env
	}
}

// putCmd represents the put command
var putCmd = &cobra.Command{
	Use:   "put",
	Short: "put parameters in the parameter store",
	Long:  ``,
	Run:   putEntrypoint,
}

func init() {
	rootCmd.AddCommand(putCmd)
	putCmd.Flags().BoolVarP(&mergeFlag, "merge", "m", false, "mode to push parameters in")
	putCmd.Flags().BoolVarP(&overwriteFlag, "overwrite", "o", false, "overwrite existing parameters")
}
