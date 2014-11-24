package defsecrules

import (
	"github.com/mitchellh/mapstructure"

	"github.com/rackspace/gophercloud/openstack/compute/v2/extensions/secgroups"
	"github.com/rackspace/gophercloud/pagination"
)

type DefaultRule secgroups.Rule

// DefaultRulePage is a single page of a DefaultRule collection.
type DefaultRulePage struct {
	pagination.SinglePageBase
}

// IsEmpty determines whether or not a page of default rules contains any results.
func (page DefaultRulePage) IsEmpty() (bool, error) {
	users, err := ExtractDefaultRules(page)
	if err != nil {
		return false, err
	}
	return len(users) == 0, nil
}

// ExtractDefaultRules returns a slice of DefaultRules contained in a single
// page of results.
func ExtractDefaultRules(page pagination.Page) ([]DefaultRule, error) {
	casted := page.(DefaultRulePage).Body
	var response struct {
		Rules []DefaultRule `mapstructure:"security_group_default_rules"`
	}

	err := mapstructure.Decode(casted, &response)

	return response.Rules, err
}
