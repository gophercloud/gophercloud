package drivers

import (
	"context"
	"encoding/json"
	"io"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListDriversOptsBuilder allows extensions to add additional parameters to the
// ListDrivers request.
type ListDriversOptsBuilder interface {
	ToListDriversOptsQuery() (string, error)
}

// ListDriversOpts defines query options that can be passed to ListDrivers
type ListDriversOpts struct {
	// Provide detailed information about the drivers
	Detail bool `q:"detail"`

	// Filter the list by the type of the driver
	Type string `q:"type"`
}

// ToListDriversOptsQuery formats a ListOpts into a query string
func (opts ListDriversOpts) ToListDriversOptsQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// ListDrivers makes a request against the API to list all drivers
func ListDrivers(client *gophercloud.ServiceClient, opts ListDriversOptsBuilder) pagination.Pager {
	url := driversURL(client)
	if opts != nil {
		query, err := opts.ToListDriversOptsQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return DriverPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// GetDriverDetails Shows details for a driver
func GetDriverDetails(ctx context.Context, client *gophercloud.ServiceClient, driverName string) (r GetDriverResult) {
	resp, err := client.Get(ctx, driverDetailsURL(client, driverName), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetDriverProperties Shows the required and optional parameters that
// driverName expects to be supplied in the driver_info field for every
// Node it manages
func GetDriverProperties(ctx context.Context, client *gophercloud.ServiceClient, driverName string) (r GetPropertiesResult) {
	resp, err := client.Get(ctx, driverPropertiesURL(client, driverName), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetDriverDiskProperties Show the required and optional parameters that
// driverName expects to be supplied in the node’s raid_config field, if a
// RAID configuration change is requested.
func GetDriverDiskProperties(ctx context.Context, client *gophercloud.ServiceClient, driverName string) (r GetDiskPropertiesResult) {
	resp, err := client.Get(ctx, driverDiskPropertiesURL(client, driverName), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListVendorPassthruMethods retrieves all vendor passthru methods available for the given Driver.
func ListVendorPassthruMethods(ctx context.Context, client *gophercloud.ServiceClient, driverName string) (r ListVendorPassthruMethodsResult) {
	resp, err := client.Get(ctx, vendorPassthruMethodsURL(client, driverName), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// VendorPassthruCallOpts defines options for invoking a vendor passthru method.
type VendorPassthruCallOpts struct {
	// Method is the driver-specific passthru method name.
	Method string

	// Body is the JSON body forwarded to the driver.
	Body any
}

// CallVendorPassthru invokes a vendor passthru method for the given Driver.
func CallVendorPassthru(ctx context.Context, client *gophercloud.ServiceClient, driverName string, httpMethod string, opts VendorPassthruCallOpts) (r CallVendorPassthruResult) {
	query, err := vendorPassthruCallQuery(opts.Method)
	if err != nil {
		r.Err = err
		return
	}

	url := vendorPassthruCallURL(client, driverName) + query
	resp, err := client.Request(ctx, httpMethod, url, &gophercloud.RequestOpts{
		JSONBody:         opts.Body,
		OkCodes:          []int{200, 202, 204},
		KeepResponseBody: true,
	})
	if resp != nil {
		r.Header = resp.Header
		r.StatusCode = resp.StatusCode
	}
	if err != nil {
		r.Err = err
		return
	}
	if resp == nil || resp.Body == nil {
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		r.Err = err
		return
	}
	if len(responseBody) == 0 {
		return
	}
	r.Err = json.Unmarshal(responseBody, &r.Body)
	return
}

func vendorPassthruCallQuery(method string) (string, error) {
	q, err := gophercloud.BuildQueryString(struct {
		Method string `q:"method" required:"true"`
	}{Method: method})
	return q.String(), err
}
