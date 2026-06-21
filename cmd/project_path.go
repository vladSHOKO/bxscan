package cmd

import (
	"fmt"
	"os"
)

func resolveProjectPath(args []string) (string, error) {
	path := "./local"

	if len(args) > 0 {
		path = args[0]
	}

	dir, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	if !dir.IsDir() {
		return "", fmt.Errorf("%s is not a directory", path)
	}

	return path, nil
}
