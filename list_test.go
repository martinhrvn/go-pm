package main

import (
	"strings"
	"testing"
)

func TestListCommands(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		expected []string
	}{
		{
			name: "single location with multiple commands",
			config: &Config{
				Locations: []Location{
					{
						Name:     "frontend",
						Location: "packages/frontend",
						Commands: []string{"start", "build", "test"},
					},
				},
			},
			expected: []string{
				"frontend:start",
				"frontend:build", 
				"frontend:test",
			},
		},
		{
			name: "multiple locations with commands",
			config: &Config{
				Locations: []Location{
					{
						Name:     "frontend",
						Location: "packages/frontend",
						Commands: []string{"start", "build"},
					},
					{
						Name:     "backend",
						Location: "packages/backend",
						Commands: []string{"run", "test"},
					},
				},
			},
			expected: []string{
				"frontend:start",
				"frontend:build",
				"backend:run",
				"backend:test",
			},
		},
		{
			name: "location without name uses location path",
			config: &Config{
				Locations: []Location{
					{
						Location: "packages/frontend",
						Commands: []string{"start", "build"},
					},
					{
						Name:     "backend",
						Location: "packages/backend",
						Commands: []string{"run"},
					},
				},
			},
			expected: []string{
				"packages/frontend:start",
				"packages/frontend:build",
				"backend:run",
			},
		},
		{
			name: "location without commands",
			config: &Config{
				Locations: []Location{
					{
						Name:     "frontend",
						Location: "packages/frontend",
					},
					{
						Name:     "backend",
						Location: "packages/backend",
						Commands: []string{"run"},
					},
				},
			},
			expected: []string{
				"backend:run",
			},
		},
		{
			name: "empty config",
			config: &Config{
				Locations: []Location{},
			},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ListCommands(tt.config)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d commands, got %d", len(tt.expected), len(result))
				return
			}

			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("Command[%d] = %q, expected %q", i, result[i], expected)
				}
			}
		})
	}
}

func TestFormatForFzf(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		expected []string
	}{
		{
			name: "single location with commands",
			config: &Config{
				Locations: []Location{
					{
						Name:     "frontend",
						Location: "packages/frontend",
						Commands: []string{"start", "build"},
					},
				},
			},
			expected: []string{
				"[frontend] start",
				"[frontend] build",
			},
		},
		{
			name: "multiple locations",
			config: &Config{
				Locations: []Location{
					{
						Name:     "frontend",
						Location: "packages/frontend",
						Commands: []string{"start"},
					},
					{
						Location: "packages/backend",
						Commands: []string{"run"},
					},
				},
			},
			expected: []string{
				"[frontend] start",
				"[packages/backend] run",
			},
		},
		{
			name: "empty config",
			config: &Config{
				Locations: []Location{},
			},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatForFzf(tt.config)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d commands, got %d", len(tt.expected), len(result))
				return
			}

			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("Command[%d] = %q, expected %q", i, result[i], expected)
				}
			}
		})
	}
}

func TestListCommandsOutput(t *testing.T) {
	config := &Config{
		Locations: []Location{
			{
				Name:     "frontend",
				Location: "packages/frontend",
				Commands: []string{"start", "build"},
			},
			{
				Location: "packages/backend",
				Commands: []string{"run", "test"},
			},
		},
	}

	tests := []struct {
		name           string
		format         string
		expectedOutput string
	}{
		{
			name:   "default format",
			format: "default",
			expectedOutput: `frontend:start
frontend:build
packages/backend:run
packages/backend:test`,
		},
		{
			name:   "fzf format",
			format: "fzf",
			expectedOutput: `[frontend] start
[frontend] build
[packages/backend] run
[packages/backend] test`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result string
			
			if tt.format == "fzf" {
				commands := FormatForFzf(config)
				result = strings.Join(commands, "\n")
			} else {
				commands := ListCommands(config)
				result = strings.Join(commands, "\n")
			}

			if result != tt.expectedOutput {
				t.Errorf("Expected output:\n%s\n\nGot:\n%s", tt.expectedOutput, result)
			}
		})
	}
}