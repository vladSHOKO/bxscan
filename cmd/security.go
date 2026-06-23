package cmd

import (
	"bxscan/internal/scanner"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewSecurityCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "security [path]",
		Short: "Run analysis of security",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			path, err := resolveProjectPath(args)
			if err != nil {
				return err
			}

			fmt.Println("Analysing security...")

			result, err := scanner.Scan(os.DirFS(path), scanner.ScanSecurity, path)
			if err != nil {
				return err
			}

			fmt.Println(result)

			return nil
		},
	}
}
