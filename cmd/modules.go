package cmd

import (
	"bxscan/internal/scanner"
	"fmt"

	"github.com/spf13/cobra"
)

func NewModulesCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "modules [path]",
		Short: "Run analysis of modules",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			path, err := resolveProjectPath(args)
			if err != nil {
				return err
			}

			fmt.Println("Analysing modules...")

			result, err := scanner.Scan(path, scanner.ScanModules)
			if err != nil {
				return err
			}

			fmt.Println(result)

			return nil
		},
	}
}
