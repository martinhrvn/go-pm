package parsers

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ParserConfig represents a single parser configuration
type ParserConfig struct {
	// DetectFiles are the files to look for to detect this project type
	DetectFiles []string `yaml:"detect_files"`
	
	// BaseCommands are commands that are always available, regardless of parsing
	BaseCommands map[string]string `yaml:"base_commands"`
	
	// BuiltinParser specifies a built-in parser to use (e.g., "package_json_scripts")
	BuiltinParser string `yaml:"builtin_parser,omitempty"`
	
	// ParserCommand is a shell command that outputs available commands (one per line)
	ParserCommand string `yaml:"parser_command,omitempty"`
	
	// CommandTemplate is how to construct the final command (e.g., "npm run {key}")
	CommandTemplate string `yaml:"command_template,omitempty"`
}

// ParsersFile represents the entire parsers.yaml configuration
type ParsersFile struct {
	Parsers map[string]ParserConfig `yaml:"parsers"`
}

// LoadParsersConfig loads parser configuration from ~/.gopm/parsers.yaml
func LoadParsersConfig() (*ParsersFile, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".gopm", "parsers.yaml")
	
	// If config doesn't exist, return default parsers
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return getDefaultParsers(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read parsers config: %w", err)
	}

	var config ParsersFile
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse parsers config: %w", err)
	}

	// Merge with defaults
	defaults := getDefaultParsers()
	if config.Parsers == nil {
		config.Parsers = make(map[string]ParserConfig)
	}
	
	// Add default parsers if not overridden
	for name, parser := range defaults.Parsers {
		if _, exists := config.Parsers[name]; !exists {
			config.Parsers[name] = parser
		}
	}

	return &config, nil
}

// getDefaultParsers returns the default built-in parser configurations
func getDefaultParsers() *ParsersFile {
	return &ParsersFile{
		Parsers: map[string]ParserConfig{
			"npm": {
				DetectFiles: []string{"package.json"},
				BaseCommands: map[string]string{
					"install": "npm install",
					"audit":   "npm audit",
				},
				BuiltinParser:   "package_json_scripts",
				CommandTemplate: "npm run {key}",
			},
			"yarn": {
				DetectFiles: []string{"package.json"},
				BaseCommands: map[string]string{
					"install": "yarn install",
					"audit":   "yarn audit",
				},
				BuiltinParser:   "package_json_scripts",
				CommandTemplate: "yarn {key}",
			},
			"pnpm": {
				DetectFiles: []string{"package.json"},
				BaseCommands: map[string]string{
					"install": "pnpm install",
					"audit":   "pnpm audit",
				},
				BuiltinParser:   "package_json_scripts",
				CommandTemplate: "pnpm run {key}",
			},
			"go": {
				DetectFiles: []string{"go.mod"},
				BaseCommands: map[string]string{
					"build": "go build ./...",
					"test":  "go test ./...",
					"fmt":   "go fmt ./...",
					"vet":   "go vet ./...",
					"mod":   "go mod tidy",
				},
			},
		},
	}
}

// GetParser returns a parser configuration by name
func (p *ParsersFile) GetParser(name string) (ParserConfig, bool) {
	parser, exists := p.Parsers[name]
	return parser, exists
}

// FindParserForDirectory finds a parser that matches files in the given directory
func (p *ParsersFile) FindParserForDirectory(directory string) (string, ParserConfig, error) {
	for name, parser := range p.Parsers {
		for _, detectFile := range parser.DetectFiles {
			filePath := filepath.Join(directory, detectFile)
			if _, err := os.Stat(filePath); err == nil {
				return name, parser, nil
			}
		}
	}
	return "", ParserConfig{}, fmt.Errorf("no parser found for directory: %s", directory)
}