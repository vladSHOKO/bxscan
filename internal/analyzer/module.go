package analyzer

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

type Module struct {
	ID              string
	HasInstallIndex bool
	HasLib          bool
	HasInclude      bool
	HasOptions      bool
}

func (m *Module) String() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("\tModule ID: %s\n", m.ID))
	builder.WriteString(fmt.Sprintf("\tHas install/index.php: %v\n", m.HasInstallIndex))
	builder.WriteString(fmt.Sprintf("\tHas lib dir: %v\n", m.HasLib))
	builder.WriteString(fmt.Sprintf("\tHas include.php: %v\n", m.HasInclude))
	builder.WriteString(fmt.Sprintf("\tHas options.php: %v\n", m.HasOptions))

	return builder.String()
}

func AnalyzeModule(relPath string, d fs.DirEntry, modulesMap map[string]*Module) error {
	parts := strings.Split(relPath, string(filepath.Separator))

	if len(parts) < 2 {
		return nil
	}

	module, ok := modulesMap[parts[1]]
	if !ok {
		if len(parts) == 2 && d.IsDir() {
			module = &Module{ID: parts[1]}
			modulesMap[parts[1]] = module
		} else {
			return nil
		}
	}

	if len(parts) < 3 {
		return nil
	}

	if parts[2] == "include.php" && !d.IsDir() {
		module.HasInclude = true
	}

	if parts[2] == "lib" && d.IsDir() {
		module.HasLib = true
	}

	if parts[2] == "options.php" && !d.IsDir() {
		module.HasOptions = true
	}

	if len(parts) == 4 && parts[2] == "install" && parts[3] == "index.php" && !d.IsDir() {
		module.HasInstallIndex = true
	}

	return nil
}
