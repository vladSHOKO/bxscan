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
	ScanMode    ScanMode
	RootPath    string
	Files       int
	Directories int
	PHPFiles    int

	ComponentsDirExists bool
	ModulesExists       bool
	TemplatesExists     bool

	Components map[string]*analyzer.Component

	Modules map[string]*analyzer.Module

	SecurityFindings []analyzer.SecurityFinding
}

func (r Result) String() string {
	var builder strings.Builder

	switch r.ScanMode {
	case ScanFull:
		r.fullString(&builder)
	case ScanComponents:
		r.componentsString(&builder)
	case ScanModules:
		r.modulesString(&builder)
	case ScanSecurity:
		r.securityString(&builder)
	}

	return builder.String()
}

func (r Result) fullString(builder *strings.Builder) {
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

	builder.WriteString(fmt.Sprintf("Modules: %d\n", len(r.Modules)))
	keys = slices.Collect(maps.Keys(r.Modules))
	slices.Sort(keys)

	for _, key := range keys {
		builder.WriteString(fmt.Sprintf("%s\n", r.Modules[key].String()))
	}

	builder.WriteString(fmt.Sprintf("Security findings: %d\n", len(r.SecurityFindings)))

	for _, finding := range r.SecurityFindings {
		builder.WriteString(fmt.Sprintf("%s\n", finding.String()))
	}
}

func (r Result) componentsString(builder *strings.Builder) {
	builder.WriteString(fmt.Sprintf("Components: %d\n", len(r.Components)))
	keys := slices.Collect(maps.Keys(r.Components))
	slices.Sort(keys)

	for _, key := range keys {
		builder.WriteString(fmt.Sprintf("%s\n", r.Components[key].String()))
	}
}

func (r Result) modulesString(builder *strings.Builder) {
	builder.WriteString(fmt.Sprintf("Modules: %d\n", len(r.Modules)))
	keys := slices.Collect(maps.Keys(r.Modules))
	slices.Sort(keys)

	for _, key := range keys {
		builder.WriteString(fmt.Sprintf("%s\n", r.Modules[key].String()))
	}
}

func (r Result) securityString(builder *strings.Builder) {
	builder.WriteString(fmt.Sprintf("Security findings: %d\n", len(r.SecurityFindings)))

	for _, finding := range r.SecurityFindings {
		builder.WriteString(fmt.Sprintf("%s\n", finding.String()))
	}
}

var targetList = map[string]struct{}{
	"components":    {},
	"templates":     {},
	"modules":       {},
	"php_interface": {},
}

func Scan(path string, mod ScanMode) (*Result, error) {
	result := Result{
		ScanMode:         mod,
		RootPath:         path,
		Components:       map[string]*analyzer.Component{},
		Modules:          map[string]*analyzer.Module{},
		SecurityFindings: []analyzer.SecurityFinding{},
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

		switch mod {
		case ScanFull:
			err = runFullScan(&result, d, section, relPath, currentPath)
			if err != nil {
				return err
			}
		case ScanComponents:
			err = runComponentsScan(&result, d, section, relPath, currentPath)
			if err != nil {
				return err
			}
		case ScanModules:
			err = runModulesScan(&result, d, section, relPath)
			if err != nil {
				return err
			}
		case ScanSecurity:
			err = runSecurityScan(&result, d, section, currentPath, relPath)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func runFullScan(result *Result, d fs.DirEntry, section, relPath, currentPath string) error {
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
		if err := analyzer.AnalyzeComponent(currentPath, relPath, d, result.Components); err != nil {
			return err
		}
	case "templates":
		result.TemplatesExists = true
	case "modules":
		result.ModulesExists = true
		if err := analyzer.AnalyzeModule(relPath, d, result.Modules); err != nil {
			return err
		}
	}

	if !d.IsDir() && filepath.Ext(relPath) == ".php" {
		securityAnalyzeResults, err := analyzer.AnalyzeSecurity(currentPath, relPath)
		if err != nil {
			return err
		}
		for _, securityResult := range securityAnalyzeResults {
			result.SecurityFindings = append(result.SecurityFindings, securityResult)
		}
	}

	return nil
}

func runComponentsScan(result *Result, d fs.DirEntry, section, relPath, currentPath string) error {
	if section != "components" {
		if d.IsDir() {
			return filepath.SkipDir
		}
		return nil
	}

	if err := analyzer.AnalyzeComponent(currentPath, relPath, d, result.Components); err != nil {
		return err
	}

	return nil
}

func runModulesScan(result *Result, d fs.DirEntry, section, relPath string) error {
	if section != "modules" {
		if d.IsDir() {
			return filepath.SkipDir
		}
		return nil
	}

	if err := analyzer.AnalyzeModule(relPath, d, result.Modules); err != nil {
		return err
	}

	return nil
}

func runSecurityScan(result *Result, d fs.DirEntry, section, currentPath, relPath string) error {
	if _, ok := targetList[section]; !ok {
		if d.IsDir() {
			return filepath.SkipDir
		}

		return nil
	}

	if !d.IsDir() && filepath.Ext(relPath) == ".php" {
		securityAnalyzeResults, err := analyzer.AnalyzeSecurity(currentPath, relPath)
		if err != nil {
			return err
		}
		for _, securityResult := range securityAnalyzeResults {
			result.SecurityFindings = append(result.SecurityFindings, securityResult)
		}
	}

	return nil
}
