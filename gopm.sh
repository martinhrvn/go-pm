#!/usr/bin/env bash

# gopm.sh - Shell script wrapper for gopm
# This script calls the gopm binary to select a command and executes it

set -e

# Function to show usage
show_usage() {
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  run        Interactive command selection and execution"
    echo "  list       List all available commands"
    echo "  help       Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 run"
    echo "  $0 list"
}

# Function to run command interactively
run_command() {
    # Check if jq is available
    if ! command -v jq &> /dev/null; then
        echo "Error: jq is required but not installed. Please install jq to use gopm." >&2
        exit 1
    fi

    # Get the gopm binary path (assume it's in the same directory as this script)
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    GOPM_BINARY="$SCRIPT_DIR/gopm"
    
    # Check if gopm binary exists
    if [ ! -f "$GOPM_BINARY" ]; then
        echo "Error: gopm binary not found at $GOPM_BINARY" >&2
        echo "Please build the Go binary first: go build -o gopm" >&2
        exit 1
    fi

    # Call gopm select to get the selection as JSON
    echo "Selecting command..." >&2
    SELECTION_JSON=$("$GOPM_BINARY" select 2>/dev/null)
    
    # Check if selection was successful
    if [ $? -ne 0 ] || [ -z "$SELECTION_JSON" ]; then
        echo "Selection cancelled or failed." >&2
        exit 1
    fi

    # Parse JSON to extract directory and command
    DIRECTORY=$(echo "$SELECTION_JSON" | jq -r '.directory')
    COMMAND=$(echo "$SELECTION_JSON" | jq -r '.command')

    # Validate parsed values
    if [ "$DIRECTORY" = "null" ] || [ "$COMMAND" = "null" ]; then
        echo "Error: Failed to parse selection from gopm output" >&2
        exit 1
    fi

    # Check if directory exists
    if [ ! -d "$DIRECTORY" ]; then
        echo "Error: Directory '$DIRECTORY' does not exist" >&2
        exit 1
    fi

    # Show what we're about to do
    echo "Running command: $COMMAND"
    echo "In directory: $DIRECTORY"
    echo ""

    # Change to the directory and run the command
    cd "$DIRECTORY"
    exec bash -c "$COMMAND"
}

# Function to list commands
list_commands() {
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    GOPM_BINARY="$SCRIPT_DIR/gopm"
    
    if [ ! -f "$GOPM_BINARY" ]; then
        echo "Error: gopm binary not found at $GOPM_BINARY" >&2
        echo "Please build the Go binary first: go build -o gopm" >&2
        exit 1
    fi

    "$GOPM_BINARY" list
}

# Main script logic
case "${1:-run}" in
    run)
        run_command
        ;;
    list)
        list_commands
        ;;
    help|--help|-h)
        show_usage
        ;;
    *)
        echo "Error: Unknown command '$1'" >&2
        show_usage
        exit 1
        ;;
esac