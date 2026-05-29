package metrics

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
)

// QueryOptsBuilder allows extensions to add parameters to the Query request.
type QueryOptsBuilder interface {
	ToMetricQueryQuery() (string, error)
}

// QueryOpts contains the options for a Prometheus instant query.
type QueryOpts struct {
	// Query is the PromQL query string.
	Query string `q:"query" required:"true"`

	// Time is the evaluation timestamp (RFC3339 or Unix timestamp).
	Time string `q:"time"`

	// Timeout is the evaluation timeout.
	Timeout string `q:"timeout"`
}

// ToMetricQueryQuery formats QueryOpts into a query string.
func (opts QueryOpts) ToMetricQueryQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// Query performs a Prometheus instant query.
func Query(ctx context.Context, client *gophercloud.ServiceClient, opts QueryOptsBuilder) (r QueryResult) {
	url := queryURL(client)
	if opts != nil {
		query, err := opts.ToMetricQueryQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	resp, err := client.Get(ctx, url, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// LabelsOptsBuilder allows extensions to add parameters to the Labels request.
type LabelsOptsBuilder interface {
	ToMetricLabelsQuery() (string, error)
}

// LabelsOpts contains the options for listing label names.
type LabelsOpts struct {
	// Match is a list of series selectors to filter labels by.
	Match []string `q:"match[]"`

	// Start is the start timestamp (RFC3339 or Unix timestamp).
	Start string `q:"start"`

	// End is the end timestamp (RFC3339 or Unix timestamp).
	End string `q:"end"`
}

// ToMetricLabelsQuery formats LabelsOpts into a query string.
func (opts LabelsOpts) ToMetricLabelsQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// Labels returns a list of label names.
func Labels(ctx context.Context, client *gophercloud.ServiceClient, opts LabelsOptsBuilder) (r LabelsResult) {
	url := labelsURL(client)
	if opts != nil {
		query, err := opts.ToMetricLabelsQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	resp, err := client.Get(ctx, url, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// LabelValuesOptsBuilder allows extensions to add parameters to the LabelValues request.
type LabelValuesOptsBuilder interface {
	ToMetricLabelValuesQuery() (string, error)
}

// LabelValuesOpts contains the options for listing label values.
type LabelValuesOpts struct {
	// Match is a list of series selectors to filter label values by.
	Match []string `q:"match[]"`

	// Start is the start timestamp (RFC3339 or Unix timestamp).
	Start string `q:"start"`

	// End is the end timestamp (RFC3339 or Unix timestamp).
	End string `q:"end"`
}

// ToMetricLabelValuesQuery formats LabelValuesOpts into a query string.
func (opts LabelValuesOpts) ToMetricLabelValuesQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// LabelValues returns a list of label values for a given label name.
func LabelValues(ctx context.Context, client *gophercloud.ServiceClient, name string, opts LabelValuesOptsBuilder) (r LabelValuesResult) {
	url := labelValuesURL(client, name)
	if opts != nil {
		query, err := opts.ToMetricLabelValuesQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	resp, err := client.Get(ctx, url, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// SeriesOptsBuilder allows extensions to add parameters to the Series request.
type SeriesOptsBuilder interface {
	ToMetricSeriesQuery() (string, error)
}

// SeriesOpts contains the options for finding series by label matchers.
type SeriesOpts struct {
	// Match is a list of series selectors. At least one must be provided.
	Match []string `q:"match[]" required:"true"`

	// Start is the start timestamp (RFC3339 or Unix timestamp).
	Start string `q:"start"`

	// End is the end timestamp (RFC3339 or Unix timestamp).
	End string `q:"end"`
}

// ToMetricSeriesQuery formats SeriesOpts into a query string.
func (opts SeriesOpts) ToMetricSeriesQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// Series returns the list of time series that match certain label sets.
func Series(ctx context.Context, client *gophercloud.ServiceClient, opts SeriesOptsBuilder) (r SeriesResult) {
	url := seriesURL(client)
	if opts != nil {
		query, err := opts.ToMetricSeriesQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	resp, err := client.Get(ctx, url, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// TargetsOptsBuilder allows extensions to add parameters to the Targets request.
type TargetsOptsBuilder interface {
	ToMetricTargetsQuery() (string, error)
}

// TargetsOpts contains the options for listing targets.
type TargetsOpts struct {
	// State filters targets by state ("active", "dropped", "any").
	State string `q:"state"`
}

// ToMetricTargetsQuery formats TargetsOpts into a query string.
func (opts TargetsOpts) ToMetricTargetsQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// Targets returns an overview of the current state of Prometheus target discovery.
func Targets(ctx context.Context, client *gophercloud.ServiceClient, opts TargetsOptsBuilder) (r TargetsResult) {
	url := targetsURL(client)
	if opts != nil {
		query, err := opts.ToMetricTargetsQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	resp, err := client.Get(ctx, url, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// RuntimeInfo returns runtime information about the Prometheus server.
func RuntimeInfo(ctx context.Context, client *gophercloud.ServiceClient) (r RuntimeInfoResult) {
	resp, err := client.Get(ctx, runtimeInfoURL(client), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CleanTombstones removes deleted data from disk and cleans up the existing tombstones.
func CleanTombstones(ctx context.Context, client *gophercloud.ServiceClient) (r CleanTombstonesResult) {
	resp, err := client.Post(ctx, cleanTombstonesURL(client), nil, nil, &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteSeriesOptsBuilder allows extensions to add parameters to the DeleteSeries request.
type DeleteSeriesOptsBuilder interface {
	ToMetricDeleteSeriesQuery() (string, error)
}

// DeleteSeriesOpts contains the options for deleting time series.
type DeleteSeriesOpts struct {
	// Match is a list of series selectors. At least one must be provided.
	Match []string `q:"match[]" required:"true"`

	// Start is the start timestamp (RFC3339 or Unix timestamp).
	Start string `q:"start"`

	// End is the end timestamp (RFC3339 or Unix timestamp).
	End string `q:"end"`
}

// ToMetricDeleteSeriesQuery formats DeleteSeriesOpts into a query string.
func (opts DeleteSeriesOpts) ToMetricDeleteSeriesQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// DeleteSeries deletes data for a selection of series in a time range.
func DeleteSeries(ctx context.Context, client *gophercloud.ServiceClient, opts DeleteSeriesOptsBuilder) (r DeleteSeriesResult) {
	url := deleteSeriesURL(client)
	if opts != nil {
		query, err := opts.ToMetricDeleteSeriesQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	resp, err := client.Post(ctx, url, nil, nil, &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Snapshot creates a snapshot of all current data into snapshots/<datetime>-<rand>
// under the TSDB's data directory and returns the directory as response.
func Snapshot(ctx context.Context, client *gophercloud.ServiceClient) (r SnapshotResult) {
	resp, err := client.Post(ctx, snapshotURL(client), nil, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
