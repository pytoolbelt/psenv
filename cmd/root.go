/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// flags for config command
var envName string
var decryptFlag bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "psenv",
	Short: "A tool for managing secrets in AWS Parameter Store",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&envName, "env", "e", "", "Environment name")
	rootCmd.PersistentFlags().BoolVarP(&decryptFlag, "decrypt", "d", false, "Decrypt secure strings")
}
