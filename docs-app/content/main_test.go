package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateManifest_IndexSupport(t *testing.T) {
	// Setup temporary source directory
	tmpDir := t.TempDir()
	srcRoot := filepath.Join(tmpDir, "src")
	outRoot := filepath.Join(tmpDir, "out")
	os.MkdirAll(srcRoot, 0755)

	enDir := filepath.Join(srcRoot, "en")
	sectionDir := filepath.Join(enDir, "010_section")
	os.MkdirAll(sectionDir, 0755)

	// Create files
	files := map[string]string{
		filepath.Join(enDir, "_index.md"):      "---\ntitle: \"Home\"\nweight: 10\n---\n# Welcome",
		filepath.Join(sectionDir, "_index.md"): "---\ntitle: \"Introduction\"\nweight: 20\n---\n# Intro",
		filepath.Join(sectionDir, "page.md"):   "---\ntitle: \"The Page\"\n---\n# Page Content",
	}

	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test file %s: %v", path, err)
		}
	}

	manifest, err := generateManifest(srcRoot, outRoot)
	if err != nil {
		t.Fatalf("generateManifest failed: %v", err)
	}

	en, ok := manifest["en"]
	if !ok {
		t.Fatal("manifest missing 'en' root")
	}

	// 1. Verify 'en' root node has properties from en/_index.md
	if en.Title != "Home" {
		t.Errorf("expected en title 'Home', got %q", en.Title)
	}
	if en.Weight != 10 {
		t.Errorf("expected en weight 10, got %d", en.Weight)
	}
	if en.Path != "_index.html" {
		t.Errorf("expected en path '_index.html', got %q", en.Path)
	}

	// 2. Verify _index.md is NOT a child of 'en'
	for _, child := range en.Children {
		if child.Title == "_index.md" || child.Path == "_index.html" {
			t.Errorf("_index.md should not be in children list")
		}
	}

	var section *Document
	for _, child := range en.Children {
		if child.Name == "010_section" {
			section = child
			break
		}
	}

	if section == nil {
		t.Fatal("manifest missing '010_section' node")
	}

	// 3. Verify 'section' node has properties from en/010_section/_index.md
	if section.Title != "Introduction" {
		t.Errorf("expected section title 'Introduction', got %q", section.Title)
	}
	if section.Weight != 20 {
		t.Errorf("expected section weight 20, got %d", section.Weight)
	}
	if section.Path != "010_section/_index.html" {
		t.Errorf("expected section path '010_section/_index.html', got %q", section.Path)
	}

	// 4. Verify 'section' children
	if len(section.Children) != 1 {
		t.Errorf("expected 1 child for section, got %d", len(section.Children))
	}
	if section.Children[0].Title != "The Page" {
		t.Errorf("expected child title 'The Page', got %q", section.Children[0].Title)
	}
}
