package confmanager

import (
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
		Type    string `yaml:"TYPE"`
		Options struct {
			Patterns string `yaml:"PATTERN"`
			Ignore_Case bool   `yaml:"IGNORE_CASE"`
		} `yaml:"OPTIONS"`
	} `yaml:"FILTER"`
	Output struct {
		Type string `yaml:"TYPE"`
	} `yaml:"OUTPUT"`
} 


func ReadConfig(filepath string) (Config, error) {
	yamlFile, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	if len(yamlFile) == 0 {
		panic("Config file is empty")
	}
	var config Config
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		panic(err)
	}
	return config, nil
}