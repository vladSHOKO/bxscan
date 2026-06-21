package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "bxscan",
		Short: "Static analyser for bitrix projects",
		Long:  "bxscan is CLI tool for analysing Bitrix projects.",
	}

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	return rootCmd
}

func Execute() {
	rootCmd := NewRootCommand()
	rootCmd.AddCommand(NewFullCommand())
	rootCmd.AddCommand(NewComponentsCommand())
	rootCmd.AddCommand(NewModulesCommand())
	rootCmd.AddCommand(NewSecurityCommand())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
