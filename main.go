package main

import (
	"flag"
	"fmt"
)

var (
	flagConfigFile = flag.String("config", "./config.yml", "Path to configuration file")
	flagDryRun     = flag.Bool("dry-run", false, "Dry run mode")
	cfg            *Config
)

func InArray(k string, arr []string) bool {
	for _, val := range arr {
		if val == k {
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
	projects, err := client.GetGroupProjects(cfg.GroupID)
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
			if !IsEqual(projVal, cfgVal) {
				fmt.Printf("\t%s = %v (%v)\n", setting, projVal, cfgVal)
			}
		}
		if *flagDryRun {
			fmt.Println()
			continue
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
