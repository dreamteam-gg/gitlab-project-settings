package main

import "encoding/json"

type Access int

//go:generate enumer -type=Access -yaml
const (
	NoAccess   Access = 0
	Guest      Access = 10
	Reporter   Access = 20
	Developer  Access = 30
	Maintainer Access = 40
	Owner      Access = 50
	Admin      Access = 60
)

// MarshalJSON implements the json.Marshaler interface for Access
// when we marshal to json for GitlabApi we need int representation
func (i Access) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(i))
}

func FloatToAccess(i interface{}) Access {
	return Access(int(i.(map[string]interface{})["access_level"].(float64)))
}
