package meters

import (
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

func TestListMeters(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleMeterListSuccessfully(t)

	res, err := List(client.ServiceClient(), ListOpts{})
	th.AssertNoErr(t, err)
	list, err := res.Extract()
	th.AssertNoErr(t, err)

	if len(*list) != 2 {
		t.Fatalf("Expected 2 meters, got %d", len(*list))
	}
	th.CheckDeepEquals(t, MeterHerp, (*list)[0])
	th.CheckDeepEquals(t, MeterDerp, (*list)[1])
}

func TestMeterStatistics(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleMeterStatisticsSuccessfully(t)

	res, err := MeterStatistics(client.ServiceClient(), "memory", MeterStatisticsOpts{})
	th.AssertNoErr(t, err)
	list, err := res.Extract()
	th.AssertNoErr(t, err)

	if len(*list) != 2 {
		t.Fatalf("Expected 2 statistics, got %d", len(*list))
	}
	th.CheckDeepEquals(t, StatisticsHerp, (*list)[0])
	th.CheckDeepEquals(t, StatisticsDerp, (*list)[1])
}
