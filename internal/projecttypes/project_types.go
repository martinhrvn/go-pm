package projecttypes

import (
	"fmt"
	"os"
	"path/filepath"
)

// ProjectType represents a project type that can discover and provide commands
type ProjectType interface {
	// Name returns the name of the project type (e.g., "npm", "go")
	Name() string
	
	// DetectConfigFile returns the config file name to look for (e.g., "package.json", "go.mod")
	DetectConfigFile() string
	
	// ParseCommands parses the config file and returns available commands
	ParseCommands(configPath string) ([]string, error)
	
	// GetCommandPrefix returns the command prefix for this project type (e.g., "npm run", "go")
	GetCommandPrefix() string
}

// ProjectTypeRegistry holds all registered project types
var ProjectTypeRegistry = map[string]ProjectType{
	"npm":  &NpmProjectType{},
	"yarn": &YarnProjectType{},
	"pnpm": &PnpmProjectType{},
}

// GetProjectType returns a project type by name
func GetProjectType(name string) (ProjectType, error) {
	projectType, exists := ProjectTypeRegistry[name]
	if !exists {
		return nil, fmt.Errorf("unknown project type: %s", name)
	}
	return projectType, nil
}

// DiscoverProjectType attempts to discover the project type in a directory
func DiscoverProjectType(directory string) (ProjectType, error) {
	for _, projectType := range ProjectTypeRegistry {
		configFile := filepath.Join(directory, projectType.DetectConfigFile())
		if fileExists(configFile) {
			return projectType, nil
		}
	}
	return nil, fmt.Errorf("no project type detected in directory: %s", directory)
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

