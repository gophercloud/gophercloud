package qos

import "github.com/gophercloud/gophercloud"

// The result of listing the qos rule types
type ListRuleTypesResult struct {
	gophercloud.Result
}

func (l ListRuleTypesResult) Extract() (out []string, err error) {
	// the response from the server is in this format, but it's not an especially useful format
	// so we will return the rule types, in a slice of strings
	var list struct {
		RuleTypes []struct {
			Type string `json:"type"`
		} `json:"rule_types"`
	}

	// take the JSON and put it into the list struct
	if err = l.ExtractInto(&list); err != nil {
		return
	}

	// now make the out slice
	out = make([]string, len(list.RuleTypes))
	// and copy the result to the return value
	for i, t := range list.RuleTypes {
		out[i] = t.Type
	}

	return
}
