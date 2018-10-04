package trunks

import (
	"github.com/gophercloud/gophercloud"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToTrunkCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents the attributes used when creating a new trunk.
type CreateOpts struct {
	TenantID     string    `json:"tenant_id,omitempty"`
	ProjectID    string    `json:"project_id,omitempty"`
	PortID       string    `json:"port_id" required:"true"`
	Name         string    `json:"name,omitempty"`
	Description  string    `json:"description,omitempty"`
	AdminStateUp *bool     `json:"admin_state_up,omitempty"`
	Subports     []Subport `json:"sub_ports"`
}

// ToTrunkCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToTrunkCreateMap() (map[string]interface{}, error) {
	if opts.Subports == nil {
		opts.Subports = []Subport{}
	}
	return gophercloud.BuildRequestBody(opts, "trunk")
}

func Create(c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	body, err := opts.ToTrunkCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = c.Post(createURL(c), body, &r.Body, nil)
	return
}

// Delete accepts a unique ID and deletes the trunk associated with it.
func Delete(c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(deleteURL(c, id), nil)
	return
}
