package endpoints

import (
	"strconv"

	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/utils"
	"github.com/rackspace/gophercloud/pagination"
)

// EndpointOpts contains the subset of Endpoint attributes that should be used to create or update an Endpoint.
type EndpointOpts struct {
	Availability gophercloud.Availability
	Name         string
	Region       string
	URL          string
	ServiceID    string
}

// Create inserts a new Endpoint into the service catalog.
// Within EndpointOpts, Region may be omitted by being left as "", but all other fields are required.
func Create(client *gophercloud.ServiceClient, opts EndpointOpts) CreateResult {
	// Redefined so that Region can be re-typed as a *string, which can be omitted from the JSON output.
	type endpoint struct {
		Interface string  `json:"interface"`
		Name      string  `json:"name"`
		Region    *string `json:"region,omitempty"`
		URL       string  `json:"url"`
		ServiceID string  `json:"service_id"`
	}

	type request struct {
		Endpoint endpoint `json:"endpoint"`
	}

	// Ensure that EndpointOpts is fully populated.
	if opts.Availability == "" {
		return createErr(ErrAvailabilityRequired)
	}
	if opts.Name == "" {
		return createErr(ErrNameRequired)
	}
	if opts.URL == "" {
		return createErr(ErrURLRequired)
	}
	if opts.ServiceID == "" {
		return createErr(ErrServiceIDRequired)
	}

	// Populate the request body.
	reqBody := request{
		Endpoint: endpoint{
			Interface: string(opts.Availability),
			Name:      opts.Name,
			URL:       opts.URL,
			ServiceID: opts.ServiceID,
		},
	}
	reqBody.Endpoint.Region = gophercloud.MaybeString(opts.Region)

	var result CreateResult
	_, result.Err = perigee.Request("POST", listURL(client), perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &result.Resp,
		OkCodes:     []int{201},
	})
	return result
}

// ListOpts allows finer control over the the endpoints returned by a List call.
// All fields are optional.
type ListOpts struct {
	Availability gophercloud.Availability
	ServiceID    string
	Page         int
	PerPage      int
}

// List enumerates endpoints in a paginated collection, optionally filtered by ListOpts criteria.
func List(client *gophercloud.ServiceClient, opts ListOpts) pagination.Pager {
	q := make(map[string]string)
	if opts.Availability != "" {
		q["interface"] = string(opts.Availability)
	}
	if opts.ServiceID != "" {
		q["service_id"] = opts.ServiceID
	}
	if opts.Page != 0 {
		q["page"] = strconv.Itoa(opts.Page)
	}
	if opts.PerPage != 0 {
		q["per_page"] = strconv.Itoa(opts.Page)
	}

	createPage := func(r pagination.LastHTTPResponse) pagination.Page {
		return EndpointPage{pagination.LinkedPageBase{LastHTTPResponse: r}}
	}

	u := listURL(client) + utils.BuildQuery(q)
	return pagination.NewPager(client, u, createPage)
}

// Update changes an existing endpoint with new data.
// All fields are optional in the provided EndpointOpts.
func Update(client *gophercloud.ServiceClient, endpointID string, opts EndpointOpts) UpdateResult {
	type endpoint struct {
		Interface *string `json:"interface,omitempty"`
		Name      *string `json:"name,omitempty"`
		Region    *string `json:"region,omitempty"`
		URL       *string `json:"url,omitempty"`
		ServiceID *string `json:"service_id,omitempty"`
	}

	type request struct {
		Endpoint endpoint `json:"endpoint"`
	}

	reqBody := request{Endpoint: endpoint{}}
	reqBody.Endpoint.Interface = gophercloud.MaybeString(string(opts.Availability))
	reqBody.Endpoint.Name = gophercloud.MaybeString(opts.Name)
	reqBody.Endpoint.Region = gophercloud.MaybeString(opts.Region)
	reqBody.Endpoint.URL = gophercloud.MaybeString(opts.URL)
	reqBody.Endpoint.ServiceID = gophercloud.MaybeString(opts.ServiceID)

	var result UpdateResult
	_, result.Err = perigee.Request("PATCH", endpointURL(client, endpointID), perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &result.Resp,
		OkCodes:     []int{200},
	})
	return result
}

// Delete removes an endpoint from the service catalog.
func Delete(client *gophercloud.ServiceClient, endpointID string) error {
	_, err := perigee.Request("DELETE", endpointURL(client, endpointID), perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{204},
	})
	return err
}
