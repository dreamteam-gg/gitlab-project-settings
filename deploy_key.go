package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (c *Client) GetDeployKeyIdByName(name string) (int, error) {
	resp, err := c.doRequest(http.MethodGet, "deploy_keys", nil)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return 0, fmt.Errorf("Error listing deploy keys. Return code not 2XX: %s", resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var keys []map[string]interface{}
	if err := json.Unmarshal(b, &keys); err != nil {
		return 0, err
	}

	for _, k := range keys {
		if k["title"].(string) == name {
			return int(k["id"].(float64)), nil
		}
	}

	return 0, fmt.Errorf("No key with name %s found", name)
}
