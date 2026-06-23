package cmd

import (
	"bxscan/internal/scanner"
	"fmt"
	"os"

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

			result, err := scanner.Scan(os.DirFS(path), scanner.ScanModules, path)
			if err != nil {
				return err
			}

			fmt.Println(result)

			return nil
		},
	}
}
