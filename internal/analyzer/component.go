package analyzer

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

type Component struct {
	Vendor   string
	Name     string
	FullName string

	ClassExists       bool
	ComponentExists   bool
	DescriptionExists bool
	ParametersExists  bool
}

func AnalyzeComponent(relPath string, d fs.DirEntry, componentsMap map[string]*Component) {
	parts := strings.Split(relPath, string(filepath.Separator))
	if len(parts) < 3 {
		return
	}

	componentKey := filepath.Join(parts[1], parts[2])

	if len(parts) == 3 && d.IsDir() {
		if _, ok := componentsMap[componentKey]; !ok {
			componentsMap[componentKey] = &Component{
				Vendor:   parts[1],
				Name:     parts[2],
				FullName: componentKey,
			}
		}

		return
	}

	if len(parts) != 4 || d.IsDir() {
		return
	}

	result, ok := componentsMap[componentKey]
	if !ok {
		return
	}

	switch parts[3] {
	case "class.php":
		result.ClassExists = true
	case ".parameters.php":
		result.ParametersExists = true
	case "component.php":
		result.ComponentExists = true
	case ".description.php":
		result.DescriptionExists = true
	}
}

func (c *Component) String() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%s\n", c.FullName))
	builder.WriteString(fmt.Sprintf("\tclass.php: %t\n", c.ClassExists))
	builder.WriteString(fmt.Sprintf("\tcomponent.php: %t\n", c.ComponentExists))
	builder.WriteString(fmt.Sprintf("\t.description.php: %t\n", c.DescriptionExists))
	builder.WriteString(fmt.Sprintf("\t.parameters.php: %t\n", c.ParametersExists))

	return builder.String()
}
