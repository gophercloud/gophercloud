package meters

import (
	"fmt"
	"net/http"

	"github.com/rackspace/gophercloud"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToMeterListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the server attributes you want to see returned.
type ListOpts struct {
}

// ToServerListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToMeterListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

// List makes a request against the API to list meters accessible to you.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) listResult {
	var res listResult
	url := listURL(client)

	if opts != nil {
		query, err := opts.ToMeterListQuery()
		if err != nil {
			res.Err = err
			return res
		}
		url += query
	}

	_, res.Err = client.Get(url, &res.Body, &gophercloud.RequestOpts{})
	return res
}

// StatisticsOptsBuilder allows extensions to add additional parameters to the
// List request.
type MeterStatisticsOptsBuilder interface {
	ToMeterStatisticsQuery() (string, error)
}

// StatisticsOpts allows the filtering and sorting of collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the server attributes you want to see returned.
type MeterStatisticsOpts struct {
	QueryField string `q:"q.field"`
	QueryOp    string `q:"q.op"`
	QueryValue string `q:"q.value"`

	// Optional group by
	GroupBy string `q:"groupby"`

	// Optional number of seconds in a period
	Period int `q:"period"`
}

// ToStatisticsQuery formats a StatisticsOpts into a query string.
func (opts MeterStatisticsOpts) ToMeterStatisticsQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

// List makes a request against the API to list meters accessible to you.
func MeterStatistics(client *gophercloud.ServiceClient, n string, opts MeterStatisticsOptsBuilder) statisticsResult {
	var res statisticsResult
	url := statisticsURL(client, n)

	if opts != nil {
		query, err := opts.ToMeterStatisticsQuery()
		if err != nil {
			res.Err = err
			return res
		}
		url += query
	}

	var b *http.Response
	b, res.Err = client.Get(url, &res.Body, &gophercloud.RequestOpts{})
	fmt.Printf("%+v\n%+v\n", res, b)
	return res
}
