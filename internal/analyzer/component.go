package analyzer

import (
	"bufio"
	"fmt"
	"io/fs"
	"maps"
	"path"
	"regexp"
	"slices"
	"sort"
	"strings"
)

const MaxTemplateLines = 500

type Component struct {
	Vendor    string
	Name      string
	FullName  string
	Templates map[string]*Template

	ClassExists       bool
	ComponentExists   bool
	DescriptionExists bool
	ParametersExists  bool
}

type Template struct {
	Name      string
	RelPath   string
	Lines     int
	TooLarge  bool
	SQLExists bool
}

func AnalyzeComponent(fsys fs.FS, relPath string, d fs.DirEntry, componentsMap map[string]*Component) error {
	parts := strings.Split(relPath, "/")
	if len(parts) < 3 {
		return nil
	}

	componentKey := path.Join(parts[1], parts[2])

	if len(parts) == 3 && d.IsDir() {
		if _, ok := componentsMap[componentKey]; !ok {
			componentsMap[componentKey] = &Component{
				Vendor:   parts[1],
				Name:     parts[2],
				FullName: componentKey,

				Templates: map[string]*Template{},
			}
		}

		return nil
	}

	if len(parts) == 4 && !d.IsDir() {
		result, ok := componentsMap[componentKey]
		if !ok {
			return nil
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

	if len(parts) == 6 && parts[3] == "templates" && parts[5] == "template.php" && !d.IsDir() {
		templateKey := parts[4]
		analyzeResult, err := AnalyzeTemplate(relPath, parts, fsys)
		if err != nil {
			return err
		}

		componentsMap[componentKey].Templates[templateKey] = analyzeResult
	}

	return nil
}

func (c *Component) String() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%s\n", c.FullName))
	builder.WriteString(fmt.Sprintf("\tclass.php: %t\n", c.ClassExists))
	builder.WriteString(fmt.Sprintf("\tcomponent.php: %t\n", c.ComponentExists))
	builder.WriteString(fmt.Sprintf("\t.description.php: %t\n", c.DescriptionExists))
	builder.WriteString(fmt.Sprintf("\t.parameters.php: %t\n", c.ParametersExists))

	builder.WriteString(fmt.Sprintf("\tTemplates: %d\n", len(c.Templates)))

	keys := slices.Collect(maps.Keys(c.Templates))
	sort.Strings(keys)

	for _, key := range keys {
		builder.WriteString(fmt.Sprintf("\t%s\n", c.Templates[key].String()))
	}

	return builder.String()
}

func AnalyzeTemplate(relPath string, parts []string, fsys fs.FS) (*Template, error) {
	var builder strings.Builder
	result := &Template{
		Name:    parts[4],
		RelPath: relPath,
	}

	file, err := fsys.Open(relPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	for scanner.Scan() {
		result.Lines++
		builder.WriteString(scanner.Text())
		builder.WriteByte('\n')
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	result.SQLExists = ContainsSQL(builder.String())

	if result.Lines > MaxTemplateLines {
		result.TooLarge = true
	}

	return result, nil
}

var sqlPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?is)\bselect\b.{1,1000}\bfrom\b`),
	regexp.MustCompile(`(?is)\binsert\s+into\b`),
	regexp.MustCompile(`(?is)\bupdate\b.{1,500}\bset\b`),
	regexp.MustCompile(`(?is)\bdelete\s+from\b`),
	regexp.MustCompile(`(?i)\$DB\s*->\s*(Query|QueryBind)\s*\(`),
	regexp.MustCompile(`(?i)->\s*query\s*\(`),
}

func ContainsSQL(content string) bool {
	for _, pattern := range sqlPatterns {
		if pattern.MatchString(content) {
			return true
		}
	}

	return false
}

func (t *Template) String() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("\t%s\n", t.Name))
	builder.WriteString(fmt.Sprintf("\t\t\tlines: %d\n", t.Lines))
	builder.WriteString(fmt.Sprintf("\t\t\ttoo large: %t\n", t.TooLarge))
	builder.WriteString(fmt.Sprintf("\t\t\tsql exists: %t\n", t.SQLExists))

	return builder.String()
}
