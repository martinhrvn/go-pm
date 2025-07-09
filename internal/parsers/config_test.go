package parsers

import (
	"testing"
)

func TestGetDefaultParsers(t *testing.T) {
	defaults := getDefaultParsers()
	
	if defaults.Parsers == nil {
		t.Fatal("Expected parsers to be initialized")
	}
	
	// Test that default parsers are present
	expectedParsers := []string{"npm", "yarn", "pnpm", "go"}
	for _, name := range expectedParsers {
		if _, exists := defaults.Parsers[name]; !exists {
			t.Errorf("Expected parser %s to exist in defaults", name)
		}
	}
	
	// Test npm parser configuration
	npmParser := defaults.Parsers["npm"]
	if len(npmParser.DetectFiles) == 0 {
		t.Error("Expected npm parser to have detect files")
	}
	if npmParser.DetectFiles[0] != "package.json" {
		t.Errorf("Expected npm parser to detect package.json, got %s", npmParser.DetectFiles[0])
	}
	if npmParser.BuiltinParser != "package_json_scripts" {
		t.Errorf("Expected npm parser to use package_json_scripts, got %s", npmParser.BuiltinParser)
	}
	
	// Test go parser configuration
	goParser := defaults.Parsers["go"]
	if len(goParser.DetectFiles) == 0 {
		t.Error("Expected go parser to have detect files")
	}
	if goParser.DetectFiles[0] != "go.mod" {
		t.Errorf("Expected go parser to detect go.mod, got %s", goParser.DetectFiles[0])
	}
	if len(goParser.BaseCommands) == 0 {
		t.Error("Expected go parser to have base commands")
	}
}

func TestParserConfigGetParser(t *testing.T) {
	parsers := getDefaultParsers()
	
	// Test existing parser
	parser, exists := parsers.GetParser("npm")
	if !exists {
		t.Error("Expected npm parser to exist")
	}
	if parser.BuiltinParser != "package_json_scripts" {
		t.Errorf("Expected npm parser to use package_json_scripts, got %s", parser.BuiltinParser)
	}
	
	// Test non-existing parser
	_, exists = parsers.GetParser("nonexistent")
	if exists {
		t.Error("Expected nonexistent parser to not exist")
	}
}