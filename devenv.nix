{ pkgs, ... }:

{
  packages = with pkgs; [
    go
    gopls
    delve
    golangci-lint
    gotools
  ];

  languages.go = {
    enable = true;
  };

  scripts.build.exec = "go build -o bin/go-pm ./cmd/go-pm";
  scripts.test.exec = "go test ./...";
  scripts.lint.exec = "golangci-lint run";
  scripts.dev.exec = "go run ./cmd/go-pm";

  enterShell = ''
    echo "Go development environment loaded"
    echo "Available commands: build, test, lint, dev"
  '';
}