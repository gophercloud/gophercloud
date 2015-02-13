package users

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
	db "github.com/rackspace/gophercloud/openstack/db/v1/databases"
	"github.com/rackspace/gophercloud/pagination"
)

// User represents a database user
type User struct {
	// The user name
	Name string

	// The user password
	Password string

	Host string

	// The databases associated with this user
	Databases []db.Database
}

type UpdatePasswordsResult struct {
	gophercloud.ErrResult
}

type UpdateResult struct {
	gophercloud.ErrResult
}

type GetResult struct {
	gophercloud.Result
}

func (r GetResult) Extract() (*User, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var response struct {
		User User `mapstructure:"user"`
	}

	err := mapstructure.Decode(r.Body, &response)
	return &response.User, err
}

// AccessPage represents a single page of a paginated user collection.
type AccessPage struct {
	pagination.LinkedPageBase
}

// IsEmpty checks to see whether the collection is empty.
func (page AccessPage) IsEmpty() (bool, error) {
	users, err := ExtractDBs(page)
	if err != nil {
		return true, err
	}
	return len(users) == 0, nil
}

// NextPageURL will retrieve the next page URL.
func (page AccessPage) NextPageURL() (string, error) {
	type resp struct {
		Links []gophercloud.Link `mapstructure:"databases_links"`
	}

	var r resp
	err := mapstructure.Decode(page.Body, &r)
	if err != nil {
		return "", err
	}

	return gophercloud.ExtractNextURL(r.Links)
}

// ExtractDBs will convert a generic pagination struct into a more
// relevant slice of DB structs.
func ExtractDBs(page pagination.Page) ([]db.Database, error) {
	casted := page.(AccessPage).Body

	var response struct {
		DBs []db.Database `mapstructure:"databases"`
	}

	err := mapstructure.Decode(casted, &response)
	return response.DBs, err
}

type GrantAccessResult struct {
	gophercloud.ErrResult
}

type RevokeAccessResult struct {
	gophercloud.ErrResult
}
