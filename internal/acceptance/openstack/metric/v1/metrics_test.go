//go:build acceptance || metric || metrics

package v1

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/metric/v1/metrics"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestMetricQuery(t *testing.T) {
	client, err := clients.NewMetricV1Client()
	th.AssertNoErr(t, err)

	result, err := metrics.Query(context.TODO(), client, metrics.QueryOpts{
		Query: "up",
	}).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, result)
	th.AssertEquals(t, result.ResultType, "vector")
}

func TestMetricLabels(t *testing.T) {
	client, err := clients.NewMetricV1Client()
	th.AssertNoErr(t, err)

	result, err := metrics.Labels(context.TODO(), client, nil).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, result)
	if len(result) == 0 {
		t.Fatal("expected at least one label")
	}
}

func TestMetricLabelValues(t *testing.T) {
	client, err := clients.NewMetricV1Client()
	th.AssertNoErr(t, err)

	result, err := metrics.LabelValues(context.TODO(), client, "__name__", nil).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, result)
	if len(result) == 0 {
		t.Fatal("expected at least one label value")
	}
}

func TestMetricSeries(t *testing.T) {
	client, err := clients.NewMetricV1Client()
	th.AssertNoErr(t, err)

	result, err := metrics.Series(context.TODO(), client, metrics.SeriesOpts{
		Match: []string{`{__name__="up"}`},
	}).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, result)
	if len(result) == 0 {
		t.Fatal("expected at least one series")
	}
}

func TestMetricTargets(t *testing.T) {
	client, err := clients.NewMetricV1Client()
	th.AssertNoErr(t, err)

	result, err := metrics.Targets(context.TODO(), client, nil).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, result)
	if len(result.ActiveTargets) == 0 {
		t.Fatal("expected at least one active target")
	}
}

func TestMetricRuntimeInfo(t *testing.T) {
	client, err := clients.NewMetricV1Client()
	th.AssertNoErr(t, err)

	result, err := metrics.RuntimeInfo(context.TODO(), client).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, result)
	if result.StorageRetention == "" {
		t.Fatal("expected non-empty StorageRetention")
	}
}
