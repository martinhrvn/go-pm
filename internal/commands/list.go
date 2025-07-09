package commands

import (
	"fmt"

	"github.com/martin/go-pm/internal/config"
)

// ListCommands returns a slice of location:command pairs
func ListCommands(cfg *config.Config) []string {
	var commands []string
	
	for _, location := range cfg.Locations {
		// Use name if available, otherwise use location path
		displayName := location.Name
		if displayName == "" {
			displayName = location.Location
		}
		
		// Add each command for this location
		for _, command := range location.Commands {
			commands = append(commands, fmt.Sprintf("%s:%s", displayName, command))
		}
	}
	
	return commands
}

// FormatForFzf returns a slice of commands formatted for fzf selection
// Format: [location-or-name] command
func FormatForFzf(cfg *config.Config) []string {
	var commands []string
	
	for _, location := range cfg.Locations {
		// Use name if available, otherwise use location path
		displayName := location.Name
		if displayName == "" {
			displayName = location.Location
		}
		
		// Add each command for this location in fzf format
		for _, command := range location.Commands {
			commands = append(commands, fmt.Sprintf("[%s] %s", displayName, command))
		}
	}
	
	return commands
}