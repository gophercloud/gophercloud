package secgroups

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

func commonList(client *gophercloud.ServiceClient, url string) pagination.Pager {
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return SecurityGroupPage{pagination.SinglePageBase(r)}
	})
}

// List will return a collection of all the security groups for a particular
// tenant.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	return commonList(client, rootURL(client))
}

// ListByServer will return a collection of all the security groups which are
// associated with a particular server.
func ListByServer(client *gophercloud.ServiceClient, serverID string) pagination.Pager {
	return commonList(client, listByServerURL(client, serverID))
}

// CreateOpts is the struct responsible for creating a security group.
type CreateOpts struct {
	// the name of your security group.
	Name string `json:"name" required:"true"`
	// the description of your security group.
	Description string `json:"description,omitempty"`
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToSecGroupCreateMap() (map[string]any, error)
}

// ToSecGroupCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToSecGroupCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "security_group")
}

// Create will create a new security group.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToSecGroupCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, rootURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOpts is the struct responsible for updating an existing security group.
type UpdateOpts struct {
	// the name of your security group.
	Name string `json:"name,omitempty"`
	// the description of your security group.
	Description *string `json:"description,omitempty"`
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToSecGroupUpdateMap() (map[string]any, error)
}

// ToSecGroupUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToSecGroupUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "security_group")
}

// Update will modify the mutable properties of a security group, notably its
// name and description.
func Update(ctx context.Context, client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToSecGroupUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(ctx, resourceURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get will return details for a particular security group.
func Get(ctx context.Context, client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(ctx, resourceURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete will permanently delete a security group from the project.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(ctx, resourceURL(client, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateRuleOpts represents the configuration for adding a new rule to an
// existing security group.
type CreateRuleOpts struct {
	// ID is the ID of the group that this rule will be added to.
	ParentGroupID string `json:"parent_group_id" required:"true"`

	// FromPort is the lower bound of the port range that will be opened.
	// Use -1 to allow all ICMP traffic.
	FromPort int `json:"from_port"`

	// ToPort is the upper bound of the port range that will be opened.
	// Use -1 to allow all ICMP traffic.
	ToPort int `json:"to_port"`

	// IPProtocol the protocol type that will be allowed, e.g. TCP.
	IPProtocol string `json:"ip_protocol" required:"true"`

	// CIDR is the network CIDR to allow traffic from.
	// This is ONLY required if FromGroupID is blank. This represents the IP
	// range that will be the source of network traffic to your security group.
	// Use 0.0.0.0/0 to allow all IP addresses.
	CIDR string `json:"cidr,omitempty" or:"FromGroupID"`

	// FromGroupID represents another security group to allow access.
	// This is ONLY required if CIDR is blank. This value represents the ID of a
	// group that forwards traffic to the parent group. So, instead of accepting
	// network traffic from an entire IP range, you can instead refine the
	// inbound source by an existing security group.
	FromGroupID string `json:"group_id,omitempty" or:"CIDR"`
}

// CreateRuleOptsBuilder allows extensions to add additional parameters to the
// CreateRule request.
type CreateRuleOptsBuilder interface {
	ToRuleCreateMap() (map[string]any, error)
}

// ToRuleCreateMap builds a request body from CreateRuleOpts.
func (opts CreateRuleOpts) ToRuleCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "security_group_rule")
}

// CreateRule will add a new rule to an existing security group (whose ID is
// specified in CreateRuleOpts). You have the option of controlling inbound
// traffic from either an IP range (CIDR) or from another security group.
func CreateRule(ctx context.Context, client *gophercloud.ServiceClient, opts CreateRuleOptsBuilder) (r CreateRuleResult) {
	b, err := opts.ToRuleCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, rootRuleURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteRule will permanently delete a rule from a security group.
func DeleteRule(ctx context.Context, client *gophercloud.ServiceClient, id string) (r DeleteRuleResult) {
	resp, err := client.Delete(ctx, resourceRuleURL(client, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

func actionMap(prefix, groupName string) map[string]map[string]string {
	return map[string]map[string]string{
		prefix + "SecurityGroup": {"name": groupName},
	}
}

// AddServer will associate a server and a security group, enforcing the
// rules of the group on the server.
func AddServer(ctx context.Context, client *gophercloud.ServiceClient, serverID, groupName string) (r AddServerResult) {
	resp, err := client.Post(ctx, serverActionURL(client, serverID), actionMap("add", groupName), nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// RemoveServer will disassociate a server from a security group.
func RemoveServer(ctx context.Context, client *gophercloud.ServiceClient, serverID, groupName string) (r RemoveServerResult) {
	resp, err := client.Post(ctx, serverActionURL(client, serverID), actionMap("remove", groupName), nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
