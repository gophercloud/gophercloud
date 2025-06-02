package leases

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

type CreateOpts struct {
	Name          string        `json:"name"`
	StartDate     string        `json:"start_date"`
	EndDate       string        `json:"end_date"`
	BeforeEndDate string        `json:"before_end_date"`
	Reservations  []Reservation `json:"reservations"`
	ResourceType  string        `json:"resource_type"`
	Events        []Event       `json:"events"`
}

type CreateOptsBuilder interface {
	ToLeaseCreateMap() (map[string]any, error)
}

// ToLeaseCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToLeaseCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "lease")
}

func List(client *gophercloud.ServiceClient) pagination.Pager {
	url := listURL(client)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return LeasePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

func Get(ctx context.Context, client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	body, err := opts.ToLeaseCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(ctx, createURL(client), body, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
