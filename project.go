package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Project map[string]interface{}

func (p Project) Get(key string) interface{} {
	return p[key]
}

// ref: https://docs.gitlab.com/ee/api/groups.html#list-a-group-s-projects
func (c *Client) GetGroupProjects(id int) ([]*Project, error) {
	var projects []*Project
	pagedURL, err := url.Parse(fmt.Sprintf("groups/%d/projects", id))
	if err != nil {
		return nil, err
	}
	values := url.Values{
		"page":     []string{"1"},
		"per_page": []string{"50"},
	}
	pagedURL.RawQuery = values.Encode()
	for {
		resp, err := c.doRequest(http.MethodGet, pagedURL.String(), nil)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			return nil, fmt.Errorf("return code not 2XX: %s", resp.Status)
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		prs := []*Project{}
		if err := json.Unmarshal(b, &prs); err != nil {
			return nil, err
		}
		for _, project := range prs {
			projects = append(projects, project)
		}
		// accumulate(obj)
		page := resp.Header.Get("x-next-page")
		if len(page) == 0 {
			break
		}
		q := pagedURL.Query()
		q.Set("page", page)
		pagedURL.RawQuery = q.Encode()
	}
	return projects, nil
}

// ref: https://docs.gitlab.com/ee/api/projects.html#edit-project
func (c *Client) UpdateProject(project *Project, settings map[string]interface{}) error {
	id := int(project.Get("id").(float64))
	name := project.Get("name").(string)
	var projectSettings map[string]interface{}

	if v, ok := settings["project"]; ok {
		projectSettings = v.(map[string]interface{})
	}

	fmt.Println(formatter.Bold(name).String())
	diff, equal := computeDiff((map[string]interface{})(*project), projectSettings)
	if !equal {
		fmt.Println(diff)
	}
	if !*flagDryRun && !equal {
		_, err := c.doFormRequest(http.MethodPut, fmt.Sprintf("projects/%v", id), projectSettings)
		if err != nil {
			return err
		}
	}

	err := c.UpdateProjectApprovals(project, settings)
	if err != nil {
		return err
	}

	err = c.UpdateProjectProtectedBranches(project, settings)
	if err != nil {
		return err
	}

	err = c.UpdateProjectServices(project, settings)
	if err != nil {
		return err
	}

	err = c.UpdateProjectWebHooks(project, settings)
	if err != nil {
		return err
	}

	err = c.UpdateProjectDeployKeys(project, settings)
	if err != nil {
		return err
	}

	fmt.Println("ok")
	return nil
}

// ref: https://docs.gitlab.com/ee/api/merge_request_approvals.html#change-configuration
func (c *Client) UpdateProjectApprovals(project *Project, settings map[string]interface{}) error {
	id := int(project.Get("id").(float64))

	var approvalSettings map[string]interface{}
	if v, ok := settings["approvals"]; ok {
		approvalSettings = v.(map[string]interface{})
	}

	var approversSettings map[string]interface{}
	if v, ok := settings["approvers"]; ok {
		approversSettings = v.(map[string]interface{})
	}

	existingApprovals, err := c.GetProjectApprovals(project)
	if err != nil {
		return err
	}

	diff, equal := computeDiff(existingApprovals, approvalSettings)
	if !equal {
		fmt.Println("\t Updating approvals")
		fmt.Println(diff)
	}
	if !*flagDryRun && !equal {
		_, err := c.doFormRequest(http.MethodPost, fmt.Sprintf("projects/%d/approvals", int(id)), approvalSettings)
		if err != nil {
			return err
		}
	}

	approvalsNames, err := c.ConvertApprovalIdsToNames(existingApprovals)
	if err != nil {
		return err
	}

	approvalIds, err := c.ConvertApprovalNamesToIds(approversSettings)
	if err != nil {
		return err
	}

	diff, equal = computeDiff(approvalsNames, approversSettings)
	if !equal {
		fmt.Println("\t Updating approvers")
		fmt.Println(diff)
	}

	if !*flagDryRun && !equal {
		_, err := c.doFormRequest(http.MethodPut, fmt.Sprintf("projects/%d/approvers", int(id)), approvalIds)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) UpdateProjectProtectedBranches(project *Project, settings map[string]interface{}) error {
	id := int(project.Get("id").(float64))
	existingBranches, err := c.GetProjectProtectedBranches(project)
	if err != nil {
		return err
	}

	var settingBranches map[string]interface{}
	if v, ok := settings["protected_branches"]; ok {
		settingBranches = v.(map[string]interface{})
	} else {
		return nil
	}

	protectedBranches, err := UnmarshalBranches(settingBranches)
	if err != nil {
		return err
	}

	for _, b := range protectedBranches {
		name := b["name"].(string)
		diff, equal := computeDiff(FinBranchInList(name, existingBranches), b)
		if !equal {
			fmt.Printf("\t Updating branch '%s'\n", name)
			fmt.Println(diff)
		}

		if !*flagDryRun && !equal {
			// we need to unprotect branch before protecting, otherwise we get 409 "Protected branch '*' already exists"
			c.doFormRequest(http.MethodDelete, fmt.Sprintf("projects/%d/protected_branches/%s", id, name), b)

			_, err = c.doFormRequest(http.MethodPost, fmt.Sprintf("projects/%d/protected_branches", id), b)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Client) UpdateProjectServices(project *Project, settings map[string]interface{}) error {
	id := int(project.Get("id").(float64))

	var services map[string]interface{}
	if v, ok := settings["services"]; ok {
		services = v.(map[string]interface{})
	} else {
		return nil
	}

	for n, s := range services {
		existingSettings, err := c.GetProjectService(project, n)
		if err != nil {
			return err
		}
		newSettings := s.(map[string]interface{})

		diff, equal := computeDiff(existingSettings, newSettings)
		if !equal {
			fmt.Printf("\t Updating service '%s'\n", n)
			fmt.Println(diff)
		}

		if !*flagDryRun && !equal {
			_, err = c.doFormRequest(http.MethodPut, fmt.Sprintf("projects/%d/services/%s", id, n), newSettings)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Client) UpdateProjectWebHooks(project *Project, settings map[string]interface{}) error {
	id := int(project.Get("id").(float64))

	var hooks map[string]interface{}
	if v, ok := settings["webhooks"]; ok {
		hooks = v.(map[string]interface{})
	} else {
		return nil
	}

	for u, h := range hooks {
		hookId, existingSettings, err := c.GetProjectWebHook(project, u)
		if err != nil {
			return err
		}

		diff, equal := computeDiff(existingSettings, h.(map[string]interface{}))
		if !equal {
			fmt.Printf("\t Updating hook '%s'\n", u)
			fmt.Println(diff)
		}

		h.(map[string]interface{})["url"] = u

		if !*flagDryRun && id == 0 {
			// create hook if it does not exist
			if !*flagDryRun {
				_, err = c.doFormRequest(http.MethodPost, fmt.Sprintf("projects/%d/hooks", id), h.(map[string]interface{}))
				if err != nil {
					return err
				}
			}
		} else if !*flagDryRun && !equal {
			// update existing hook
			if !*flagDryRun {
				_, err = c.doFormRequest(http.MethodPut, fmt.Sprintf("projects/%d/hooks/%d", id, hookId), h.(map[string]interface{}))
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (c *Client) UpdateProjectDeployKeys(project *Project, settings map[string]interface{}) error {
	id := int(project.Get("id").(float64))

	keys := make(map[string]interface{})
	if v, ok := settings["deploy_keys"]; ok {
		keys["deploy_keys"] = v.([]interface{})
	} else {
		return nil
	}

	existingKeys, err := c.GetProjectDeployKeys(project)
	if err != nil {
		return err
	}

	newKeyIds := make([]int, len(keys["deploy_keys"].([]interface{})))
	for i, k := range keys["deploy_keys"].([]interface{}) {
		kId, err := c.GetDeployKeyIdByName(k.(string))
		if err != nil {
			return err
		}
		newKeyIds[i] = kId
	}

	diff, equal := computeDiff(existingKeys, keys)
	if !equal {
		fmt.Println("\tUpdating deploy keys")
		fmt.Println(diff)
	}

	if !*flagDryRun && !equal {
		for _, k := range newKeyIds {
			_, err = c.doFormRequest(http.MethodPost, fmt.Sprintf("projects/%d/deploy_keys/%d/enable", id, k), nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Client) GetProjectApprovals(project *Project) (map[string]interface{}, error) {
	var approvals map[string]interface{}
	id := int(project.Get("id").(float64))

	resp, err := c.doRequest(http.MethodGet, fmt.Sprintf("projects/%v/approvals", id), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("Error getting approval settings. Return code not 2XX: %s", resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &approvals); err != nil {
		return nil, err
	}

	// fix unmarshalling to int
	approvals["approvals_before_merge"] = int(approvals["approvals_before_merge"].(float64))

	// get approvers
	g := []int{}
	for _, v := range approvals["approver_groups"].([]interface{}) {
		g = append(g, int(v.(map[string]interface{})["group"].(map[string]interface{})["id"].(float64)))
	}
	approvals["approver_group_ids"] = g

	u := []int{}
	for _, v := range approvals["approvers"].([]interface{}) {
		u = append(u, int(v.(map[string]interface{})["user"].(map[string]interface{})["id"].(float64)))
	}
	approvals["approver_ids"] = u

	return approvals, nil
}

func (c *Client) ConvertApprovalNamesToIds(settings map[string]interface{}) (map[string]interface{}, error) {
	converted := make(map[string]interface{})
	groups := make([]int, 0)
	users := make([]int, 0)

	if v, ok := settings["approver_ids"]; ok {
		for _, name := range v.([]interface{}) {
			user, err := c.GetUserIdByName(name.(string))
			if err != nil {
				return nil, err
			}
			users = append(users, user)
		}
	}

	if v, ok := settings["approver_group_ids"]; ok {
		for _, name := range v.([]interface{}) {
			group, err := c.GetGroupIdByName(name.(string))
			if err != nil {
				return nil, err
			}
			groups = append(groups, group)
		}
	}

	converted["approver_group_ids"] = groups
	converted["approver_ids"] = users

	return converted, nil
}

func (c *Client) ConvertApprovalIdsToNames(settings map[string]interface{}) (map[string]interface{}, error) {
	converted := make(map[string]interface{})
	groups := make([]interface{}, 0)
	users := make([]interface{}, 0)

	for _, id := range settings["approver_group_ids"].([]int) {
		group, err := c.GetGroupNameById(id)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	for _, id := range settings["approver_ids"].([]int) {
		user, err := c.GetUserNameById(id)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	converted["approver_group_ids"] = groups
	converted["approver_ids"] = users

	return converted, nil
}

func (c *Client) GetProjectProtectedBranches(project *Project) ([]map[string]interface{}, error) {
	var branches []map[string]interface{}
	id := int(project.Get("id").(float64))

	resp, err := c.doRequest(http.MethodGet, fmt.Sprintf("projects/%v/protected_branches", id), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("Error getting protected branch settings. Return code not 2XX: %s", resp.Status)
	}

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	brnch := []map[string]interface{}{}
	if err := json.Unmarshal(r, &brnch); err != nil {
		return nil, err
	}

	for _, b := range brnch {
		branch := make(map[string]interface{})
		branch["name"] = b["name"].(string)
		if v := b["push_access_levels"].([]interface{}); len(v) > 0 {
			branch["push_access_level"] = FloatToBranchAccess(v[0])
		}
		if v := b["merge_access_levels"].([]interface{}); len(v) > 0 {
			branch["merge_access_level"] = FloatToBranchAccess(v[0])
		}
		if v := b["unprotect_access_levels"].([]interface{}); len(v) > 0 {
			branch["unprotect_access_level"] = FloatToBranchAccess(v[0])
		}

		branches = append(branches, branch)
	}

	return branches, nil
}

// ref https://docs.gitlab.com/ee/api/services.html#get-slack-service-settings
func (c *Client) GetProjectService(project *Project, service string) (map[string]interface{}, error) {
	id := int(project.Get("id").(float64))

	resp, err := c.doRequest(http.MethodGet, fmt.Sprintf("/projects/%d/services/%s", id, service), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("Error getting %s settings. Return code not 2XX: %s", service, resp.Status)
	}

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	srv := map[string]interface{}{}
	if err := json.Unmarshal(r, &srv); err != nil {
		return nil, err
	}

	err = MergeConfig(srv, srv["properties"].(map[string]interface{}))
	if err != nil {
		return nil, err
	}

	return srv, nil
}

// ref https://docs.gitlab.com/ee/api/projects.html#get-project-hook
func (c *Client) GetProjectWebHook(project *Project, url string) (int, map[string]interface{}, error) {
	id := int(project.Get("id").(float64))

	resp, err := c.doRequest(http.MethodGet, fmt.Sprintf("/projects/%d/hooks", id), nil)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return 0, nil, fmt.Errorf("Error getting project %d hooks. Return code not 2XX: %s", id, resp.Status)
	}

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	hooks := []map[string]interface{}{}
	if err := json.Unmarshal(r, &hooks); err != nil {
		return 0, nil, err
	}

	for _, h := range hooks {
		if h["url"].(string) == url {
			return int(h["id"].(float64)), h, nil
		}
	}

	return 0, nil, nil
}

func (c *Client) GetProjectDeployKeys(project *Project) (map[string]interface{}, error) {
	id := int(project.Get("id").(float64))
	keys := make(map[string]interface{})

	resp, err := c.doRequest(http.MethodGet, fmt.Sprintf("/projects/%d/deploy_keys", id), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("Error getting project %d hooks. Return code not 2XX: %s", id, resp.Status)
	}

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	k := []map[string]interface{}{}
	if err := json.Unmarshal(r, &k); err != nil {
		return nil, err
	}

	keyNames := make([]interface{}, len(k))
	for i, n := range k {
		keyNames[i] = n["title"].(string)
	}

	keys["deploy_keys"] = keyNames

	return keys, nil
}

// ref: https://docs.gitlab.com/ee/api/projects.html#create-project
func (c *Client) CreateProject(name string, namespace int, settings map[string]interface{}) error {
	projectSettings := settings["project"].(map[string]interface{})
	projectSettings["name"] = name
	projectSettings["namespace_id"] = namespace

	fmt.Printf("Will create missing project '%s'\n", name)
	if *flagDryRun {
		fmt.Println()
		return nil
	}

	resp, err := c.doFormRequest(http.MethodPost, "projects", projectSettings)
	if err != nil {
		return err
	}

	p := Project{}
	if err := json.Unmarshal(resp, &p); err != nil {
		return err
	}

	err = c.UpdateProject(&p, settings)
	if err != nil {
		return err
	}

	fmt.Println("ok")
	return nil
}
