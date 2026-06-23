package analyzer

import (
	"bufio"
	"fmt"
	"io/fs"
	"regexp"
	"strings"
)

type SecurityFinding struct {
	RelPath string
	Line    int
	Kind    string
	Snippet string
}

type SecurityPattern struct {
	Kind    string
	Pattern *regexp.Regexp
}

func (sf *SecurityFinding) String() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%s:%d\n", sf.RelPath, sf.Line))
	builder.WriteString(fmt.Sprintf("\t%s\n", sf.Kind))
	builder.WriteString(fmt.Sprintf("\t%s\n", sf.Snippet))

	return builder.String()
}

var securityPatterns = []SecurityPattern{
	{
		Kind:    "eval usage",
		Pattern: regexp.MustCompile(`\beval\s*\(`),
	},
	{
		Kind:    "unserialize usage",
		Pattern: regexp.MustCompile(`\bunserialize\s*\(`),
	},
	{
		Kind:    "$_REQUEST usage",
		Pattern: regexp.MustCompile(`\$_REQUEST\b`),
	},
	{
		Kind:    "$_GET usage",
		Pattern: regexp.MustCompile(`\$_GET\b`),
	},
	{
		Kind:    "$_POST usage",
		Pattern: regexp.MustCompile(`\$_POST\b`),
	},
}

func AnalyzeSecurity(fsys fs.FS, relPath string) ([]SecurityFinding, error) {
	var securityFindings []SecurityFinding

	file, err := fsys.Open(relPath)
	if err != nil {
		return securityFindings, fmt.Errorf("failed to open file %s. %w", relPath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		for _, securityPattern := range securityPatterns {
			if securityPattern.Pattern.MatchString(line) {
				securityFindings = append(securityFindings, SecurityFinding{
					RelPath: relPath,
					Line:    lineNumber,
					Kind:    securityPattern.Kind,
					Snippet: makeSnippet(line, *securityPattern.Pattern),
				})
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return securityFindings, fmt.Errorf("failed to scan file %s. %w", relPath, err)
	}

	return securityFindings, nil
}

func makeSnippet(line string, regex regexp.Regexp) string {
	endSize := 20

	loc := regex.FindStringIndex(line)
	if loc == nil {
		return line
	}

	start, end := loc[0], loc[1]

	snippetEnd := end + endSize
	if snippetEnd > len(line) {
		snippetEnd = len(line)
	}

	snippet := line[start:snippetEnd]

	return snippet
}
