package cmd

import (
	"bxscan/internal/scanner"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewComponentsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "components [path]",
		Short: "Run analysis of components",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			path, err := resolveProjectPath(args)
			if err != nil {
				return err
			}

			fmt.Println("Analysing components...")

			result, err := scanner.Scan(os.DirFS(path), scanner.ScanComponents, path)
			if err != nil {
				return err
			}

			fmt.Println(result)

			return nil
		},
	}
}
