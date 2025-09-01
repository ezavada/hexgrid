package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadYAMLConfig(t *testing.T) {
	// Create a temporary YAML file for testing
	testYAML := `default: "#FFFFFF"
items:
  - name: "Test Item"
    percentage: 50.0
    style: "fill"
    color: "#FF0000"
  - name: "Test Dot"
    percentage: 50.0
    style: "dot"
    color: "#00FF00"`

	tmpFile, err := os.CreateTemp("", "test_config_*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(testYAML)
	if err != nil {
		t.Fatalf("Failed to write test YAML: %v", err)
	}
	tmpFile.Close()

	// Test loading the configuration
	config, err := LoadYAMLConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load YAML config: %v", err)
	}

	// Verify the configuration
	if len(config.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(config.Items))
	}

	if config.Items[0].Name != "Test Item" {
		t.Errorf("Expected first item name to be 'Test Item', got '%s'", config.Items[0].Name)
	}

	if config.Items[0].Percentage != 50.0 {
		t.Errorf("Expected first item percentage to be 50.0, got %f", config.Items[0].Percentage)
	}

	if config.Items[0].Style != "fill" {
		t.Errorf("Expected first item style to be 'fill', got '%s'", config.Items[0].Style)
	}

	if config.Items[1].Style != "dot" {
		t.Errorf("Expected second item style to be 'dot', got '%s'", config.Items[1].Style)
	}
}

func TestYAMLFileDiscovery(t *testing.T) {
	// Test that we can read YAML files from the grid-specs directory
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	gridSpecsDir := filepath.Join(currentDir, "grid-specs")

	// Read all YAML files from the grid-specs directory
	files, err := os.ReadDir(gridSpecsDir)
	if err != nil {
		t.Fatalf("Failed to read grid-specs directory: %v", err)
	}

	var yamlFiles []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".yaml" {
			yamlFiles = append(yamlFiles, file.Name())
		}
	}

	// We should have at least the sample files
	if len(yamlFiles) < 3 {
		t.Errorf("Expected at least 3 YAML files, got %d", len(yamlFiles))
	}

	// Check for expected files
	expectedFiles := []string{"sample_config.yaml", "fantasy-world.yaml", "desert-world.yaml", "hellstella-space.yaml"}
	for _, expected := range expectedFiles {
		found := false
		for _, file := range yamlFiles {
			if file == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected file %s not found in grid-specs directory", expected)
		}
	}
}
