package testing

import (
	"net/url"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/workflow/v2/crontriggers"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestCreateTimezone(t *testing.T) {
	build := func(t *testing.T, timestamp *time.Time) map[string]any {
		t.Helper()
		b, err := (crontriggers.CreateOpts{WorkflowID: "workflow", Name: "trigger", FirstExecutionTime: timestamp}).ToCronTriggerCreateMap()
		th.AssertNoErr(t, err)
		return b
	}
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
			timestamp := time.Date(2032, 12, 31, 23, 22, 22, 123456789, time.UTC).In(loc)
			original := timestamp
			b := build(t, &timestamp)
			th.CheckEquals(t, "2032-12-31 23:22", b["first_execution_time"])
			th.CheckDeepEquals(t, original, timestamp)
		})
	}
	t.Run("omitted", func(t *testing.T) {
		_, present := build(t, nil)["first_execution_time"]
		th.CheckEquals(t, false, present)
	})
}

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
			for _, operator := range []crontriggers.FilterType{"", crontriggers.FilterGTE} {
				t.Run(string(operator), func(t *testing.T) {
					timestamp := time.Date(2032, 12, 31, 23, 22, 22, 123456789, time.UTC).In(loc)
					filter := crontriggers.ListDateFilter{Filter: operator, Value: timestamp}
					original := filter
					expected := "2032-12-31 23:22:22"
					if operator != "" {
						expected = string(operator) + ":" + expected
					}
					th.CheckEquals(t, expected, filter.String())
					opts := crontriggers.ListOpts{
						FirstExecutionTime: &filter,
						NextExecutionTime:  &filter,
						CreatedAt:          &filter,
						UpdatedAt:          &filter,
					}
					raw, err := opts.ToCronTriggerListQuery()
					th.AssertNoErr(t, err)
					query, err := url.Parse(raw)
					th.AssertNoErr(t, err)
					th.CheckDeepEquals(t, url.Values{
						"first_execution_time": {expected},
						"next_execution_time":  {expected},
						"created_at":           {expected},
						"updated_at":           {expected},
					}, query.Query())
					th.CheckDeepEquals(t, original, filter)
				})
			}
		})
	}
	t.Run("omitted", func(t *testing.T) {
		raw, err := (crontriggers.ListOpts{}).ToCronTriggerListQuery()
		th.AssertNoErr(t, err)
		th.CheckEquals(t, "", raw)
	})
}
