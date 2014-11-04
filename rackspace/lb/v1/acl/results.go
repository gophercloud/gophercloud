package acl

import (
	"github.com/mitchellh/mapstructure"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

type AccessList []NetworkItem

type NetworkItem struct {
	Address string
	ID      int
	Type    Type
}

type Type string

const (
	ALLOW Type = "ALLOW"
	DENY  Type = "DENY"
)

// AccessListPage is the page returned by a pager when traversing over a collection of
// network items in an access list.
type AccessListPage struct {
	pagination.SinglePageBase
}

// IsEmpty checks whether an AccessListPage struct is empty.
func (p AccessListPage) IsEmpty() (bool, error) {
	is, err := ExtractAccessList(p)
	if err != nil {
		return true, nil
	}
	return len(is) == 0, nil
}

// ExtractAccessList accepts a Page struct, specifically an AccessListPage
// struct, and extracts the elements into a slice of NetworkItem structs. In
// other words, a generic collection is mapped into a relevant slice.
func ExtractAccessList(page pagination.Page) (AccessList, error) {
	var resp struct {
		List AccessList `mapstructure:"accessList" json:"accessList"`
	}

	err := mapstructure.Decode(page.(AccessListPage).Body, &resp)

	return resp.List, err
}

type CreateResult struct {
	gophercloud.ErrResult
}

type DeleteResult struct {
	gophercloud.ErrResult
}
