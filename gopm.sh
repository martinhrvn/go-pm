#!/usr/bin/env bash

# gopm - Go Project Manager
# Shell wrapper for gopm binary that handles command selection and execution

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_error() { echo -e "${RED}Error:${NC} $1" >&2; }
print_success() { echo -e "${GREEN}$1${NC}"; }
print_info() { echo -e "${YELLOW}$1${NC}"; }

# Function to find gopm binary
find_gopm_binary() {
    # Try different locations in order of preference
    local candidates=(
        "$(dirname "$0")/gopm"           # Same directory as script
        "$(dirname "$0")/gopm-bin"       # Same directory with -bin suffix
        "$(command -v gopm-bin 2>/dev/null)" # In PATH with -bin suffix
        "$HOME/.local/bin/gopm-bin"      # User local bin
        "/usr/local/bin/gopm-bin"        # System local bin
    )

    for candidate in "${candidates[@]}"; do
        if [ -n "$candidate" ] && [ -f "$candidate" ] && [ -x "$candidate" ]; then
            echo "$candidate"
            return 0
        fi
    done

    return 1
}

# Function to show usage
show_usage() {
    echo "gopm - Go Project Manager"
    echo
    echo "USAGE:"
    echo "    gopm [command]"
    echo
    echo "COMMANDS:"
    echo "    run        Interactive command selection and execution (default)"
    echo "    list       List all available commands"
    echo "    help       Show this help message"
    echo
    echo "EXAMPLES:"
    echo "    gopm           # Interactive selection and execution"
    echo "    gopm run       # Same as above"
    echo "    gopm list      # List all commands"
    echo
    echo "CONFIGURATION:"
    echo "    gopm looks for .gopmrc files starting from the current directory"
    echo "    and traversing up the directory tree until it finds one."
}

# Function to check dependencies
check_dependencies() {
    if ! command -v jq &> /dev/null; then
        print_error "jq is required but not installed."
        echo "Please install jq to use gopm:"
        echo "  Ubuntu/Debian: sudo apt-get install jq"
        echo "  macOS: brew install jq"
        echo "  CentOS/RHEL: sudo yum install jq"
        exit 1
    fi
}

# Function to run command interactively
run_command() {
    check_dependencies

    # Find the gopm binary
    GOPM_BINARY=$(find_gopm_binary)
    if [ $? -ne 0 ]; then
        print_error "gopm binary not found."
        echo "Please ensure gopm is installed or the binary is in your PATH."
        echo "If you're in development, run: go build -o gopm"
        exit 1
    fi

    # Call gopm select to get the selection as JSON
    print_info "Loading configuration and starting selection..."
    if ! SELECTION_JSON=$("$GOPM_BINARY" select 2>/dev/null); then
        print_error "Selection cancelled or failed."
        exit 1
    fi

    # Check if we got valid JSON
    if [ -z "$SELECTION_JSON" ]; then
        print_error "No selection made."
        exit 1
    fi

    # Parse JSON to extract directory and command
    DIRECTORY=$(echo "$SELECTION_JSON" | jq -r '.directory')
    COMMAND=$(echo "$SELECTION_JSON" | jq -r '.command')

    # Validate parsed values
    if [ "$DIRECTORY" = "null" ] || [ "$COMMAND" = "null" ]; then
        print_error "Failed to parse selection from gopm output."
        exit 1
    fi

    # Check if directory exists
    if [ ! -d "$DIRECTORY" ]; then
        print_error "Directory '$DIRECTORY' does not exist."
        exit 1
    fi

    # Show what we're about to do
    print_info "Running: $COMMAND"
    print_info "In: $DIRECTORY"
    echo

    # Change to the directory and run the command
    cd "$DIRECTORY"
    exec bash -c "$COMMAND"
}

# Function to list commands
list_commands() {
    GOPM_BINARY=$(find_gopm_binary)
    if [ $? -ne 0 ]; then
        print_error "gopm binary not found."
        echo "Please ensure gopm is installed or the binary is in your PATH."
        exit 1
    fi

    "$GOPM_BINARY" list
}

# Main script logic
case "${1:-run}" in
    run|"")
        run_command
        ;;
    list)
        list_commands
        ;;
    help|--help|-h)
        show_usage
        ;;
    *)
        print_error "Unknown command '$1'"
        show_usage
        exit 1
        ;;
esac