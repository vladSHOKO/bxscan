package scanner

import (
	"bxscan/internal/analyzer"
	"fmt"
	"io/fs"
	"maps"
	"path/filepath"
	"slices"
	"strings"
)

type Result struct {
	RootPath    string
	Files       int
	Directories int
	PHPFiles    int

	ComponentsDirExists bool
	ModulesExists       bool
	TemplatesExists     bool

	Components map[string]*analyzer.Component
}

func (r Result) String() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Root path: %s\n", r.RootPath))
	builder.WriteString(fmt.Sprintf("Files: %d\n", r.Files))
	builder.WriteString(fmt.Sprintf("Directories: %d\n", r.Directories))
	builder.WriteString(fmt.Sprintf("PHP files: %d\n", r.PHPFiles))
	builder.WriteString(fmt.Sprintf("Components dir exists: %t\n", r.ComponentsDirExists))
	builder.WriteString(fmt.Sprintf("Modules exists: %t\n", r.ModulesExists))
	builder.WriteString(fmt.Sprintf("Templates exists: %t\n", r.TemplatesExists))

	builder.WriteString(fmt.Sprintf("Components: %d\n", len(r.Components)))
	keys := slices.Collect(maps.Keys(r.Components))
	slices.Sort(keys)

	for _, key := range keys {
		builder.WriteString(fmt.Sprintf("%s\n", r.Components[key].String()))
	}

	return builder.String()
}

var targetList = map[string]struct{}{
	"components": {},
	"templates":  {},
	"modules":    {},
}

func Scan(path string) (*Result, error) {
	result := Result{
		RootPath:   path,
		Components: map[string]*analyzer.Component{},
	}
	err := filepath.WalkDir(path, func(currentPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(path, currentPath)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		parts := strings.Split(relPath, string(filepath.Separator))
		section := parts[0]

		if _, ok := targetList[section]; !ok {
			if d.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}

		if d.IsDir() {
			result.Directories++
		} else {
			result.Files++

			if filepath.Ext(relPath) == ".php" {
				result.PHPFiles++
			}
		}

		switch section {
		case "components":
			result.ComponentsDirExists = true
			analyzer.AnalyzeComponent(relPath, d, result.Components)
		case "templates":
			result.TemplatesExists = true
		case "modules":
			result.ModulesExists = true
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &result, nil
}
