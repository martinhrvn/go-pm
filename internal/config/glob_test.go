package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandGlobPatterns(t *testing.T) {
	tests := []struct {
		name        string
		setupDirs   []string
		setupFiles  []string
		locations   []Location
		expected    []Location
		wantErr     bool
	}{
		{
			name: "simple glob pattern",
			setupDirs: []string{
				"packages/frontend",
				"packages/backend", 
				"packages/shared",
			},
			locations: []Location{
				{
					Name:     "services",
					Location: "packages/*",
					Type:     "npm",
					Commands: []string{"start", "build"},
				},
			},
			expected: []Location{
				{
					Name:     "services",
					Location: "packages/backend",
					Type:     "npm",
					Commands: []string{"start", "build"},
				},
				{
					Name:     "services",
					Location: "packages/frontend",
					Type:     "npm",
					Commands: []string{"start", "build"},
				},
				{
					Name:     "services",
					Location: "packages/shared",
					Type:     "npm",
					Commands: []string{"start", "build"},
				},
			},
			wantErr: false,
		},
		{
			name: "no glob pattern",
			setupDirs: []string{
				"packages/frontend",
			},
			locations: []Location{
				{
					Name:     "frontend",
					Location: "packages/frontend",
					Type:     "npm",
					Commands: []string{"start"},
				},
			},
			expected: []Location{
				{
					Name:     "frontend",
					Location: "packages/frontend",
					Type:     "npm",
					Commands: []string{"start"},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple glob patterns",
			setupDirs: []string{
				"apps/web",
				"apps/mobile",
				"packages/ui",
				"packages/utils",
			},
			locations: []Location{
				{
					Location: "apps/*",
					Type:     "npm",
					Commands: []string{"start"},
				},
				{
					Location: "packages/*",
					Type:     "npm", 
					Commands: []string{"build"},
				},
			},
			expected: []Location{
				{
					Location: "apps/mobile",
					Type:     "npm",
					Commands: []string{"start"},
				},
				{
					Location: "apps/web",
					Type:     "npm",
					Commands: []string{"start"},
				},
				{
					Location: "packages/ui",
					Type:     "npm",
					Commands: []string{"build"},
				},
				{
					Location: "packages/utils",
					Type:     "npm",
					Commands: []string{"build"},
				},
			},
			wantErr: false,
		},
		{
			name: "glob pattern with no matches",
			setupDirs: []string{
				"packages/frontend",
			},
			locations: []Location{
				{
					Location: "services/*",
					Type:     "go",
					Commands: []string{"run"},
				},
			},
			expected: []Location{},
			wantErr: false,
		},
		{
			name: "glob pattern ignores files",
			setupDirs: []string{
				"packages/frontend",
				"packages/backend",
			},
			setupFiles: []string{
				"packages/file.txt",
				"packages/README.md",
			},
			locations: []Location{
				{
					Location: "packages/*",
					Type:     "npm",
					Commands: []string{"test"},
				},
			},
			expected: []Location{
				{
					Location: "packages/backend",
					Type:     "npm",
					Commands: []string{"test"},
				},
				{
					Location: "packages/frontend",
					Type:     "npm",
					Commands: []string{"test"},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid glob pattern - multiple asterisks",
			locations: []Location{
				{
					Location: "foo/*/bar/*",
					Type:     "npm",
					Commands: []string{"test"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid glob pattern - asterisk not at end",
			locations: []Location{
				{
					Location: "packages/*/src",
					Type:     "npm",
					Commands: []string{"test"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid glob pattern - asterisk without slash",
			locations: []Location{
				{
					Location: "packages*",
					Type:     "npm",
					Commands: []string{"test"},
				},
			},
			wantErr: true,
		},
		{
			name: "valid glob pattern - single asterisk at end",
			setupDirs: []string{
				"apps/web",
				"apps/mobile",
			},
			locations: []Location{
				{
					Location: "apps/*",
					Type:     "npm",
					Commands: []string{"start"},
				},
			},
			expected: []Location{
				{
					Location: "apps/mobile",
					Type:     "npm",
					Commands: []string{"start"},
				},
				{
					Location: "apps/web",
					Type:     "npm",
					Commands: []string{"start"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "gopm-glob-test")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			// Create test directories
			for _, dir := range tt.setupDirs {
				err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755)
				if err != nil {
					t.Fatalf("Failed to create dir %s: %v", dir, err)
				}
			}

			// Create test files
			for _, file := range tt.setupFiles {
				err := os.WriteFile(filepath.Join(tmpDir, file), []byte("test"), 0644)
				if err != nil {
					t.Fatalf("Failed to create file %s: %v", file, err)
				}
			}

			// Change to temp directory for glob expansion
			oldWd, err := os.Getwd()
			if err != nil {
				t.Fatalf("Failed to get working directory: %v", err)
			}
			defer os.Chdir(oldWd)

			err = os.Chdir(tmpDir)
			if err != nil {
				t.Fatalf("Failed to change to temp directory: %v", err)
			}

			result, err := ExpandGlobPatterns(tt.locations)

			if (err != nil) != tt.wantErr {
				t.Errorf("ExpandGlobPatterns() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(result) != len(tt.expected) {
					t.Errorf("Expected %d locations, got %d", len(tt.expected), len(result))
					return
				}

				for i, loc := range result {
					expected := tt.expected[i]
					if loc.Name != expected.Name {
						t.Errorf("Location[%d].Name = %q, expected %q", i, loc.Name, expected.Name)
					}
					if loc.Location != expected.Location {
						t.Errorf("Location[%d].Location = %q, expected %q", i, loc.Location, expected.Location)
					}
					if loc.Type != expected.Type {
						t.Errorf("Location[%d].Type = %q, expected %q", i, loc.Type, expected.Type)
					}
					if len(loc.Commands) != len(expected.Commands) {
						t.Errorf("Location[%d] has %d commands, expected %d", i, len(loc.Commands), len(expected.Commands))
					}
					for j, cmd := range loc.Commands {
						if cmd != expected.Commands[j] {
							t.Errorf("Location[%d].Commands[%d] = %q, expected %q", i, j, cmd, expected.Commands[j])
						}
					}
				}
			}
		})
	}
}