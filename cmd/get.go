// /*
// Copyright © 2024 Jesse Maitland jesse@pytoolbelt.com
// */
package cmd

import (
	"fmt"
	"github.com/pytoolbelt/psenv/internal/config"
	"github.com/pytoolbelt/psenv/internal/parameterstore"
	"github.com/spf13/cobra"
	"os"
	"sync"
)

func getEntryPoint(cmd *cobra.Command, args []string) {
	var envChan = make(chan string, 10)
	var paramsChan = make(chan map[string]string, 25)
	var errorChan = make(chan error, 10)
	var wg sync.WaitGroup
	var numberOfWorkers int = 1
	var environmentsToGet []string
	var secretsConfig *config.SecretsConfig

	fmt.Println("getting parameters from the parameter store")

	projectConfig, err := config.LoadProjectConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	secretsConfig, err = config.LoadSecretsConfig()
	if err != nil {
		secretsConfig, err = config.CreateNewSecretsConfigFile()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if envName == "" {
		numberOfWorkers = len(projectConfig.Environments)
	}

	if numberOfWorkers == 0 {
		fmt.Println("no environments found in the psenv-project.yml file")
		os.Exit(1)
	}

	// start the workers to put the parameters
	for i := 0; i < numberOfWorkers; i++ {
		wg.Add(1)
		go mainGetWorker(envChan, errorChan, paramsChan, &wg, projectConfig, getCommandDecryptFlag)
	}

	// put the configured paths on the channel. These will be used
	// by the parameter store to get params by path.
	if envName != "" {
		environmentsToGet = append(environmentsToGet, envName)
	} else {
		for _, env := range projectConfig.Environments {
			environmentsToGet = append(environmentsToGet, env)
		}
	}

	// put the environments on the channel
	for _, env := range environmentsToGet {
		envChan <- env
	}

	// close the channel so the workers know when they are done
	close(envChan)

	// wait for the workers to finish
	wg.Wait()

	// close the error channel and check for errors
	close(errorChan)
	for err := range errorChan {
		fmt.Println(err)
		os.Exit(1)
	}

	// close the params channel
	close(paramsChan)

	// clear all the environments from the secrets file
	secretsConfig.ClearEnvironments()

	// get the params and update the secrets file
	for params := range paramsChan {
		err = secretsConfig.UpdateSecretsConfigFromParameters(params)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// save the secrets file
	err = secretsConfig.Save()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func mainGetWorker(envChan <-chan string, errorChan chan<- error, paramsChan chan<- map[string]string, wg *sync.WaitGroup, projectConfig *config.ProjectConfig, decrypt bool) {
	defer wg.Done()

	// get the parameter store. If we can't make one for some reason,
	// just exit as there is nothing to do.
	ps, err := parameterstore.New()
	if err != nil {
		errorChan <- err
		return
	}

	for env := range envChan {
		// get the parameters for a given environment path
		path := projectConfig.GetEnvironmentPath(env)
		remoteParams, err := ps.GetParameters(path, decrypt)
		if err != nil {
			errorChan <- err
			return
		}
		paramsChan <- remoteParams
	}
}

var getCommandDecryptFlag bool

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get parameters from the AWS Parameter Store",
	Long:  ``,
	Run:   getEntryPoint,
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().BoolVarP(&getCommandDecryptFlag, "decrypt", "d", false, "Decrypt secure string parameters")
}
