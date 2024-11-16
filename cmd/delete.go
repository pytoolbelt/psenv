// /*
// Copyright © 2024 Jesse Maitland jesse@pytoolbelt.com
// */
package cmd

import (
	"fmt"
	"github.com/pytoolbelt/psenv/internal/parameterstore"
	"github.com/spf13/cobra"
	"os"
)

import (
	"github.com/pytoolbelt/psenv/internal/config"
)

func deleteEntryPoint(cmd *cobra.Command, args []string) {

	if deleteEnvName == "" {
		fmt.Printf("Must specify an environment name to delete.\n")
		os.Exit(1)
	}

	projectConfig, err := config.LoadProjectConfig()
	if err != nil {
		fmt.Printf("error loading project config %s\n", err)
		os.Exit(1)
	}

	if !projectConfig.HasEnvironment(deleteEnvName) {
		fmt.Printf("environment %s does not exist in the project configuration.\n", deleteEnvName)
		os.Exit(1)
	}

	ps, err := parameterstore.New()
	if err != nil {
		fmt.Printf("error creating ssm paramstore %s\n", err)
		os.Exit(1)
	}

	remoteParameterDescriptions, err := ps.DescribeParameters(projectConfig.GetEnvironmentPath(deleteEnvName))
	if err != nil {
		fmt.Printf("error describing parameters %s\n", err)
		os.Exit(1)
	}

	if len(remoteParameterDescriptions) == 0 {
		fmt.Printf("No parameters found in the parameter store on path %s\n", projectConfig.GetEnvironmentPath(deleteEnvName))
		os.Exit(0)
	}

	err = ps.DeleteParameters(remoteParameterDescriptions)
	if err != nil {
		fmt.Printf("error deleting parameters %s\n", err)
		os.Exit(1)
	}

	projectConfig.RemoveEnvironment(deleteEnvName)

	err = projectConfig.Save()
	if err != nil {
		fmt.Printf("error saving project config %s\n", err)
		os.Exit(1)
	}

	secretsConfig, err := config.LoadSecretsConfig()
	if err != nil {
		fmt.Println("No secrets config file found. Nothing to update")
		os.Exit(0)
	}

	secretsConfig.ClearEnvironment(deleteEnvName)
	err = secretsConfig.Save()
	if err != nil {
		fmt.Printf("error saving secrets config %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete parameters from the AWS Parameter Store",
	Long:  ``,
	Run:   deleteEntryPoint,
}

var deleteEnvName string

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVarP(&deleteEnvName, "env", "e", "", "The environment to delete parameters from")
}
