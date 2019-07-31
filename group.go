package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Group map[string]interface{}

func (p Group) Get(key string) interface{} {
	return p[key]
}

func (c *Client) GetGroupIdByName(name string) (int, error) {
	resp, err := c.doRequest(http.MethodGet, fmt.Sprintf("groups/%s", url.QueryEscape(name)), nil)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return 0, fmt.Errorf("Error searching for group %s. Return code not 2XX: %s", name, resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	grp := Group{}
	if err := json.Unmarshal(b, &grp); err != nil {
		return 0, err
	}

	return int(grp.Get("id").(float64)), nil
}

func (c *Client) GetGroupNameById(id int) (string, error) {
	resp, err := c.doRequest(http.MethodGet, fmt.Sprintf("groups/%d", id), nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", fmt.Errorf("Error searching for group %d. Return code not 2XX: %s", id, resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	grp := Group{}
	if err := json.Unmarshal(b, &grp); err != nil {
		return "", err
	}

	return grp.Get("full_path").(string), nil
}

func normalizeGroupName(name string) string {
	return url.QueryEscape(name)
}
