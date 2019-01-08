package main

import (
	"flag"
	"fmt"

	"github.com/logrusorgru/aurora"
)

var (
	flagConfigFile = flag.String("config", "./config.yml", "Path to configuration file")
	flagDryRun     = flag.Bool("dry-run", false, "Dry run mode")
	cfg            *Config
	formatter      aurora.Aurora
)

func init() {
	formatter = aurora.NewAurora(isTerminal())
}

func InArray(k string, arr []string) bool {
	for _, val := range arr {
		if val == k {
			return true
		}
	}
	return false
}

func InProjects(k string, arr []*Project) bool {
	for _, p := range arr {
		if k == p.Get("name").(string) {
			return true
		}
	}
	return false
}

func CanProcessProject(name string) bool {
	if InArray(name, cfg.ExcludeProjects) {
		return false
	}
	if len(cfg.OnlyProject) == 0 {
		return true
	}
	if InArray(name, cfg.OnlyProject) {
		return true
	}
	return false
}

func IsEqual(a, b interface{}) bool {
	if fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b) {
		return true
	}
	return false
}

func main() {
	var err error

	flag.Parse()

	cfg, err = ConfigFromFile(*flagConfigFile)
	if err != nil {
		panic(err)
	}

	client := NewClient(cfg.GitLabUrl, cfg.GitLabToken)

	groupId, err := client.GetGroupIdByName(cfg.GroupID)
	if err != nil {
		panic(err)
	}

	projects, err := client.GetGroupProjects(groupId)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Found %d projects\n", len(projects))
	for _, project := range projects {
		name := project.Get("name").(string)
		if !CanProcessProject(name) {
			continue
		}
		settings := cfg.Settings
		if v, ok := cfg.Overrides[name]; ok {
			err := MergeConfig(settings, v.(map[string]interface{}))
			if err != nil {
				panic(err)
			}
		}
		err = client.UpdateProject(project, settings)
		if err != nil {
			fmt.Println(err)
			if cfg.StopOnError {
				break
			}
		}
	}

	if cfg.CreateMissing {
		for _, name := range cfg.OnlyProject {
			if InProjects(name, projects) {
				continue
			}

			settings := cfg.Settings
			if v, ok := cfg.Overrides[name]; ok {
				err := MergeConfig(settings, v.(map[string]interface{}))
				if err != nil {
					panic(err)
				}
			}

			err = client.CreateProject(name, groupId, cfg.Settings)
			if err != nil {
				fmt.Println(err)
				if cfg.StopOnError {
					break
				}
			}
		}
	}
}
