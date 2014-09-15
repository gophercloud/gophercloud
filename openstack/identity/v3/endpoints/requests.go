package endpoints

import (
	"strconv"

	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/utils"
)

// maybeString returns nil for empty strings and nil for empty.
func maybeString(original string) *string {
	if original != "" {
		return &original
	}
	return nil
}

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
func Create(client *gophercloud.ServiceClient, opts EndpointOpts) (*Endpoint, error) {
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

	type response struct {
		Endpoint Endpoint `json:"endpoint"`
	}

	// Ensure that EndpointOpts is fully populated.
	if opts.Availability == "" {
		return nil, ErrAvailabilityRequired
	}
	if opts.Name == "" {
		return nil, ErrNameRequired
	}
	if opts.URL == "" {
		return nil, ErrURLRequired
	}
	if opts.ServiceID == "" {
		return nil, ErrServiceIDRequired
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
	reqBody.Endpoint.Region = maybeString(opts.Region)

	var respBody response
	_, err := perigee.Request("POST", getListURL(client), perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &respBody,
		OkCodes:     []int{201},
	})
	if err != nil {
		return nil, err
	}

	return &respBody.Endpoint, nil
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
func List(client *gophercloud.ServiceClient, opts ListOpts) gophercloud.Pager {
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

	u := getListURL(client) + utils.BuildQuery(q)
	return gophercloud.NewLinkedPager(client, u)
}

// Update changes an existing endpoint with new data.
// All fields are optional in the provided EndpointOpts.
func Update(client *gophercloud.ServiceClient, endpointID string, opts EndpointOpts) (*Endpoint, error) {
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

	type response struct {
		Endpoint Endpoint `json:"endpoint"`
	}

	reqBody := request{Endpoint: endpoint{}}
	reqBody.Endpoint.Interface = maybeString(string(opts.Availability))
	reqBody.Endpoint.Name = maybeString(opts.Name)
	reqBody.Endpoint.Region = maybeString(opts.Region)
	reqBody.Endpoint.URL = maybeString(opts.URL)
	reqBody.Endpoint.ServiceID = maybeString(opts.ServiceID)

	var respBody response
	_, err := perigee.Request("PATCH", getEndpointURL(client, endpointID), perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &respBody,
		OkCodes:     []int{200},
	})
	if err != nil {
		return nil, err
	}

	return &respBody.Endpoint, nil
}

// Delete removes an endpoint from the service catalog.
func Delete(client *gophercloud.ServiceClient, endpointID string) error {
	_, err := perigee.Request("DELETE", getEndpointURL(client, endpointID), perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{204},
	})
	return err
}
