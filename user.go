package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type User map[string]interface{}

func (p User) Get(key string) interface{} {
	return p[key]
}

func (c *Client) GetUserIdByName(name string) (int, error) {
	resp, err := c.doRequest(http.MethodGet, fmt.Sprintf("users?username=%s", name), nil)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return 0, fmt.Errorf("Error searching for user %s. Return code not 2XX: %s", name, resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	usr := []*User{}
	if err := json.Unmarshal(b, &usr); err != nil {
		return 0, err
	}

	if l := len(usr); l != 1 {
		return 0, fmt.Errorf("Found %d %s users, need one", l, name)
	}

	return int(usr[0].Get("id").(float64)), nil
}

func (c *Client) GetUserNameById(id int) (string, error) {
	resp, err := c.doRequest(http.MethodGet, fmt.Sprintf("users/%d", id), nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", fmt.Errorf("Error searching for user %d. Return code not 2XX: %s", id, resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	usr := User{}
	if err := json.Unmarshal(b, &usr); err != nil {
		return "", err
	}

	return usr.Get("username").(string), nil
}
