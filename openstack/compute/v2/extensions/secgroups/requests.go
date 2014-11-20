package secgroups

import (
	"errors"

	"github.com/racker/perigee"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

func commonList(client *gophercloud.ServiceClient, url string) pagination.Pager {
	createPage := func(r pagination.PageResult) pagination.Page {
		return SecurityGroupPage{pagination.SinglePageBase(r)}
	}

	return pagination.NewPager(client, url, createPage)
}

func List(client *gophercloud.ServiceClient) pagination.Pager {
	return commonList(client, rootURL(client))
}

func ListByServer(client *gophercloud.ServiceClient, serverID string) pagination.Pager {
	return commonList(client, listByServerURL(client, serverID))
}

type GroupOpts struct {
	// Optional - the name of your security group. If no value provided, null
	// will be set.
	Name string `json:"name,omitempty"`

	// Optional - the description of your security group. If no value provided,
	// null will be set.
	Description string `json:"description,omitempty"`
}

type CreateOpts GroupOpts

func Create(client *gophercloud.ServiceClient, opts CreateOpts) CreateResult {
	var result CreateResult

	reqBody := struct {
		CreateOpts `json:"security_group"`
	}{opts}

	_, result.Err = perigee.Request("POST", rootURL(client), perigee.Options{
		Results:     &result.Body,
		ReqBody:     &reqBody,
		MoreHeaders: client.AuthenticatedHeaders(),
		OkCodes:     []int{200},
	})

	return result
}

type UpdateOpts GroupOpts

func Update(client *gophercloud.ServiceClient, id string, opts UpdateOpts) UpdateResult {
	var result UpdateResult

	reqBody := struct {
		UpdateOpts `json:"security_group"`
	}{opts}

	_, result.Err = perigee.Request("PUT", resourceURL(client, id), perigee.Options{
		Results:     &result.Body,
		ReqBody:     &reqBody,
		MoreHeaders: client.AuthenticatedHeaders(),
		OkCodes:     []int{200},
	})

	return result
}

func Get(client *gophercloud.ServiceClient, id string) GetResult {
	var result GetResult

	_, result.Err = perigee.Request("GET", resourceURL(client, id), perigee.Options{
		Results:     &result.Body,
		MoreHeaders: client.AuthenticatedHeaders(),
		OkCodes:     []int{200},
	})

	return result
}

func Delete(client *gophercloud.ServiceClient, id string) gophercloud.ErrResult {
	var result gophercloud.ErrResult

	_, result.Err = perigee.Request("DELETE", resourceURL(client, id), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		OkCodes:     []int{202},
	})

	return result
}

type AddRuleOpts struct {
	// Required - the ID of the group that this rule will be added to.
	ParentGroupID string `json:"parent_group_id"`

	// Required - the lower bound of the port range that will be opened.
	FromPort int `json:"from_port"`

	// Required - the upper bound of the port range that will be opened.
	ToPort int `json:"to_port"`

	// Required - the protocol type that will be allowed, e.g. TCP.
	IPProtocol string `json:"ip_protocol"`

	// ONLY required if FromGroupID is blank. This represents the IP range that
	// will be the source of network traffic to your security group. Use
	// 0.0.0.0/0 to allow all IP addresses.
	CIDR string `json:"cidr,omitempty"`

	// ONLY required if CIDR is blank. This value represents the ID of a group
	// that forwards traffic to the parent group. So, instead of accepting
	// network traffic from an entire IP range, you can instead refine the
	// inbound source by an existing security group.
	FromGroupID string `json:"group_id,omitempty"`
}

func AddRule(client *gophercloud.ServiceClient, opts AddRuleOpts) AddRuleResult {
	var result AddRuleResult

	if opts.ParentGroupID == "" {
		result.Err = errors.New("A ParentGroupID must be set")
		return result
	}
	if opts.FromPort == 0 {
		result.Err = errors.New("A FromPort must be set")
		return result
	}
	if opts.ToPort == 0 {
		result.Err = errors.New("A ToPort must be set")
		return result
	}
	if opts.IPProtocol == "" {
		result.Err = errors.New("A IPProtocol must be set")
		return result
	}
	if opts.CIDR == "" && opts.FromGroupID == "" {
		result.Err = errors.New("A CIDR or FromGroupID must be set")
		return result
	}

	reqBody := struct {
		AddRuleOpts `json:"security_group_rule"`
	}{opts}

	_, result.Err = perigee.Request("POST", rootRuleURL(client), perigee.Options{
		Results:     &result.Body,
		ReqBody:     &reqBody,
		MoreHeaders: client.AuthenticatedHeaders(),
		OkCodes:     []int{200},
	})

	return result
}

func DeleteRule(client *gophercloud.ServiceClient, id string) gophercloud.ErrResult {
	var result gophercloud.ErrResult

	_, result.Err = perigee.Request("DELETE", resourceRuleURL(client, id), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		OkCodes:     []int{202},
	})

	return result
}

func actionMap(prefix, groupName string) map[string]map[string]string {
	return map[string]map[string]string{
		prefix + "SecurityGroup": map[string]string{"name": groupName},
	}
}

func AddServerToGroup(client *gophercloud.ServiceClient, serverID, groupName string) gophercloud.ErrResult {
	var result gophercloud.ErrResult

	_, result.Err = perigee.Request("POST", serverActionURL(client, serverID), perigee.Options{
		Results:     &result.Body,
		ReqBody:     actionMap("add", groupName),
		MoreHeaders: client.AuthenticatedHeaders(),
		OkCodes:     []int{202},
	})

	return result
}

func RemoveServerFromGroup(client *gophercloud.ServiceClient, serverID, groupName string) gophercloud.ErrResult {
	var result gophercloud.ErrResult

	_, result.Err = perigee.Request("POST", serverActionURL(client, serverID), perigee.Options{
		Results:     &result.Body,
		ReqBody:     actionMap("remove", groupName),
		MoreHeaders: client.AuthenticatedHeaders(),
		OkCodes:     []int{202},
	})

	return result
}
