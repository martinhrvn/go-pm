package main

import (
	"fmt"
	"os"

	"github.com/martin/go-pm/internal/commands"
	"github.com/martin/go-pm/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		showUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "list":
		handleListCommand()
	case "select":
		handleSelectCommand()
	case "help":
		showUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		showUsage()
		os.Exit(1)
	}
}

func handleListCommand() {
	// Check for format flag
	format := "default"
	if len(os.Args) > 2 && os.Args[2] == "--format=fzf" {
		format = "fzf"
	}

	// Load config from discovery
	cfg, err := config.LoadConfigFromDiscovery()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Generate command list
	var cmdList []string
	if format == "fzf" {
		cmdList = commands.FormatForFzf(cfg)
	} else {
		cmdList = commands.ListCommands(cfg)
	}

	// Output commands
	for _, cmd := range cmdList {
		fmt.Println(cmd)
	}
}

func handleSelectCommand() {
	// Check for enhanced flag
	enhanced := false
	for _, arg := range os.Args[2:] {
		if arg == "--enhanced" || arg == "-e" {
			enhanced = true
			break
		}
	}

	// Load config from discovery
	cfg, err := config.LoadConfigFromDiscovery()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Run fzf selection (enhanced or regular)
	var result *commands.SelectionResult
	if enhanced {
		result, err = commands.RunEnhancedFzf(cfg)
	} else {
		result, err = commands.RunFzf(cfg)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error with selection: %v\n", err)
		os.Exit(1)
	}

	// Output as JSON for shell script parsing
	fmt.Printf(`{"directory":"%s","command":"%s"}`, result.Directory, result.Command)
	fmt.Println()
}

func showUsage() {
	fmt.Println("gopm - Go Project Manager")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("    gopm <command> [options]")
	fmt.Println()
	fmt.Println("COMMANDS:")
	fmt.Println("    list                     List all available location:command pairs")
	fmt.Println("    list --format=fzf        List commands in fzf format")
	fmt.Println("    select                   Interactive command selection with fzf")
	fmt.Println("    select --enhanced        Enhanced TUI selection with location filtering")
	fmt.Println("    help                     Show this help message")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("    gopm list")
	fmt.Println("    gopm list --format=fzf")
	fmt.Println("    gopm select")
	fmt.Println("    gopm select --enhanced")
}
