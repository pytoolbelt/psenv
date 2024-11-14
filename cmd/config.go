/*
Copyright © 2024 Jesse Maitland jesse@pytoolbelt.com
*/
package cmd

import (
	"fmt"
	"github.com/pytoolbelt/psenv/internal/config"
	"github.com/spf13/cobra"
	"os"
)

var newProjectFlag bool
var newSecretsFlag bool
var printFlag bool

func configEntrypoint(cmd *cobra.Command, args []string) {
	validateConfigFlags()

	if printFlag {
		projectConfig, err := config.LoadProjectConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		projectConfig.PrintTable()
		os.Exit(0)
	}

	if newProjectFlag {
		_, err := config.CreateNewProjectConfigFile()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Created new project config file psenv-project.yml")
		os.Exit(0)
	}

	//if newSecretsFlag {
	//	_, err := config.CreateNewSecretsConfigFile()
	//	if err != nil {
	//		fmt.Println(err)
	//		os.Exit(1)
	//	}
	//	fmt.Println("Created new secrets file psenv-secrets.yml")
	//	os.Exit(0)
	//}
}

func validateConfigFlags() {
	if newProjectFlag && printFlag || newSecretsFlag && printFlag || newProjectFlag && newSecretsFlag {
		fmt.Println("Cannot use both --new and --print flags")
		os.Exit(1)
	}
	if !newProjectFlag && !printFlag && !newSecretsFlag {
		fmt.Println("Must use either --new-project, --new-secrets or --print flag")
		os.Exit(1)
	}

}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Do things with config",
	Long:  ``,
	Run:   configEntrypoint,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().BoolVarP(&newProjectFlag, "new-project", "", false, "Create a new project config file")
	configCmd.Flags().BoolVarP(&newSecretsFlag, "new-secrets", "", false, "Create a new secrets file")
	configCmd.Flags().BoolVarP(&printFlag, "print", "p", false, "Print the project config")
}
