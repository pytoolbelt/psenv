package cmd

import (
	"fmt"
	"github.com/pytoolbelt/psenv/internal/config"
	"github.com/pytoolbelt/psenv/internal/parameterstore"
	"os"
	"sync"
)

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
