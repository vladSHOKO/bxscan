package analyzer

import (
	"strings"
	"testing"
	"testing/fstest"
)

func TestContainsSQL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "simple query",
			input:    "select id, name, something from table",
			expected: true,
		},
		{
			name:     "many strings select query",
			input:    "select id, name, something \n from table",
			expected: true,
		},
		{
			name:     "query in upper case",
			input:    "SELECT id, name, something FROM table",
			expected: true,
		},
		{
			name:     "bitrix execute simple query",
			input:    "$DB->Query('select id, name, something from table')",
			expected: true,
		},
		{
			name:     "empty content",
			input:    "",
			expected: false,
		},
		{
			name:     "simple php content",
			input:    "<?php\n" + "echo 'Hello World!'",
			expected: false,
		},
		{
			name:     "hmtl select tag",
			input:    "<select>one</select>",
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ContainsSQL(test.input)
			if result != test.expected {
				t.Errorf("expected: %v, got: %v", test.expected, result)
			}
		})
	}
}

func TestAnalyzeTemplate(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		wantLines    int
		wantTooLarge bool
	}{
		{
			name:         "empty template",
			content:      "",
			wantLines:    0,
			wantTooLarge: false,
		},
		{
			name:         "template at line limit",
			content:      strings.Repeat("some line \n", 500),
			wantLines:    500,
			wantTooLarge: false,
		},
		{
			name:         "template over line limit",
			content:      strings.Repeat("some line \n", 501),
			wantLines:    501,
			wantTooLarge: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fsys := fstest.MapFS{
				"components/acme/news/templates/.default/template.php": {
					Data: []byte(test.content),
				},
			}

			parts := []string{
				"components",
				"acme",
				"news",
				"templates",
				".default",
				"template.php",
			}

			result, err := AnalyzeTemplate(
				"components/acme/news/templates/.default/template.php",
				parts,
				fsys,
			)

			if err != nil {
				t.Fatalf("AnalyzeTemplate error: %v", err)
			}

			if result.Lines != test.wantLines {
				t.Errorf("lines want: %d, got: %d", test.wantLines, result.Lines)
			}

			if result.TooLarge != test.wantTooLarge {
				t.Errorf("too large: want: %v, got: %v", test.wantTooLarge, result.TooLarge)
			}
		})
	}
}
