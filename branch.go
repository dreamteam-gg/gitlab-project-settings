package main

import (
	"encoding/json"
)

type BranchAccess int

//go:generate enumer -type=BranchAccess -yaml
const (
	NoAccess   BranchAccess = 0
	Developer  BranchAccess = 30
	Maintainer BranchAccess = 40
	Admin      BranchAccess = 60
)

// MarshalJSON implements the json.Marshaler interface for BranchAccess
// when we marshal to json for GitlabApi we need int representation
func (i BranchAccess) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(i))
}

type Branch map[string]interface{}

func (p Branch) Get(key string) interface{} {
	return p[key]
}

func UnmarshalBranches(settings map[string]interface{}) ([]map[string]interface{}, error) {
	branches := make([]map[string]interface{}, 0)

	for n, b := range settings {
		branch := make(map[string]interface{})
		access := b.(map[string]interface{})
		branch["name"] = n
		if v, err := BranchAccessString(access["push_access_level"].(string)); err == nil {
			branch["push_access_level"] = v
		} else {
			return nil, err
		}
		if v, err := BranchAccessString(access["merge_access_level"].(string)); err == nil {
			branch["merge_access_level"] = v
		} else {
			return nil, err
		}
		if v, err := BranchAccessString(access["unprotect_access_level"].(string)); err == nil {
			branch["unprotect_access_level"] = v
		} else {
			return nil, err
		}
		branches = append(branches, branch)
	}

	return branches, nil
}

func FloatToBranchAccess(i interface{}) BranchAccess {
	return BranchAccess(int(i.(map[string]interface{})["access_level"].(float64)))
}

func FinBranchInList(name string, branches []map[string]interface{}) map[string]interface{} {
	for _, b := range branches {
		if b["name"].(string) == name {
			return b
		}
	}
	// return empty if not found
	return map[string]interface{}{}
}
