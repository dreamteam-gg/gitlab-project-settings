package main

type Branch map[string]interface{}

func (p Branch) Get(key string) interface{} {
	return p[key]
}

func (c *Client) UnmarshalBranches(settings map[string]interface{}) ([]map[string]interface{}, error) {
	branches := make([]map[string]interface{}, 0)

	for n, b := range settings {
		branch := make(map[string]interface{})
		access := b.(map[string]interface{})

		branch["name"] = n
		if v, err := AccessString(access["push_access_level"].(string)); err == nil {
			branch["push_access_level"] = v
		} else {
			return nil, err
		}
		if v, err := AccessString(access["merge_access_level"].(string)); err == nil {
			branch["merge_access_level"] = v
		} else {
			return nil, err
		}
		if v, err := AccessString(access["unprotect_access_level"].(string)); err == nil {
			branch["unprotect_access_level"] = v
		} else {
			return nil, err
		}

		if v, ok := access["allowed_to_push"]; ok {
			branch["allowed_to_push"] = v
		}
		if v, ok := access["allowed_to_merge"]; ok {
			branch["allowed_to_merge"] = v
		}
		if v, ok := access["allowed_to_unprotect"]; ok {
			branch["allowed_to_unprotect"] = v
		}

		branches = append(branches, branch)
	}

	return branches, nil
}

func (c *Client) BranchAllowedToIds(branch map[string]interface{}) (map[string]interface{}, error) {

	push, err := c.UnmarshallAllowed(branch["allowed_to_push"])
	if err != nil {
		return nil, err
	}
	branch["allowed_to_push"] = push
	merge, err := c.UnmarshallAllowed(branch["allowed_to_merge"])
	if err != nil {
		return nil, err
	}
	branch["allowed_to_merge"] = merge
	unpro, err := c.UnmarshallAllowed(branch["allowed_to_unprotect"])
	branch["allowed_to_unprotect"] = unpro

	return branch, nil
}

func (c *Client) UnmarshallAllowed(a interface{}) ([]map[string]interface{}, error) {
	var allowed []map[string]interface{}
	if a == nil {
		return nil, nil
	}
	for _, m := range a.([]interface{}) {
		for k, v := range InterfaceMapToStringMap(m.(map[interface{}]interface{})) {
			var i int
			var err error
			switch k {
			case "user_id":
				i, err = c.GetUserIdByName(v)
				if err != nil {
					return nil, err
				}
			case "group_id":
				i, err = c.GetGroupIdByName(v)
				if err != nil {
					return nil, err
				}
			}
			allowed = append(allowed, map[string]interface{}{
				k: i,
			})
		}
	}

	return allowed, nil
}

func FindBranchInList(name string, branches []map[string]interface{}) map[string]interface{} {
	for _, b := range branches {
		if b["name"].(string) == name {
			return b
		}
	}
	// return empty if not found
	return map[string]interface{}{}
}
