package config

import (
	"fmt"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

type yamlConfig struct {
	Repositories []yamlRepository `yaml:"repositories"`
}

type yamlRepository struct {
	RemoteURL string `yaml:"url"`
	Path      string `yaml:"path"`
}

// Parse a configuration content into a valid Config struct
func Parse(configurationContent []byte) (*Config, error) {
	rawConfiguration := yamlConfig{}
	err := yaml.Unmarshal(configurationContent, &rawConfiguration)
	if err != nil {
		return nil, fmt.Errorf("can't load configuration: %v", err)
	}

	configuration := Config{
		Repositories: []Repository{},
	}

	parseRepositories(&configuration, rawConfiguration)

	return &configuration, nil
}

func parseRepositories(config *Config, rawConfig yamlConfig) error {
	for _, rawRepository := range rawConfig.Repositories {
		rawRepositoryPath := rawRepository.Path
		rawURL := rawRepository.RemoteURL

		validRepositoryPath, err := homedir.Expand(rawRepositoryPath)
		if err != nil {
			return fmt.Errorf("invalid repository path '%s': %s", rawRepositoryPath, err.Error())
		}

		// TODO: should add validations here
		validURL := rawURL

		repository := Repository{RemoteURL: validURL, Path: validRepositoryPath}

		config.Repositories = append(config.Repositories, repository)
	}

	return nil
}
