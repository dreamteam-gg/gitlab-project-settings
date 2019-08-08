package main

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
		branches = append(branches, branch)
	}

	return branches, nil
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
