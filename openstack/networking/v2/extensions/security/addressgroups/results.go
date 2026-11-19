package addressgroups

import (
	"encoding/json"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// AddressGroup represents a container for address groups.
type AddressGroup struct {
	// Unique identifier for the address_group object.
	ID string `json:"id"`

	// Human readable name for the address group (255 characters limit). Does not have to be unique.
	Name string `json:"name"`

	// Human readable description for the address group (255 characters limit).
	Description string `json:"description"`

	// ProjectID is the project owner of this address group.
	ProjectID string `json:"project_id"`

	// Array of address. It supports both CIDR and IP range objects.
	// An example of addresses: [“132.168.4.12/24”, “132.168.5.12-132.168.5.24”, “2001::db8::f00/64”]
	Addresses []string `json:"addresses"`

	// Timestamp when the address group was created
	CreatedAt time.Time `json:"-"`

	// Timestamp when the address group was last updated
	UpdatedAt time.Time `json:"-"`
}

func (r *AddressGroup) UnmarshalJSON(b []byte) error {
	type tmp AddressGroup

	// Support for older neutron time format
	var s1 struct {
		tmp
		CreatedAt gophercloud.JSONRFC3339NoZ `json:"created_at"`
		UpdatedAt gophercloud.JSONRFC3339NoZ `json:"updated_at"`
	}

	err := json.Unmarshal(b, &s1)
	if err == nil {
		*r = AddressGroup(s1.tmp)
		r.CreatedAt = time.Time(s1.CreatedAt)
		r.UpdatedAt = time.Time(s1.UpdatedAt)

		return nil
	}

	// Support for newer neutron time format
	var s2 struct {
		tmp
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	err = json.Unmarshal(b, &s2)
	if err != nil {
		return err
	}

	*r = AddressGroup(s2.tmp)
	r.CreatedAt = time.Time(s2.CreatedAt)
	r.UpdatedAt = time.Time(s2.UpdatedAt)

	return nil
}

// AddressGroupPage is the page returned by a pager when traversing over a
// collection of address group addresses.
type AddressGroupPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of address groups has
// reached the end of a page and the pager seeks to traverse over a new one. In
// order to do this, it needs to construct the next page's URL.
func (r AddressGroupPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"address_groups_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a AddressGroupPage struct is empty.
func (r AddressGroupPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	is, err := ExtractGroups(r)
	return len(is) == 0, err
}

// ExtractGroups accepts a Page struct, specifically a AddressGroupPage struct,
// and extracts the elements into a slice of AddressGroup structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractGroups(r pagination.Page) ([]AddressGroup, error) {
	var s struct {
		AddressGroups []AddressGroup `json:"address_groups"`
	}
	err := (r.(AddressGroupPage)).ExtractInto(&s)
	return s.AddressGroups, err
}

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a address group.
func (r commonResult) Extract() (*AddressGroup, error) {
	var s struct {
		AddressGroup *AddressGroup `json:"address_group"`
	}
	err := r.ExtractInto(&s)
	return s.AddressGroup, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a AddressGroup.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a AddressGroup.
type GetResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its
// ExtractErr method to determine if the request succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

// UpdateResult represents the result of an update address group operation. Call its Extract
// method to interpret it as a AddressGroup.
type UpdateResult struct {
	commonResult
}

// AddAddressesResult represents the result of an add addresses operation. Call its Extract
// method to interpret it as a AddressGroup.
type AddAddressesResult struct {
	commonResult
}

// RemoveAddressesResult represents the result of a remove addresses operation. Call its Extract
// method to interpret it as a AddressGroup.
type RemoveAddressesResult struct {
	commonResult
}
