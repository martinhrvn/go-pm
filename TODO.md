# CLAUDE.md - gopm Development Checklist

## Project Overview
**gopm** (Go Project Manager) - A utility for quickly running commands in monorepos using fzf selection.

## Core Requirements

### Configuration File (.gopmrc)
- [x] Support YAML-like format parsing
- [x] Support for `locations` array
- [x] Support for `name` field (optional display name)
- [x] Support for `location` field (path, can include globs)
- [x] Support for `type` field (optional, for package manager context)
  - initially support `npm`, `yarn`, `pnpm`, `go`
  - this basically means in the folder it will automatically find the appropriate package manager config file
  - generate list of commands based on the type
- [x] Support for `commands` array
  - this can be specified as a list of commands, if type is specified for project it will add this as extra commands
- [x] Support glob patterns in location paths (e.g., `packages/bar/*`)
    - This should be simple eg. expand * to all directories in the path, but not recurse into subdirectories
- [ ] Config file discovery (search current dir and parents)
- [x] Validate config file structure
- [x] Handle malformed config gracefully

### Command Selection
- [ ] Integration with fzf for fuzzy selection
- [ ] Display format: `[location-or-name] command`
- [ ] Support for keyboard interrupts (Ctrl+C)
- [ ] Show helpful error messages
- [ ] The cli interface should have following commands
-  gopm list - output all available location:command pairs
   gopm list --format=fzf - format for fzf selection
   gopm get --location=X --command=Y - get execution details as JSON
   gopm help - show usage and available commands
   Handle command-line argument parsing

### Command Execution
The go part should handle config parsing, command selection, the execution part should be done in some shell script.
- [ ] Change directory to the specified location before execution
- [ ] Execute the selected command
- [ ] Pass through stdin/stdout/stderr
- [ ] Handle command failures appropriately
- [ ] Support complex commands with arguments

## Nice-to-Have Features

### Enhanced Configuration
- [ ] Support for environment variables in commands
- [ ] Support for command aliases/shortcuts
- [ ] Support for default commands per location
- [ ] Support for inheriting/extending configurations
- [ ] Support JSON/TOML formats as alternatives
- [ ] Support for comments in config file

### Enhanced UI/UX
- [ ] Colorized output
- [ ] Preview window in fzf showing command details
- [ ] Most recently used (MRU) commands at top
- [ ] Command history tracking
- [ ] Dry-run mode to preview what will be executed
- [ ] Verbose mode for debugging

### Advanced Features
- [ ] Support for pre/post command hooks
- [ ] Support for command templates/variables
- [ ] Support for running commands in parallel
- [ ] Support for command groups/categories
- [ ] Integration with different package managers based on `type`
- [ ] Watch mode for repeated command execution
- [ ] Shell completion (bash/zsh/fish)

## Implementation Details

### Language & Dependencies
- [ ] Written in Go
- [ ] Minimal external dependencies
- [ ] Use standard library where possible
- [ ] Ensure cross-platform compatibility (Linux/Mac/Windows)

### Code Structure
- [x] Clear separation of concerns
- [x] Config parsing module
- [ ] Command execution module
- [ ] FZF integration module
- [x] Error handling throughout
- [x] Unit tests for core functionality

### Distribution
- [ ] Single binary distribution
- [ ] Installation instructions
- [ ] Usage documentation
- [ ] Example .gopmrc files
- [ ] Homebrew formula (optional)
- [ ] AUR package (optional)

## Testing Checklist

### Unit Tests
- [x] Config file parsing tests
- [x] Glob pattern expansion tests
- [ ] Command parsing tests
- [x] Error handling tests

### Integration Tests
- [x] Test with real monorepo structure
- [ ] Test with various .gopmrc formats
- [ ] Test with missing fzf
- [x] Test with malformed configs
- [ ] Test with non-existent locations

### Edge Cases
- [x] Empty config file
- [ ] No commands defined
- [x] Invalid glob patterns
- [ ] Commands with special characters
- [ ] Very long command lists
- [ ] Nested monorepo structures

## Documentation

### README.md
- [ ] Clear installation instructions
- [ ] Usage examples
- [ ] Configuration format documentation
- [ ] Troubleshooting section
- [ ] Contributing guidelines

### Examples
- [ ] Example .gopmrc for npm/yarn/pnpm monorepo
- [ ] Example .gopmrc for Go workspace
- [ ] Example .gopmrc for mixed-language monorepo
- [ ] Advanced configuration examples

## Release Checklist
- [ ] Version tagging
- [ ] Changelog updates
- [ ] Binary releases for major platforms
- [ ] Release notes
- [ ] Update documentation

## Current Status
- [x] Initial concept defined
- [x] Basic implementation started
- [ ] Core features working
- [ ] Testing complete
- [ ] Documentation complete
- [ ] First release

## Notes for Development
- Start with MVP: config parsing, fzf selection, command execution
- Iterate based on real usage patterns
- Keep the tool fast and responsive
- Maintain backwards compatibility with config format