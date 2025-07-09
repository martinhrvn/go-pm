package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const ConfigFileName = ".gopmrc"

// FindConfigFile searches for .gopmrc starting from the current working directory
// and traversing up the directory tree until it finds one or reaches the root.
func FindConfigFile() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	return findConfigFileFromPath(cwd)
}

// findConfigFileFromPath searches for .gopmrc starting from the given path
// and traversing up the directory tree.
func findConfigFileFromPath(startPath string) (string, error) {
	currentPath := startPath
	
	for {
		configPath := filepath.Join(currentPath, ConfigFileName)
		
		// Check if config file exists
		if _, err := os.Stat(configPath); err == nil {
			return configPath, nil
		}
		
		// Move to parent directory
		parentPath := filepath.Dir(currentPath)
		
		// Check if we've reached the root directory
		if parentPath == currentPath {
			break
		}
		
		currentPath = parentPath
	}
	
	return "", fmt.Errorf("no %s found in current directory or any parent directories", ConfigFileName)
}

// LoadConfigFromDiscovery finds and loads the nearest .gopmrc file
func LoadConfigFromDiscovery() (*Config, error) {
	configPath, err := FindConfigFile()
	if err != nil {
		return nil, err
	}
	
	// Change to the directory containing the config file for glob expansion
	configDir := filepath.Dir(configPath)
	oldWd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current working directory: %w", err)
	}
	defer os.Chdir(oldWd)
	
	err = os.Chdir(configDir)
	if err != nil {
		return nil, fmt.Errorf("failed to change to config directory %q: %w", configDir, err)
	}
	
	return LoadConfig(configPath)
}