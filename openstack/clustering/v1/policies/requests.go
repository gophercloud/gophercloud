package policies

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOptsBuilder Builder.
type ListOptsBuilder interface {
	ToPolicyListQuery() (string, error)
}

// ListOpts params
type ListOpts struct {
	// Limit limits the number of Policies to return.
	Limit int `q:"limit"`

	// Marker and Limit control paging. Marker instructs List where to start listing from.
	Marker string `q:"marker"`

	// Sorts the response by one or more attribute and optional sort direction combinations.
	Sort string `q:"sort"`

	// GlobalProject indicates whether to include resources for all projects or resources for the current project
	GlobalProject *bool `q:"global_project"`

	// Name to filter the response by the specified name property of the object
	Name string `q:"name"`

	// Filter the response by the specified type property of the object
	Type string `q:"type"`
}

// ToPolicyListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToPolicyListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List instructs OpenStack to retrieve a list of policies.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := policyListURL(client)
	if opts != nil {
		query, err := opts.ToPolicyListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		p := PolicyPage{pagination.MarkerPageBase{PageResult: r}}
		p.MarkerPageBase.Owner = p
		return p
	})
}

// CreateOpts params
type CreateOpts struct {
	Name string `json:"name"`
	Spec Spec   `json:"spec"`
}

type Spec struct {
	Description string                 `json:"description"`
	Properties  map[string]interface{} `json:"properties"`
	Type        string                 `json:"type"`
	Version     Version                `json:"version"`
}

// ToPolicyCreateMap formats a CreateOpts into a body map.
func (opts CreateOpts) ToPolicyCreateMap() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"policy": b}, nil
}

type Version string

// custom unmarshal function to handle Version returned as either float64 or string
func (v *Version) UnmarshalJSON(b []byte) error {
	if b[0] == '"' {
		return json.Unmarshal(b, (*string)(v))
	}

	var f float64
	if err := json.Unmarshal(b, &f); err != nil {
		return err
	}

	if f == 1 {
		blah := fmt.Sprintf("%.1f", f)
		*v = Version(blah)
		return nil
	} else {
		blah := strconv.FormatFloat(f, 'f', -1, 64)
		*v = Version(blah)
		return nil
	}
}

// Create makes a request against the API to create a policy
func Create(client *gophercloud.ServiceClient, opts CreateOpts) (r CreateResult) {
	b, err := opts.ToPolicyCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(policyCreateURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	return
}
