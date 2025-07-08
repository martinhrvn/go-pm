package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindConfigFile(t *testing.T) {
	tests := []struct {
		name           string
		setupDirs      []string
		setupConfigs   map[string]string
		startDir       string
		expectedConfig string
		wantErr        bool
	}{
		{
			name:      "config in current directory",
			setupDirs: []string{"project"},
			setupConfigs: map[string]string{
				"project/.gopmrc": `locations:
  - location: "src"
    commands: ["build"]`,
			},
			startDir:       "project",
			expectedConfig: "project/.gopmrc",
			wantErr:        false,
		},
		{
			name:      "config in parent directory",
			setupDirs: []string{"project", "project/subdir"},
			setupConfigs: map[string]string{
				"project/.gopmrc": `locations:
  - location: "src"
    commands: ["build"]`,
			},
			startDir:       "project/subdir",
			expectedConfig: "project/.gopmrc",
			wantErr:        false,
		},
		{
			name:      "config in grandparent directory",
			setupDirs: []string{"project", "project/subdir", "project/subdir/nested"},
			setupConfigs: map[string]string{
				"project/.gopmrc": `locations:
  - location: "src"
    commands: ["build"]`,
			},
			startDir:       "project/subdir/nested",
			expectedConfig: "project/.gopmrc",
			wantErr:        false,
		},
		{
			name:      "config in closer directory takes precedence",
			setupDirs: []string{"project", "project/subdir"},
			setupConfigs: map[string]string{
				"project/.gopmrc": `locations:
  - location: "root"
    commands: ["build"]`,
				"project/subdir/.gopmrc": `locations:
  - location: "subdir"
    commands: ["test"]`,
			},
			startDir:       "project/subdir",
			expectedConfig: "project/subdir/.gopmrc",
			wantErr:        false,
		},
		{
			name:      "no config file found",
			setupDirs: []string{"project", "project/subdir"},
			startDir:  "project/subdir",
			wantErr:   true,
		},
		{
			name:      "config in root directory",
			setupDirs: []string{"project", "project/subdir", "project/subdir/nested", "project/subdir/nested/deep"},
			setupConfigs: map[string]string{
				"project/.gopmrc": `locations:
  - location: "root"
    commands: ["build"]`,
			},
			startDir:       "project/subdir/nested/deep",
			expectedConfig: "project/.gopmrc",
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "gopm-discovery-test")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			// Create directory structure
			for _, dir := range tt.setupDirs {
				err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755)
				if err != nil {
					t.Fatalf("Failed to create dir %s: %v", dir, err)
				}
			}

			// Create config files
			for configPath, content := range tt.setupConfigs {
				fullPath := filepath.Join(tmpDir, configPath)
				err := os.WriteFile(fullPath, []byte(content), 0644)
				if err != nil {
					t.Fatalf("Failed to create config file %s: %v", configPath, err)
				}
			}

			// Change to start directory
			startPath := filepath.Join(tmpDir, tt.startDir)
			oldWd, err := os.Getwd()
			if err != nil {
				t.Fatalf("Failed to get working directory: %v", err)
			}
			defer os.Chdir(oldWd)

			err = os.Chdir(startPath)
			if err != nil {
				t.Fatalf("Failed to change to start directory: %v", err)
			}

			// Test FindConfigFile
			configPath, err := FindConfigFile()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("FindConfigFile() error = %v", err)
				return
			}

			// Convert expected path to absolute for comparison
			expectedPath := filepath.Join(tmpDir, tt.expectedConfig)
			if configPath != expectedPath {
				t.Errorf("FindConfigFile() = %q, expected %q", configPath, expectedPath)
			}
		})
	}
}

func TestLoadConfigFromDiscovery(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gopm-discovery-integration-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create directory structure
	dirs := []string{"project", "project/subdir", "project/subdir/nested"}
	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755)
		if err != nil {
			t.Fatalf("Failed to create dir %s: %v", dir, err)
		}
	}

	// Create config in parent directory
	configContent := `locations:
  - name: "frontend"
    location: "packages/frontend"
    type: "npm"
    commands: ["start", "build"]`

	configPath := filepath.Join(tmpDir, "project/.gopmrc")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Create test directories for glob expansion
	testDirs := []string{"project/packages/frontend", "project/packages/backend"}
	for _, dir := range testDirs {
		err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755)
		if err != nil {
			t.Fatalf("Failed to create test dir %s: %v", dir, err)
		}
	}

	// Change to nested directory
	nestedPath := filepath.Join(tmpDir, "project/subdir/nested")
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(oldWd)

	err = os.Chdir(nestedPath)
	if err != nil {
		t.Fatalf("Failed to change to nested directory: %v", err)
	}

	// Test LoadConfigFromDiscovery
	config, err := LoadConfigFromDiscovery()
	if err != nil {
		t.Fatalf("LoadConfigFromDiscovery() error = %v", err)
	}

	// Verify config was loaded correctly
	if len(config.Locations) != 1 {
		t.Errorf("Expected 1 location, got %d", len(config.Locations))
	}

	if config.Locations[0].Name != "frontend" {
		t.Errorf("Expected location name 'frontend', got %q", config.Locations[0].Name)
	}

	if config.Locations[0].Location != "packages/frontend" {
		t.Errorf("Expected location path 'packages/frontend', got %q", config.Locations[0].Location)
	}
}