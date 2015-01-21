package base

import (
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

func TestGetHomeDocument(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetSuccessfully(t)

	actual, err := Get(fake.ServiceClient()).Extract()
	th.CheckNoErr(t, err)

	expected := HomeDocument{
		"rel/cdn": `{
        "href-template": "services{?marker,limit}",
        "href-vars": {
            "marker": "param/marker",
            "limit": "param/limit"
        },
        "hints": {
            "allow": [
                "GET"
            ],
            "formats": {
                "application/json": {}
            }
        }
    }`,
	}
	th.CheckDeepEquals(t, expected, *actual)
}

func TestPing(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandlePingSuccessfully(t)

	err := Ping(fake.ServiceClient()).ExtractErr()
	th.CheckNoErr(t, err)
}
