/*
Copyright © 2024 Jesse Maitland jesse@pytoolbelt.com
*/
package cmd

import (
	"fmt"
	"github.com/pytoolbelt/psenv/internal/config"
	"github.com/pytoolbelt/psenv/internal/paramstore"
	"github.com/spf13/cobra"
	"os"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete parameters from the AWS Parameter Store",
	Long:  ``,
	Run:   deleteEntryPoint,
}

func deleteEntryPoint(cmd *cobra.Command, args []string) {

	if envName == "all" {
		fmt.Println("Cannot delete all environments at once. Please specify an environment name with the --env flag.")
		os.Exit(1)
	}

	if envName == "" {
		fmt.Println("Must specify an environment name to delete.")
		os.Exit(1)
	}

	projectConfig, err := config.InitAndLoadProjectConfig()
	if err != nil {
		fmt.Println("Error loading project config %s\n", err)
		os.Exit(1)
	}

	if !projectConfig.HasEnvironment(envName) {
		fmt.Println("Environment %s does not exist in the project configuration.\n", envName)
		os.Exit(1)
	}

	ps, err := paramstore.NewParamStore(projectConfig.GetEnvironmentPath(envName))
	if err != nil {
		fmt.Println("Error creating ssm paramstore %s\n", err)
		os.Exit(1)
	}

	remoteParameterDescriptions, err := ps.DescribeParameters()
	if err != nil {
		fmt.Println("Error describing parameters %s\n", err)
		os.Exit(1)
	}

	err = ps.DeleteParameters(remoteParameterDescriptions)
	if err != nil {
		fmt.Println("Error deleting parameters %s\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
