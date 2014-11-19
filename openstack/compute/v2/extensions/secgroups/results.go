package secgroups

import (
	"github.com/mitchellh/mapstructure"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

type SecurityGroup struct {
	ID          string
	Name        string
	Description string
	Rules       []Rule
	TenantID    string `mapstructure:"tenant_id"`
}

type Rule struct {
	ID            string
	FromPort      int     `mapstructure:"from_port"`
	ToPort        int     `mapstructure:"to_port"`
	IPProtocol    string  `mapstructure:"ip_protocol"`
	IPRange       IPRange `mapstructure:"ip_range"`
	ParentGroupID string  `mapstructure:"parent_group_id"`
	Group         Group
}

type IPRange struct {
	CIDR string
}

type Group struct {
	TenantID string `mapstructure:"tenant_id"`
	Name     string
}

// RolePage is a single page of a user Role collection.
type SecurityGroupPage struct {
	pagination.SinglePageBase
}

// IsEmpty determines whether or not a page of Security Groups contains any results.
func (page SecurityGroupPage) IsEmpty() (bool, error) {
	users, err := ExtractSecurityGroups(page)
	if err != nil {
		return false, err
	}
	return len(users) == 0, nil
}

// ExtractSecurityGroups returns a slice of SecurityGroups contained in a single page of results.
func ExtractSecurityGroups(page pagination.Page) ([]SecurityGroup, error) {
	casted := page.(SecurityGroupPage).Body
	var response struct {
		SecurityGroups []SecurityGroup `mapstructure:"security_groups"`
	}

	err := mapstructure.Decode(casted, &response)
	return response.SecurityGroups, err
}

type commonResult struct {
	gophercloud.Result
}

type CreateResult struct {
	commonResult
}

type GetResult struct {
	commonResult
}

func (r commonResult) Extract() (*SecurityGroup, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var response struct {
		SecurityGroup SecurityGroup `mapstructure:"security_group"`
	}

	err := mapstructure.Decode(r.Body, &response)

	return &response.SecurityGroup, err
}
