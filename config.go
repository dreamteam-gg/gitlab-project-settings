package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/imdario/mergo"
	"github.com/spf13/viper"
)

type Config struct {
	GitLabUrl       string                 `yaml:"gitlab_url"`
	GitLabToken     string                 `yaml:"gitlab_private_token"`
	CreateMissing   bool                   `yaml:"create_missing"`
	GroupID         string                 `yaml:"group_id"`
	GroupSettings   map[string]interface{} `yaml:"group_settings"`
	Settings        map[string]interface{} `yaml:"project_settings"`
	Overrides       map[string]interface{} `yaml:"overrides"`
	OnlyProject     []string               `yaml:"only_projects"`
	ExcludeProjects []string               `yaml:"exclude_projects"`
}

func ConfigFromFile(file string) (*Config, error) {
	viper.SetConfigFile(file)
	viper.SetEnvPrefix("gitlab")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}

	viper.AutomaticEnv() // read in environment variables that match

	c := Config{
		GitLabUrl:       strings.TrimSuffix(viper.GetString("gitlab_url"), "/"),
		GitLabToken:     viper.GetString("gitlab_private_token"),
		CreateMissing:   viper.GetBool("create_missing"),
		GroupID:         viper.GetString("group_id"),
		GroupSettings:   viper.GetStringMap("group_settings"),
		Settings:        viper.GetStringMap("project_settings"),
		Overrides:       viper.GetStringMap("overrides"),
		OnlyProject:     viper.GetStringSlice("only_projects"),
		ExcludeProjects: viper.GetStringSlice("exclude_projects"),
	}

	return &c, nil
}

func MergeConfig(dst map[string]interface{}, src map[string]interface{}) error {
	err := mergo.Merge(&dst, src, mergo.WithOverride)
	if err != nil {
		return err
	}

	return nil
}
