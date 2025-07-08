package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Locations []Location `yaml:"locations"`
}

type Location struct {
	Name     string   `yaml:"name,omitempty"`
	Location string   `yaml:"location"`
	Type     string   `yaml:"type,omitempty"`
	Commands []string `yaml:"commands,omitempty"`
}

func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	// Expand glob patterns in locations
	expandedLocations, err := ExpandGlobPatterns(config.Locations)
	if err != nil {
		return nil, err
	}
	config.Locations = expandedLocations

	return &config, nil
}