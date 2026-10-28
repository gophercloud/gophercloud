package testing

import (
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/objectstorage/v1/objects"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestDownloadTimezone(t *testing.T) {
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
			opts := objects.DownloadOpts{IfModifiedSince: timestamp, IfUnmodifiedSince: timestamp}
			original := opts
			headers, query, err := opts.ToObjectDownloadParams()
			th.AssertNoErr(t, err)
			th.CheckEquals(t, "", query)
			for _, name := range []string{"If-Modified-Since", "If-Unmodified-Since"} {
				th.CheckEquals(t, "Fri, 31 Dec 2032 23:22:22 GMT", headers[name])
				parsed, err := http.ParseTime(headers[name])
				th.AssertNoErr(t, err)
				th.CheckEquals(t, true, parsed.Equal(timestamp.Truncate(time.Second)))
			}
			th.CheckDeepEquals(t, original, opts)
		})
	}
	t.Run("omitted", func(t *testing.T) {
		headers, _, err := (objects.DownloadOpts{}).ToObjectDownloadParams()
		th.AssertNoErr(t, err)
		for _, name := range []string{"If-Modified-Since", "If-Unmodified-Since"} {
			_, present := headers[name]
			th.CheckEquals(t, false, present)
		}
	})
}
