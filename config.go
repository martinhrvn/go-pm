package main

type Config struct {
	Locations []Location `yaml:"locations"`
}

type Location struct {
	Name     string   `yaml:"name,omitempty"`
	Location string   `yaml:"location"`
	Type     string   `yaml:"type,omitempty"`
	Commands []string `yaml:"commands,omitempty"`
}