package main

import (
	"flag"
	"fmt"
)

var (
	flagConfigFile = flag.String("config", "./config.yml", "Path to configuration file")
	cfg            *Config
)

func CanProcessProject(name string) bool {
	if len(cfg.OnlyProject) == 0 {
		return true
	}
	for _, proj := range cfg.OnlyProject {
		if proj == name {
			return true
		}
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
	projects, err := client.GetGroupProjects(cfg.NamespaceID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found %d projects\n", len(projects))
	for _, project := range projects {
		name := project.Get("name").(string)
		id := project.Get("id").(float64)
		if !CanProcessProject(name) {
			continue
		}
		fmt.Println(name)
		for setting, cfgVal := range cfg.Settings {
			projVal := project.Get(setting)
			if projVal != cfgVal {
				fmt.Printf("\t%s = %v (%v)\n", setting, projVal, cfgVal)
			}
		}
		err = client.UpdateProject(id, cfg.Settings)
		if err != nil {
			fmt.Println(err)
			if cfg.StopOnError {
				break
			}
		}
		fmt.Println("ok")
	}
}
