/*
Copyright © 2024 Jesse Maitland jesse@pytoolbelt.com
*/
package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"sync"

	"github.com/pytoolbelt/psenv/internal/config"
	"github.com/pytoolbelt/psenv/internal/paramstore"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get parameters from the AWS Parameter Store",
	Long:  ``,
	Run:   getEntryPoint,
}

func getEntryPoint(cmd *cobra.Command, args []string) {
	var wg sync.WaitGroup

	projectConfig, err := config.InitAndLoadProjectConfig()

	if err != nil {
		fmt.Printf("Error loading project config %s\n", err)
		os.Exit(1)
	}

	basePath := projectConfig.GetBasePath()

	ps, err := paramstore.NewParamStore(basePath)
	if err != nil {
		fmt.Printf("Error creating ssm paramstore %s\n", err)
		os.Exit(1)
	}

	remoteParameterDescriptions, err := ps.DescribeParameters()
	if err != nil {
		fmt.Printf("Error describing parameters %s\n", err)
		os.Exit(1)
	}

	remoteParameterDescriptions = config.RemoveVarNameFromPaths(remoteParameterDescriptions)
	remoteParameterDescriptions = paramstore.SplitAndDeduplicatePaths(remoteParameterDescriptions)

	var environments []*config.Environment

	for _, p := range remoteParameterDescriptions {
		environments = append(environments, config.NewEnvironmentFromPath(p))
	}

	envChan := make(chan config.Environment, len(environments))
	resultChan := make(chan config.Environment, len(environments))

	// put env on the channel
	for _, env := range environments {
		envChan <- *env
	}
	close(envChan)

	// Start worker Go routines
	for i := 0; i <= 3; i++ { // Number of workers
		wg.Add(1)
		go getWorker(envChan, resultChan, &wg)
	}

	wg.Wait()
	close(resultChan)

	var cfg *config.Config

	cfg, err = config.InitAndLoad()
	if err != nil {
		// if the config does not exist, create a new template file
		cfg, err = config.CreateNewSecretsConfigFile()
		if err != nil {
			fmt.Printf("Error creating new config file %s\n", err)
			os.Exit(1)
		}
	}

	cfg.ClearEnvironments()
	for env := range resultChan {
		cfg.SetEnvironment(&env)
	}

	err = viper.WriteConfig()
	if err != nil {
		fmt.Printf("Error saving config file %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func getWorker(envChan <-chan config.Environment, resultChan chan<- config.Environment, wg *sync.WaitGroup) {
	defer wg.Done()

	for env := range envChan {
		psPath := env.GetParamStorePath()
		ps, err := paramstore.NewParamStore(psPath)

		if err != nil {
			fmt.Printf("Error creating ssm paramstore %s\n", err)
			continue
		}

		newParams, err := ps.GetParameters(decryptFlag)
		if err != nil {
			fmt.Printf("Error getting parameters %s\n", err)
			continue
		}

		env.Params = newParams
		resultChan <- env
	}
}

func entrypoint(cmd *cobra.Command, args []string) {
	cfg, err := config.InitAndLoad()
	HandelError(err)

	basePath := cfg.GetBasePath()

	ps, err := paramstore.NewParamStore(basePath)
	HandelError(err)

	params, err := ps.DescribeParameters()
	HandelError(err)

	if len(params) == 0 {
		for k := range cfg.Environments {
			e, er := cfg.GetEnvironment(k)
			HandelError(er)
			e.Params = nil
			cfg.SetEnvironment(e)
		}
		err = cfg.Save()
		HandelError(err)
		os.Exit(0)
	}

	params = paramstore.SplitAndDeduplicatePaths(params)

	for _, p := range params {
		ps, err = paramstore.NewParamStore(p)
		HandelError(err)

		newParams, err := ps.GetParameters(decryptFlag)
		HandelError(err)

		newEnvName, err := paramstore.GetEnvNameFromSSMPath(p)
		HandelError(err)

		env := &config.Environment{
			Name:    newEnvName,
			Project: cfg.Project,
			Prefix:  cfg.Prefix,
			Params:  newParams,
		}

		cfg.SetEnvironment(env)
	}
	err = cfg.Save()
	HandelError(err)
	os.Exit(0)
}

func init() {
	rootCmd.AddCommand(getCmd)
}
