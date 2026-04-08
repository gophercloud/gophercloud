package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/metric/v1/metrics"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestQuery(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleQuerySuccessfully(t, fakeServer)

	opts := metrics.QueryOpts{
		Query: "up",
	}

	actual, err := metrics.Query(context.TODO(), serviceClient(fakeServer), opts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedQueryData, actual)
}

func TestLabels(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleLabelsSuccessfully(t, fakeServer)

	actual, err := metrics.Labels(context.TODO(), serviceClient(fakeServer), nil).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedLabels, actual)
}

func TestLabelValues(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleLabelValuesSuccessfully(t, fakeServer)

	actual, err := metrics.LabelValues(context.TODO(), serviceClient(fakeServer), "job", nil).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedLabelValues, actual)
}

func TestSeries(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleSeriesSuccessfully(t, fakeServer)

	opts := metrics.SeriesOpts{
		Match: []string{"up"},
	}

	actual, err := metrics.Series(context.TODO(), serviceClient(fakeServer), opts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedSeries, actual)
}

func TestTargets(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleTargetsSuccessfully(t, fakeServer)

	actual, err := metrics.Targets(context.TODO(), serviceClient(fakeServer), nil).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedTargetsData, actual)
}

func TestRuntimeInfo(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleRuntimeInfoSuccessfully(t, fakeServer)

	actual, err := metrics.RuntimeInfo(context.TODO(), serviceClient(fakeServer)).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedRuntimeInfoData, actual)
}

func TestCleanTombstones(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCleanTombstonesSuccessfully(t, fakeServer)

	err := metrics.CleanTombstones(context.TODO(), serviceClient(fakeServer)).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestDeleteSeries(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDeleteSeriesSuccessfully(t, fakeServer)

	opts := metrics.DeleteSeriesOpts{
		Match: []string{"up"},
	}

	err := metrics.DeleteSeries(context.TODO(), serviceClient(fakeServer), opts).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestSnapshot(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleSnapshotSuccessfully(t, fakeServer)

	actual, err := metrics.Snapshot(context.TODO(), serviceClient(fakeServer)).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedSnapshotData, actual)
}
