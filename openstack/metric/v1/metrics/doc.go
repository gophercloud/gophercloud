/*
Package metrics provides interaction with the Prometheus HTTP API through the
OpenStack Aetos service (metric-storage).

Aetos acts as a Keystone-authenticated reverse proxy for Prometheus, enforcing
multi-tenant RBAC on all Prometheus queries. It exposes the Prometheus HTTP API
under /api/v1/ with OpenStack authentication.

Example to Create a Metric Service Client

	metricClient, err := openstack.NewMetricV1(providerClient, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	if err != nil {
		panic(err)
	}

Example to Query Metrics

	opts := metrics.QueryOpts{
		Query: "up",
	}

	result, err := metrics.Query(ctx, metricClient, opts).Extract()
	if err != nil {
		panic(err)
	}

	for _, v := range result.Result {
		fmt.Printf("%s: %v\n", v.Metric["__name__"], v.Value)
	}

Example to List Label Names

	labels, err := metrics.Labels(ctx, metricClient, nil).Extract()
	if err != nil {
		panic(err)
	}

	for _, label := range labels {
		fmt.Println(label)
	}

Example to Find Series by Label Matchers

	opts := metrics.SeriesOpts{
		Match: []string{"up", `process_start_time_seconds{job="prometheus"}`},
	}

	series, err := metrics.Series(ctx, metricClient, opts).Extract()
	if err != nil {
		panic(err)
	}

	for _, s := range series {
		fmt.Printf("%v\n", s)
	}
*/
package metrics
