package testing

import (
	"net/url"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/usage"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestUsageTimezone(t *testing.T) {
	locations := []*time.Location{
		time.UTC,
		time.FixedZone("UTC+1", 3600),
		time.FixedZone("UTC+2", 2*3600),
		time.FixedZone("UTC-7", -7*3600),
		time.FixedZone("UTC+5:45", 5*3600+45*60),
		time.FixedZone("", 90*60),
		time.Local,
	}
	builders := map[string]func(*time.Time, *time.Time) (string, error){
		"single": func(start, end *time.Time) (string, error) {
			return (usage.SingleTenantOpts{Start: start, End: end}).ToUsageSingleTenantQuery()
		},
		"all": func(start, end *time.Time) (string, error) {
			return (usage.AllTenantsOpts{Start: start, End: end}).ToUsageAllTenantsQuery()
		},
	}
	for name, build := range builders {
		t.Run(name, func(t *testing.T) {
			for _, loc := range locations {
				t.Run(loc.String(), func(t *testing.T) {
					start := time.Date(2032, 12, 31, 23, 22, 22, 123456789, time.UTC).In(loc)
					end := start.Add(2 * time.Hour)
					originalStart, originalEnd := start, end
					raw, err := build(&start, &end)
					th.AssertNoErr(t, err)
					query, err := url.Parse(raw)
					th.AssertNoErr(t, err)
					th.CheckDeepEquals(t, url.Values{
						"start": {"2032-12-31T23:22:22.123456"},
						"end":   {"2033-01-01T01:22:22.123456"},
					}, query.Query())
					th.CheckDeepEquals(t, originalStart, start)
					th.CheckDeepEquals(t, originalEnd, end)
				})
			}
			t.Run("omitted", func(t *testing.T) {
				raw, err := build(nil, nil)
				th.AssertNoErr(t, err)
				th.CheckEquals(t, "", raw)
			})
		})
	}
}
