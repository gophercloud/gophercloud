package testing

import (
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/trusts"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestCreateTimezone(t *testing.T) {
	build := func(t *testing.T, timestamp *time.Time) map[string]any {
		t.Helper()
		b, err := (trusts.CreateOpts{TrustorUserID: "trustor", TrusteeUserID: "trustee", ExpiresAt: timestamp}).ToTrustCreateMap()
		th.AssertNoErr(t, err)
		return b["trust"].(map[string]any)
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
			th.CheckEquals(t, "2032-12-31T23:22:22.123456Z", b["expires_at"])
			th.CheckDeepEquals(t, original, timestamp)
		})
	}
	t.Run("omitted", func(t *testing.T) {
		_, present := build(t, nil)["expires_at"]
		th.CheckEquals(t, false, present)
	})
}
