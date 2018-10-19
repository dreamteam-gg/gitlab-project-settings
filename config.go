package main

import (
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	GitLabUrl       string                 `yaml:"gitlab_url"`
	GitLabToken     string                 `yaml:"gitlab_private_token"`
	StopOnError     bool                   `yaml:"stop_on_error"`
	NamespaceID     int                    `yaml:"namespace_id"`
	Settings        map[string]interface{} `yaml:"settings"`
	OnlyProject     []string               `yaml:"only_projects"`
	ExcludeProjects []string               `yaml:"exclude_projects"`
}

func ConfigFromFile(file string) (*Config, error) {
	configBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var c *Config
	if err := yaml.Unmarshal(configBytes, &c); err != nil {
		return nil, err
	}
	c.GitLabUrl = strings.TrimSuffix(c.GitLabUrl, "/")
	return c, nil
}
