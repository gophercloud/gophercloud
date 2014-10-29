package users

import (
	"errors"

	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

func List(client *gophercloud.ServiceClient) pagination.Pager {
	createPage := func(r pagination.PageResult) pagination.Page {
		return UserPage{pagination.SinglePageBase(r)}
	}

	return pagination.NewPager(client, rootURL(client), createPage)
}

// EnabledState represents whether the user is enabled or not.
type EnabledState *bool

// Useful variables to use when creating or updating users.
var (
	iTrue  = true
	iFalse = false

	Enabled  EnabledState = &iTrue
	Disabled EnabledState = &iFalse
)

// CreateOpts represents the options needed when creating new users.
type CreateOpts struct {
	// Either a name or username is required. When provided, the value must be
	// unique or a 409 conflict error will be returned. If you provide a name but
	// omit a username, the latter will be set to the former; and vice versa.
	Name, Username string

	// The ID of the tenant to which you want to assign this user.
	TenantID string

	// Indicates whether this user is enabled or not.
	Enabled EnabledState

	// The email address of this user.
	Email string
}

// CreateOptsBuilder describes struct types that can be accepted by the Create call.
type CreateOptsBuilder interface {
	ToUserCreateMap() (map[string]interface{}, error)
}

// ToUserCreateMap assembles a request body based on the contents of a CreateOpts.
func (opts CreateOpts) ToUserCreateMap() (map[string]interface{}, error) {
	m := make(map[string]interface{})

	if opts.Name == "" && opts.Username == "" {
		return m, errors.New("Either a Name or Username must be provided")
	}

	if opts.Name != "" {
		m["name"] = opts.Name
	}
	if opts.Username != "" {
		m["username"] = opts.Username
	}
	if opts.Enabled != nil {
		m["enabled"] = &opts.Enabled
	}
	if opts.Email != "" {
		m["email"] = opts.Email
	}

	return m, nil
}

// Create is the operation responsible for creating new users.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) CreateResult {
	var res CreateResult

	reqBody, err := opts.ToUserCreateMap()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = perigee.Request("POST", rootURL(client), perigee.Options{
		Results:     &res.Body,
		ReqBody:     reqBody,
		MoreHeaders: client.AuthenticatedHeaders(),
		OkCodes:     []int{200, 201},
	})

	return res
}
