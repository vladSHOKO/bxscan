package cmd

import (
	"bxscan/internal/scanner"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewFullCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "full [path]",
		Short: "Run full bitrix project analyse",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			path, err := resolveProjectPath(args)
			if err != nil {
				return err
			}

			fmt.Printf("Analysing Bitrix project: %s\n", path)

			result, err := scanner.Scan(os.DirFS(path), scanner.ScanFull, path)
			if err != nil {
				return err
			}

			fmt.Println(result)

			return nil
		},
	}
}
