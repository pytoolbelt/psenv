// /*
// Copyright © 2024 Jesse Maitland jesse@pytoolbelt.com
// */
package cmd

import (
	"fmt"
	"github.com/pytoolbelt/psenv/internal/config"
	"github.com/pytoolbelt/psenv/internal/parameterstore"
	"github.com/pytoolbelt/psenv/internal/utils"
	"github.com/spf13/cobra"
	"os"
	"sync"
	"time"
)

var overwriteFlag bool
var keyIDFlag string

func putWorker(paramsToAdd map[string]string, wg *sync.WaitGroup) {
	defer wg.Done()
	// put the parameters in the parameter store
	ps, err := parameterstore.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// put the parameters in the parameter store
	err = ps.PutParameters(paramsToAdd, keyIDFlag, overwriteFlag)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func deleteWorker(paramsToDelete []string, wg *sync.WaitGroup) {
	defer wg.Done()
	// delete the parameters in the parameter store
	ps, err := parameterstore.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// delete the parameters in the parameter store
	err = ps.DeleteParameters(paramsToDelete)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getWorker(envChan <-chan string, wg *sync.WaitGroup, paramsChan chan<- map[string]string, secretsConfig *config.SecretsConfig) {
	defer wg.Done()

	// get the parameter store. If we can't make one for some reason,
	// just exit as there is nothing to do.
	ps, err := parameterstore.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for env := range envChan {
		// get the parameters for a given environment path
		path := secretsConfig.GetEnvironmentPath(env)
		remoteParams, err := ps.GetParameters(path, decryptFlag)
		if err != nil {
			fmt.Println(err)
			continue
		}

		paramsChan <- remoteParams
	}
}

func putEntrypoint(cmd *cobra.Command, args []string) {

	var envChan = make(chan string, 10)
	var paramsChan = make(chan map[string]string, 25)
	var errorChan = make(chan error, 10)
	var wg sync.WaitGroup
	var numberOfWorkers int = 1
	var environmentsToPut []string

	fmt.Println("putting parameters in the parameter store")
	// we are doing a put operation so load the secrets config file.
	// if one is not found, just exit as there is nothing to do.
	secretsConfig, err := config.LoadSecretsConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if envName == "" {
		numberOfWorkers = len(secretsConfig.Environments)
	}

	if numberOfWorkers == 0 {
		fmt.Println("no environments found in the psenv-secrets.yml file")
		os.Exit(1)
	}

	// start the workers to put the parameters
	for i := 0; i < numberOfWorkers; i++ {
		wg.Add(1)
		go putMainWorker(envChan, errorChan, &wg, secretsConfig)
	}

	// put the configured paths on the channel. These will be used
	// by the parameter store to get params by path.
	if envName != "" {
		environmentsToPut = append(environmentsToPut, envName)
	} else {
		for env := range secretsConfig.Environments {
			environmentsToPut = append(environmentsToPut, env)
		}
	}

	// put the envs on the channel
	for _, env := range environmentsToPut {
		envChan <- env
	}

	// close the path channel as we are done with it.
	close(envChan)

	// wait for the putWorker workers to finish
	wg.Wait()

	// close the error channel and check to see if we got any errors from
	// the putWorker workers. If we did, print them out and exit.
	close(errorChan)
	for err := range errorChan {
		fmt.Println(err)
		os.Exit(1)
	}

	// it seems we got no errors, so proceed to execute the get workers
	time.Sleep(5 * time.Second)

	envChan = make(chan string, 10)

	// start the get workers to get the parameters
	for i := 0; i < numberOfWorkers; i++ {
		wg.Add(1)
		go getWorker(envChan, &wg, paramsChan, secretsConfig)
	}

	// put the environments back on the channel
	for _, env := range environmentsToPut {
		envChan <- env
	}

	// close the path channel as we are done with it.
	close(envChan)

	// wait for the getWorker workers to finish
	wg.Wait()

	// close the params channel
	close(paramsChan)

	// clear all the environments that we are processing
	for _, env := range environmentsToPut {
		secretsConfig.ClearEnvironment(env)
	}
	// update the config with the fetch fresh params
	for params := range paramsChan {
		err = secretsConfig.UpdateSecretsConfigFromParameters(params)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	err = secretsConfig.Save()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func putMainWorker(envChan <-chan string, errorChan chan<- error, wg *sync.WaitGroup, secretsConfig *config.SecretsConfig) {
	var swg sync.WaitGroup
	defer wg.Done()

	// get the parameter store. If we can't make one for some reason,
	// just exit as there is nothing to do.
	ps, err := parameterstore.New()
	if err != nil {
		fmt.Println(err)
		errorChan <- err
		os.Exit(1)
	}

	for env := range envChan {
		// first get the parameters for a given environment path
		path := secretsConfig.GetEnvironmentPath(env)
		remoteParams, err := ps.GetParameters(path, false)

		if err != nil {
			errorChan <- err
			continue
		}

		localParams := secretsConfig.GetEnvironmentParams(env)

		// create a param map to determine what we need to do.
		parameters := utils.MergeLocalAndRemoteParams(localParams, remoteParams)

		// first if we have any parameters to add, just add them.
		if len(parameters.ToAdd) > 0 {
			go putWorker(parameters.ToAdd, &swg)
			swg.Add(1)
		} else {
			fmt.Printf("no parameters to add to environment %s\n", env)
		}

		// if we have any parameters to delete, just delete them.
		if len(parameters.ToDelete) > 0 {
			go deleteWorker(parameters.ToDelete, &swg)
			swg.Add(1)
		} else {
			fmt.Printf("no parameters to update in environment %s\n", env)
		}

		// if we have any parameters to update, just update them.
		if len(parameters.ToUpdate) > 0 {
			go putWorker(parameters.ToUpdate, &swg)
			swg.Add(1)
		} else {
			fmt.Printf("no parameters to delete from environment %s\n", env)
		}

		// wait for the sub workers to finish
		swg.Wait()
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
	putCmd.Flags().BoolVarP(&overwriteFlag, "overwrite", "o", false, "overwrite existing parameters")
	putCmd.Flags().StringVarP(&keyIDFlag, "kms-name", "k", "alias/aws/ssm", "KMS key name to use for encryption")
}
