package testing

import (
	"net/url"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/workflow/v2/workflows"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestDateFilterTimezone(t *testing.T) {
	locations := []*time.Location{
		time.UTC,
		time.FixedZone("UTC+1", 3600),
		time.FixedZone("UTC+2", 2*3600),
		time.FixedZone("UTC-7", -7*3600),
		time.FixedZone("UTC+5:45", 5*3600+45*60),
		time.FixedZone("", 90*60),
		time.Local,
	}
	for _, loc := range locations {
		t.Run(loc.String(), func(t *testing.T) {
			for _, operator := range []workflows.FilterType{"", workflows.FilterGTE} {
				t.Run(string(operator), func(t *testing.T) {
					timestamp := time.Date(2032, 12, 31, 23, 22, 22, 123456789, time.UTC).In(loc)
					filter := workflows.ListDateFilter{Filter: operator, Value: timestamp}
					original := filter
					expected := "2032-12-31 23:22:22"
					if operator != "" {
						expected = string(operator) + ":" + expected
					}
					th.CheckEquals(t, expected, filter.String())
					opts := workflows.ListOpts{
						CreatedAt: &filter,
						UpdatedAt: &filter,
					}
					raw, err := opts.ToWorkflowListQuery()
					th.AssertNoErr(t, err)
					query, err := url.Parse(raw)
					th.AssertNoErr(t, err)
					th.CheckDeepEquals(t, url.Values{
						"created_at": {expected},
						"updated_at": {expected},
					}, query.Query())
					th.CheckDeepEquals(t, original, filter)
				})
			}
		})
	}
	t.Run("omitted", func(t *testing.T) {
		raw, err := (workflows.ListOpts{}).ToWorkflowListQuery()
		th.AssertNoErr(t, err)
		th.CheckEquals(t, "", raw)
	})
}
