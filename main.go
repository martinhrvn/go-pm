package main

import (
	"fmt"
	"os"
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
	config, err := LoadConfigFromDiscovery()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Generate command list
	var commands []string
	if format == "fzf" {
		commands = FormatForFzf(config)
	} else {
		commands = ListCommands(config)
	}

	// Output commands
	for _, cmd := range commands {
		fmt.Println(cmd)
	}
}

func handleSelectCommand() {
	// Load config from discovery
	config, err := LoadConfigFromDiscovery()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Run fzf selection
	result, err := RunFzf(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error with selection: %v\n", err)
		os.Exit(1)
	}

	// Display the result
	fmt.Printf("Selected:\n")
	fmt.Printf("  Directory: %s\n", result.Directory)
	fmt.Printf("  Command: %s\n", result.Command)
	fmt.Printf("  Display Name: %s\n", result.DisplayName)
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
	fmt.Println("    help                     Show this help message")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("    gopm list")
	fmt.Println("    gopm list --format=fzf")
	fmt.Println("    gopm select")
}