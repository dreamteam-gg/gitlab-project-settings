package main

type Project map[string]interface{}

func (p Project) Get(key string) interface{} {
	return p[key]
}
