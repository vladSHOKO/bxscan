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

			path := "./local"

			if len(args) > 0 {
				path = args[0]
			}

			directory, err := os.Stat(path)
			if err != nil {
				return err
			}

			isDir := directory.IsDir()
			if !isDir {
				return fmt.Errorf("%s is not a directory", path)
			}

			fmt.Printf("Analysing Bitrix project: %s\n", path)
			fmt.Printf("Directory found: %s\n", directory.Name())

			result, err := scanner.Scan(path)
			if err != nil {
				return err
			}

			fmt.Println(result)

			return nil
		},
	}
}
