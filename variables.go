package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
)

func (c *Client) UpdateVariables(entity string, id int, settings interface{}) error {
	existingVariables, err := c.GetVariables(entity, id)
	if err != nil {
		return err
	}

	for _, v := range settings.([]interface{}) {
		variable := InterfaceMapToInterfaceMap(v.(map[interface{}]interface{}))

		// if variable is masked mask its value
		oldMasking := viper.GetStringSlice("mask")
		if variable["masked"].(bool) {
			viper.Set("mask", append(viper.GetStringSlice("mask"), "value"))
		}

		if v, ok := existingVariables[variable["key"].(string)]; ok {
			diff, _, _, equal := computeDiff(v.(map[string]interface{}), variable)
			if !equal {
				fmt.Printf("\t Updating variable '%s'\n", variable["key"].(string))
				fmt.Println(diff)
			}
			if !*flagDryRun && !equal {
				_, err = c.doFormRequest(http.MethodPut, fmt.Sprintf("%s/%d/variables/%s", entity, id, variable["key"].(string)), variable)
				if err != nil {
					return fmt.Errorf("error updating variable: %v", err)
				}
			}
		} else {
			diff, _, _, equal := computeDiff(map[string]interface{}{"key": variable["key"].(string)}, variable)
			if !equal {
				fmt.Printf("\t Creating variable '%s'\n", variable["key"].(string))
				fmt.Println(diff)
			}
			if !*flagDryRun && !equal {
				_, err = c.doFormRequest(http.MethodPost, fmt.Sprintf("%s/%d/variables", entity, id), variable)
				if err != nil {
					return fmt.Errorf("error creating variable: %v", err)
				}
			}
		}

		// restore old masking config
		viper.Set("mask", oldMasking)
	}

	return nil
}

func (c *Client) GetVariables(entity string, id int) (map[string]interface{}, error) {
	resp, err := c.doRequest(http.MethodGet, fmt.Sprintf("%s/%d/variables", entity, id), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("Error getting project %d pipeline schedules. Return code not 2XX: %s", id, resp.Status)
	}

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	variables := []map[string]interface{}{}
	if err := json.Unmarshal(r, &variables); err != nil {
		return nil, err
	}

	variableMap := make(map[string]interface{}, len(variables))
	for _, v := range variables {
		variableMap[v["key"].(string)] = v
	}

	return variableMap, nil
}
