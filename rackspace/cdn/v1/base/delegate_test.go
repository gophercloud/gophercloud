package base

import (
	"testing"

	os "github.com/rackspace/gophercloud/openstack/cdn/v1/base"
	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

func TestGetHomeDocument(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	os.HandleGetSuccessfully(t)

	actual, err := Get(fake.ServiceClient()).Extract()
	th.CheckNoErr(t, err)

	expected := os.HomeDocument{
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
	os.HandlePingSuccessfully(t)

	err := Ping(fake.ServiceClient()).ExtractErr()
	th.CheckNoErr(t, err)
}
