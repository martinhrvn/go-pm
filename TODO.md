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
- [x] Config file discovery (search current dir and parents)
- [x] Validate config file structure
- [x] Handle malformed config gracefully

### Command Selection
- [x] Integration with fzf for fuzzy selection
- [x] Display format: `[location-or-name] command`
- [x] Support for keyboard interrupts (Ctrl+C)
- [x] Show helpful error messages
- [ ] The cli interface should have following commands
  - [x] gopm list - output all available location:command pairs
  - [x] gopm list --format=fzf - format for fzf selection
  - [ ] gopm get --location=X --command=Y - get execution details as JSON
  - [x] gopm help - show usage and available commands
  - [x] Handle command-line argument parsing

### Command Execution
The go part should handle config parsing, command selection, the execution part should be done in some shell script.
- [x] Change directory to the specified location before execution
- [x] Execute the selected command
- [x] Handle command failures appropriately
- [x] Support complex commands with arguments
- [x] Shell script integration for execution


### Project type implementation
The project types we should support intitally are:
- [x] npm - should find package.json and run npm commands
  - [x] should allow to set the package manager to use, eg. npm, yarn, pnpm
- [ ] go - should find go.mod and run go commands


## Nice-to-Have Features

### Enhanced Configuration
- [ ] Support for environment variables in commands
- [ ] Support for command aliases/shortcuts
- [ ] Support for default commands per location
- [ ] Support for inheriting/extending configurations
- [ ] Support JSON/TOML formats as alternatives
- [ ] Support for comments in config file

### Enhanced UI/UX
- [ ] Preview window in fzf showing command details
- [ ] Most recently used (MRU) commands at top
- [ ] Command history tracking
- [ ] Dry-run mode to preview what will be executed
- [ ] Verbose mode for debugging
- [ ] Do not show locations that do not exist
- [ ] ability to focus specific locations
  - either using .gopmrc.local.yaml (something like `only: location/foo` or command line argument
  - keyboard shortcut to focus/unfocus while searching
- [ ] aliases
- [ ] automatically detect type of a location based on presence of package.json/go.mod/etc.
- [ ] history of executed commands
   - [ ] store in a file (probably some home directory config, but per "project")
   - [ ] default sort by frecency of use
   - [ ] allow to change sorting by keboard shortcut
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
- [x] Written in Go
- [x] Minimal external dependencies (only gopkg.in/yaml.v3 and go-fuzzyfinder)
- [x] Use standard library where possible
- [ ] Ensure cross-platform compatibility (Linux/Mac/Windows)

### Distribution Strategy
- **Two-component approach**: Go binary (`gopm-bin`) + shell wrapper (`gopm`)
- **Flexible installation**: Auto-detects binary location in multiple paths
- **Shell integration**: Bash/zsh completion and colored output
- **Easy installation**: Single script installation with `install.sh`
- **User-friendly**: No need to manually manage binary paths

### Code Structure
- [x] Clear separation of concerns
- [x] Organized folder structure (cmd/, internal/)
- [x] Config parsing module (internal/config)
- [x] Command execution module (internal/commands)
- [x] FZF integration module (internal/commands)
- [x] Project type system (internal/projecttypes)
- [x] Error handling throughout
- [x] Unit tests for core functionality

### Distribution
- [x] Enhanced shell wrapper script (gopm.sh)
- [x] Installation script (install.sh)
- [x] Shell completion (bash/zsh)
- [x] Auto-detection of binary location
- [x] Colored output and better UX
- [x] Nix flake with complete package and development environment
- [x] Cross-platform support (Linux/macOS)
- [ ] Usage documentation
- [ ] Example .gopmrc files
- [ ] GitHub releases with binaries
- [ ] Homebrew formula (optional)
- [ ] AUR package (optional)

## Testing Checklist

### Unit Tests
- [x] Config file parsing tests
- [x] Glob pattern expansion tests
- [x] Project type parsing tests (npm, yarn, pnpm)
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
- [x] Core features working
- [x] Project type system implemented (npm, yarn, pnpm)
- [x] Shell script integration complete
- [x] Testing complete
- [ ] Documentation complete
- [ ] First release

## Notes for Development
- Start with MVP: config parsing, fzf selection, command execution
- Iterate based on real usage patterns
- Keep the tool fast and responsive
- Maintain backwards compatibility with config format