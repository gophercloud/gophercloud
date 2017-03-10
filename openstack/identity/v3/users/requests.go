package users

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to
// the List request
type ListOptsBuilder interface {
	ToUsersListQuery() (string, error)
}

// ListOpts allows you to query the List method.
type ListOpts struct {
	// DomainID filters the response by a domain ID.
	DomainID string `q:"domain_id"`

	// Enabled filters the response by enabled users.
	Enabled *bool `q:"enabled"`

	// IdpID filters the response by an Identity Provider ID.
	IdPID string `q:"enabled"`

	// Name filters the response by username.
	Name string `q:"name"`

	// PasswordExpiresAt filters the response based on expiring passwords.
	PasswordExpiresAt string `q:"password_expires_at"`

	// ProtocolID filters the response by protocol ID.
	ProtocolID string `q:"protocol_id"`

	// UniqueID filters the response by unique ID.
	UniqueID string `q:"unique_id"`
}

// ToUserListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToUsersListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List enumerates the Users to which the current token has access.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := baseURL(client)
	if opts != nil {
		query, err := opts.ToUsersListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return UserPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// GetOptsBuilder allows extensions to add additional parameters to
// the Get request.
type GetOptsBuilder interface {
	ToUserGetQuery() (string, error)
}

// GetOpts allows you to modify the details included in the Get request.
type GetOpts struct{}

// ToUserGetQuery formats a GetOpts into a query string.
func (opts GetOpts) ToUserGetQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// Get retrieves details on a single user, by ID.
func Get(client *gophercloud.ServiceClient, id string, opts GetOptsBuilder) (r GetResult) {
	url := userURL(client, id)
	if opts != nil {
		query, err := opts.ToUserGetQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}

	_, r.Err = client.Get(url, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to
// the Create request.
type CreateOptsBuilder interface {
	ToUserCreateMap() (map[string]interface{}, error)
}

// CreateOpts implements CreateOptsBuilder
type CreateOpts struct {
	// Name is the name of the new user.
	Name string `json:"name" required:"true"`

	// DefaultProjectID is the ID of the default project of the user.
	DefaultProjectID string `json:"default_project_id,omitempty"`

	// Description is a description of the user.
	Description string `json:"description,omitempty"`

	// DomainID is the ID of the domain the user belongs to.
	DomainID string `json:"domain_id,omitempty"`

	// Enabled sets the user status to enabled or disabled.
	Enabled *bool `json:"enabled,omitempty"`

	// Password is the password of the new user.
	Password string `json:"password,omitempty"`
}

// ToUserCreateMap formats a CreateOpts into a create request.
func (opts CreateOpts) ToUserCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "user")
}

// Create creates a new User.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToUserCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(baseURL(client), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	return
}

// Delete deletes a user.
func Delete(client *gophercloud.ServiceClient, userID string) (r DeleteResult) {
	_, r.Err = client.Delete(userURL(client, userID), nil)
	return
}
