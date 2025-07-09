{
  description = "gopm - Go Project Manager for monorepos";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        
        # Build the Go binary
        gopm-bin = pkgs.buildGoModule {
          pname = "gopm-bin";
          version = "0.1.0";
          
          src = ./.;
          
          vendorHash = "sha256-Lh/yde3GkhNDd/IH0rFj80hpkwLbLzE5N1bAUKpHMEQ=";
          
          # Build flags
          ldflags = [
            "-s"
            "-w"
            "-X main.version=0.1.0"
          ];
          
          # Test the binary
          doCheck = true;
          
          # Only build the main package
          subPackages = [ "cmd/gopm" ];
          
          meta = with pkgs.lib; {
            description = "Go Project Manager - A utility for quickly running commands in monorepos";
            homepage = "https://github.com/martin/go-pm";
            license = licenses.mit;
            maintainers = [ "martin" ];
            platforms = platforms.unix;
          };
        };
        
        # Create the shell wrapper
        gopm-wrapper = pkgs.writeShellScriptBin "gopm" ''
          # gopm - Go Project Manager
          # Shell wrapper for gopm binary that handles command selection and execution
          
          set -e
          
          # Colors for output
          RED='\033[0;31m'
          GREEN='\033[0;32m'
          YELLOW='\033[1;33m'
          NC='\033[0m' # No Color
          
          # Function to print colored output
          print_error() { echo -e "''${RED}Error:''${NC} $1" >&2; }
          print_success() { echo -e "''${GREEN}$1''${NC}"; }
          print_info() { echo -e "''${YELLOW}$1''${NC}"; }
          
          # Path to the gopm binary
          GOPM_BINARY="${gopm-bin}/bin/gopm"

          echo $GOPM_BINARY
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
              if ! command -v ${pkgs.jq}/bin/jq &> /dev/null; then
                  print_error "jq is required but not found in PATH."
                  echo "This should not happen in a Nix environment."
                  exit 1
              fi
          }
          
          # Function to run command interactively
          run_command() {
              check_dependencies
          
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
              DIRECTORY=$(echo "$SELECTION_JSON" | ${pkgs.jq}/bin/jq -r '.directory')
              COMMAND=$(echo "$SELECTION_JSON" | ${pkgs.jq}/bin/jq -r '.command')
          
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
              exec ${pkgs.bash}/bin/bash -c "$COMMAND"
          }
          
          # Function to list commands
          list_commands() {
              "$GOPM_BINARY" list
          }
          
          # Main script logic
          case "''${1:-run}" in
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
        '';
        
        # Complete gopm package with binary and wrapper
        gopm = pkgs.symlinkJoin {
          name = "gopm";
          paths = [ gopm-bin gopm-wrapper ];
          buildInputs = [ pkgs.makeWrapper ];
          postBuild = ''
            # Make sure jq is available at runtime
            wrapProgram $out/bin/gopm \
              --prefix PATH : ${pkgs.lib.makeBinPath [ pkgs.jq pkgs.bash ]}
          '';
          
          meta = with pkgs.lib; {
            description = "Go Project Manager - A utility for quickly running commands in monorepos";
            homepage = "https://github.com/martin/go-pm";
            license = licenses.mit;
            maintainers = [ "martin" ];
            platforms = platforms.unix;
          };
        };
        
        # Shell completions
        gopm-completions = pkgs.stdenv.mkDerivation {
          name = "gopm-completions";
          src = ./.;
          
          installPhase = ''
            mkdir -p $out/share/bash-completion/completions
            mkdir -p $out/share/zsh/site-functions
            
            cp completion.bash $out/share/bash-completion/completions/gopm
            cp completion.zsh $out/share/zsh/site-functions/_gopm
          '';
          
          meta = with pkgs.lib; {
            description = "Shell completions for gopm";
            platforms = platforms.unix;
          };
        };
        
        # Development shell
        devShell = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            jq
            gopls
            golangci-lint
            delve
          ];
          
          shellHook = ''
            echo "🚀 Welcome to gopm development environment!"
            echo "Available commands:"
            echo "  go build -o gopm     # Build the binary"
            echo "  go test ./...        # Run tests"
            echo "  go run .             # Run directly"
            echo "  ./gopm.sh            # Test the shell wrapper"
            echo ""
            echo "Nix environment includes:"
            echo "  - Go ${pkgs.go.version}"
            echo "  - jq ${pkgs.jq.version}"
            echo "  - gopls (Go Language Server)"
            echo "  - golangci-lint"
            echo "  - delve (Go debugger)"
          '';
        };
        
      in
      {
        packages = {
          default = gopm-wrapper;
          gopm = gopm-wrapper;
          gopm-bin = gopm-bin;
          gopm-wrapper = gopm-wrapper;
          gopm-completions = gopm-completions;
        };
        
        apps = {
          default = flake-utils.lib.mkApp {
            drv = gopm;
            name = "gopm";
          };
          
          gopm = flake-utils.lib.mkApp {
            drv = gopm;
            name = "gopm";
          };
        };
        
        devShells.default = devShell;
        
        # For backwards compatibility
        devShell = devShell;
        defaultPackage = gopm;
      }
    );
}