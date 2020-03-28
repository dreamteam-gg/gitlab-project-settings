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

func (c *Client) UpdateGroup(id int, settings map[string]interface{}) error {
	err := c.UpdateGroupMembers(id, settings)
	if err != nil {
		return err
	}

	err = c.UpdateGroupVariables(id, settings)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateGroupMembers(id int, settings map[string]interface{}) error {
	var membersSettings map[string]interface{}
	if v, ok := settings["members"]; ok {
		membersSettings = v.(map[string]interface{})
	}

	existingMembers, err := c.GetGroupMembers(id)
	if err != nil {
		return err
	}

	members, err := UnmarshalMembers(membersSettings)
	if err != nil {
		return err
	}

	diff, removedDiff, changedItems, equal := computeDiff(existingMembers, members)
	if (!equal || len(changedItems.Removed) > 0) && len(membersSettings) > 0 {
		fmt.Println("\t Updating group members")
		fmt.Println(diff + removedDiff)
	}
	if !*flagDryRun && (!equal || len(changedItems.Removed) > 0) && len(membersSettings) > 0 {
		for _, m := range changedItems.Added {
			userId, err := c.GetUserIdByName(m)
			if err != nil {
				return err
			}
			accessLevel, err := AccessString(membersSettings[m].(string))
			if err != nil {
				return err
			}

			_, err = c.doFormRequest(http.MethodPost, fmt.Sprintf("groups/%d/members", int(id)), map[string]interface{}{
				"user_id":      userId,
				"access_level": accessLevel,
			})

			if err != nil {
				return err
			}
		}

		for _, m := range changedItems.Modified {
			userId, err := c.GetUserIdByName(m)
			if err != nil {
				return err
			}
			accessLevel, err := AccessString(membersSettings[m].(string))
			if err != nil {
				return err
			}

			_, err = c.doFormRequest(http.MethodPut, fmt.Sprintf("groups/%d/members/%d", int(id), userId), map[string]interface{}{
				"access_level": accessLevel,
			})
			if err != nil {
				return err
			}
		}

		for _, m := range changedItems.Removed {
			userId, err := c.GetUserIdByName(m)
			if err != nil {
				return err
			}

			_, err = c.doFormRequest(http.MethodDelete, fmt.Sprintf("groups/%d/members/%d", int(id), userId), nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Client) UpdateGroupVariables(id int, settings map[string]interface{}) error {
	if v, ok := settings["variables"]; ok {
		err := c.UpdateVariables("groups", id, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func UnmarshalMembers(settings map[string]interface{}) (map[string]interface{}, error) {
	members := make(map[string]interface{})

	for n, a := range settings {
		access, err := AccessString(a.(string))
		if err != nil {
			return nil, err
		}
		members[n] = access
	}

	return members, nil
}

func (c *Client) GetGroupMembers(id int) (map[string]interface{}, error) {
	members := make(map[string]interface{})

	resp, err := c.doRequest(http.MethodGet, fmt.Sprintf("groups/%v/members", id), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("Error getting group members. Return code not 2XX: %s", resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	mmbr := []map[string]interface{}{}
	if err := json.Unmarshal(b, &mmbr); err != nil {
		return nil, err
	}

	for _, m := range mmbr {
		member := m
		members[member["username"].(string)] = FloatToAccess(m)
	}

	return members, nil
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
