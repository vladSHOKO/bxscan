package scanner

import (
	"io/fs"
	"testing"
	"testing/fstest"
)

func TestScanFull(t *testing.T) {
	fsys := fstest.MapFS{
		"components/acme/news/class.php": {
			Data: []byte("<?php"),
		},
		"components/acme/news/component.php": {
			Data: []byte("<?php"),
		},
		"components/acme/news/templates/.default/template.php": {
			Data: []byte("<?php SELECT ID FROM b_user;"),
		},

		"modules/acme.demo/install/index.php": {
			Data: []byte("<?php"),
		},
		"modules/acme.demo/lib": {
			Mode: fs.ModeDir,
		},
		"modules/acme.demo/include.php": {
			Data: []byte("<?php eval($someCode);"),
		},
		"upload/ignored.php": {
			Data: []byte("<?php eval($ignored);"),
		},
	}

	result, err := Scan(fsys, ScanFull, "test-local")
	if err != nil {
		t.Fatalf("Failed to scan full acme module: %v", err)
	}
	if len(result.Components) != 1 {
		t.Fatalf("Expected 1 component, got: %d", len(result.Components))
	}

	component := result.Components["acme/news"]
	if component == nil {
		t.Fatal("component acme/news was not found")
	}
	if !component.ClassExists {
		t.Fatalf("component ClassExists = false, want true")
	}
	if !component.ComponentExists {
		t.Error("component ComponentExists = false, want true")
	}

	template := component.Templates[".default"]
	if template == nil {
		t.Fatal("template .default was not found")
	}
	if !template.SQLExists {
		t.Error("template SQLExists = false, want true")
	}

	module := result.Modules["acme.demo"]
	if module == nil {
		t.Fatal("module acme.demo was not found")
	}
	if !module.HasInstallIndex {
		t.Error("module HasInstallIndex = false, want true")
	}
	if !module.HasLib {
		t.Error("module HasLib = false, want true")
	}
	if !module.HasInclude {
		t.Error("module HasInclude = false, want true")
	}
	if len(result.SecurityFindings) != 1 {
		t.Fatalf("security findings count = %d, want 1", len(result.SecurityFindings))
	}

	finding := result.SecurityFindings[0]
	if finding.RelPath != "modules/acme.demo/include.php" {
		t.Errorf("security finding path = %q, want %q", finding.RelPath, "modules/acme.demo/include.php")
	}
}
