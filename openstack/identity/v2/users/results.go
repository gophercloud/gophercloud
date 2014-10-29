package users

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud/pagination"
)

// User represents a user resource that exists on the API.
type User struct {
	// The UUID for this user.
	ID string

	// The human name for this user.
	Name string

	// The username for this user.
	Username string

	// Indicates whether the user is enabled (true) or disabled (false).
	Enabled bool

	// The email address for this user.
	Email string

	// The ID of the tenant to which this user belongs.
	TenantID string `mapstructure:"tenant_id"`
}

// UserPage is a single page of a User collection.
type UserPage struct {
	pagination.SinglePageBase
}

// IsEmpty determines whether or not a page of Tenants contains any results.
func (page UserPage) IsEmpty() (bool, error) {
	users, err := ExtractUsers(page)
	if err != nil {
		return false, err
	}
	return len(users) == 0, nil
}

// ExtractUsers returns a slice of Tenants contained in a single page of results.
func ExtractUsers(page pagination.Page) ([]User, error) {
	casted := page.(UserPage).Body
	var response struct {
		Users []User `mapstructure:"users"`
	}

	err := mapstructure.Decode(casted, &response)
	return response.Users, err
}
