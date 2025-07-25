package confmanager

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v3"
)

type Config struct {
	Inputs  []InputConfig          `yaml:"INPUTS"`
	Filters map[string]FilterGroup `yaml:"FILTERS"`
	Outputs []OutputConfig         `yaml:"OUTPUTS"`
}

type InputConfig struct {
	Name   string `yaml:"NAME"`
	Type   string `yaml:"TYPE"`
	Path   string `yaml:"PATH"`
	Parser string `yaml:"PARSER"`
}

type FilterGroup struct {
	Mode  string       `yaml:"MODE"`
	Rules []FilterRule `yaml:"RULES"`
}

type FilterRule struct {
	Type    string `yaml:"TYPE"`
	Options struct {
		IgnoreCase bool   `yaml:"IGNORE_CASE"`
		Pattern    string `yaml:"PATTERN"`
		Field      string `yaml:"FIELD"`
	} `yaml:"OPTIONS"`
}

type OutputConfig struct {
	Type    string   `yaml:"TYPE"`
	Targets []string `yaml:"TARGETS"`
	Options struct {
		Path     string `yaml:"PATH"`
		Filename string `yaml:"FILENAME"`  // e.g., "output.log"
		Rolling  string `yaml:"ROLLING"`   // daily, hourly, monthly
		MaxSize  string `yaml:"MAX_SIZE"`  // e.g., "100MB"
		MaxFiles int    `yaml:"MAX_FILES"` // e.g., 7
		BATCH_SIZE int  `yaml:"BATCH_SIZE"` // e.g., 10
		FLUSH_INTERVAL string `yaml:"FLUSH_INTERVAL"` // e.g., "5s"
	} `yaml:"OPTIONS"`
}

func ReadConfig(filepath string) (Config, error) {
	yamlFile, err := os.ReadFile(filepath)
	if err != nil {
		return Config{}, err
	}
	if len(yamlFile) == 0 {
		return Config{}, fmt.Errorf("config file %s is empty", filepath)
	}
	var config Config
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		return Config{}, fmt.Errorf("error parsing config file %s: %v", filepath, err)
	}
	return config, nil
}
