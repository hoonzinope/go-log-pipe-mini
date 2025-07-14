package confmanager

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v3"
)

// read configuration from a config.yml file
// and return struct with the configuration

type Config struct {
	Input struct {
		Type string `yaml:"TYPE"`
		Path string `yaml:"PATH"`
	} `yaml:"INPUT"`
	Filter struct {
		Mode    string `yaml:"MODE"`
		Filters []struct {
			Type    string `yaml:"TYPE"`
			Options struct {
				IgnoreCase bool   `yaml:"IGNORE_CASE"`
				Pattern    string `yaml:"PATTERN"`
			} `yaml:"OPTIONS"`
		} `yaml:"FILTERS"`
	} `yaml:"FILTER"`
	Output struct {
		Type string `yaml:"TYPE"`
	} `yaml:"OUTPUT"`
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
