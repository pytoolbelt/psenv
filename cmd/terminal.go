/*
Copyright © 2024 Jesse Maitland jesse@pytoolbelt.com
*/
package cmd

import (
	"fmt"
	"github.com/pytoolbelt/psenv/internal/terminal"
	"github.com/pytoolbelt/psenv/internal/utils"
	"os"
	"sync"

	"github.com/pytoolbelt/psenv/internal/config"
	"github.com/spf13/cobra"
)

// terminalEntryPoint is the entry point for the terminal command
func terminalEntryPoint(cmd *cobra.Command, args []string) {

	var envChan = make(chan string, 10)
	var paramsChan = make(chan map[string]string, 25)
	var errorChan = make(chan error, 10)
	var wg sync.WaitGroup
	var numberOfWorkers int = 1
	var projectConfig *config.ProjectConfig

	validateEnvName()
	validateCommand(cmd, args)

	fmt.Printf("Starting terminal session for environment %s\n", terminalEnvName)
	fmt.Println("Type 'exit' to exit the terminal session.")
	fmt.Println("")

	// load the project config
	projectConfig, err := config.LoadProjectConfig()
	if err != nil {
		fmt.Printf("Error loading project config %s\n", err)
		os.Exit(1)
	}

	// check if the environment exists in the project config
	if !projectConfig.HasEnvironment(terminalEnvName) {
		fmt.Printf("Environment %s does not exist in the project configuration.\n", terminalEnvName)
		os.Exit(1)
	}

	// load up the environments to fetch. we always fetch base, plus whatever is specified
	envChan <- "base"
	envChan <- terminalEnvName

	// create 2 workers to fetch the parameters
	for i := 0; i < numberOfWorkers; i++ {
		wg.Add(1)
		go mainGetWorker(envChan, errorChan, paramsChan, &wg, projectConfig, NoDecryptFlag)
	}

	// close the envChan channel
	close(envChan)

	// wait for the workers to finish
	wg.Wait()

	// close the paramsChan channel and errorChan channel
	close(paramsChan)
	close(errorChan)

	// check for errors
	for err := range errorChan {
		fmt.Printf("error getting parameters %s\n", err)
		os.Exit(1)
	}

	// print the parameters
	paramsToConvert := make(map[string]string)

	for params := range paramsChan {
		for k, v := range params {
			paramsToConvert[k] = v
		}
	}

	// convert the parameters to environment variables
	envVars := utils.ConvertParamsToEnvVars(paramsToConvert)

	term, err := terminal.NewSubShell(envVars, args...)
	if err != nil {
		fmt.Printf("error creating terminal %s\n", err)
		os.Exit(1)
	}

	err = term.Run()
	if err != nil {
		fmt.Printf("error starting terminal %s\n", err)
		os.Exit(1)
	}
}

func validateEnvName() {
	if terminalEnvName == "" {
		fmt.Println("Please specify an environment name with the --env flag to start a terminal session.")
		os.Exit(1)
	}
}

func validateCommand(cmd *cobra.Command, args []string) {
	if len(args) == 0 && cmd.Use == "exec" {
		fmt.Println("Please specify a command to execute with the 'exec' command using --.")
		os.Exit(1)
	}

	if len(args) > 0 && cmd.Use == "terminal" {
		fmt.Println("The 'terminal' command does not accept arguments after the -- .")
		os.Exit(1)
	}
}

// terminalCmd represents the terminal command
var terminalCmd = &cobra.Command{
	Use:   "terminal",
	Short: "Start a terminal session with a given environment",
	Long:  ``,
	Run:   terminalEntryPoint,
}

// execCmd represents the terminal command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "execute a command for ",
	Long:  ``,
	Run:   terminalEntryPoint,
}

var NoDecryptFlag bool = false
var terminalEnvName string

func init() {
	rootCmd.AddCommand(terminalCmd)
	terminalCmd.Flags().BoolVar(&NoDecryptFlag, "no-decrypt", true, "Do not decrypt secure string parameters")
	terminalCmd.Flags().StringVarP(&terminalEnvName, "env", "e", "", "The environment to start a terminal session with")

	rootCmd.AddCommand(execCmd)
	execCmd.Flags().BoolVar(&NoDecryptFlag, "no-decrypt", true, "Do not decrypt secure string parameters")
	execCmd.Flags().StringVarP(&terminalEnvName, "env", "e", "", "The environment to start a terminal session with")
}
