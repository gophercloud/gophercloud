package defsecrules

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// List will return a collection of default rules.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, rootURL(client), func(r pagination.PageResult) pagination.Page {
		return DefaultRulePage{pagination.SinglePageBase(r)}
	})
}

// CreateOpts represents the configuration for adding a new default rule.
type CreateOpts struct {
	// The lower bound of the port range that will be opened.
	FromPort int `json:"from_port" required:"true"`
	// The upper bound of the port range that will be opened.
	ToPort int `json:"to_port" required:"true"`
	// The protocol type that will be allowed, e.g. TCP.
	IPProtocol string `json:"ip_protocol" required:"true"`
	// ONLY required if FromGroupID is blank. This represents the IP range that
	// will be the source of network traffic to your security group. Use
	// 0.0.0.0/0 to allow all IP addresses.
	CIDR string `json:"cidr,omitempty"`
}

// CreateOptsBuilder builds the create rule options into a serializable format.
type CreateOptsBuilder interface {
	ToRuleCreateMap() (map[string]interface{}, error)
}

// ToRuleCreateMap builds the create rule options into a serializable format.
func (opts CreateOpts) ToRuleCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "security_group_default_rule")
}

// Create is the operation responsible for creating a new default rule.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) CreateResult {
	var r CreateResult
	b, err := opts.ToRuleCreateMap()
	if err != nil {
		r.Err = err
		return r
	}
	_, r.Err = client.Post(rootURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return r
}

// Get will return details for a particular default rule.
func Get(client *gophercloud.ServiceClient, id string) GetResult {
	var r GetResult
	_, r.Err = client.Get(resourceURL(client, id), &r.Body, nil)
	return r
}

// Delete will permanently delete a default rule from the project.
func Delete(client *gophercloud.ServiceClient, id string) gophercloud.ErrResult {
	var r gophercloud.ErrResult
	_, r.Err = client.Delete(resourceURL(client, id), nil)
	return r
}
