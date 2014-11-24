package defsecrules

import (
	"testing"

	"github.com/rackspace/gophercloud/openstack/compute/v2/extensions/secgroups"
	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockListRulesResponse(t)

	count := 0

	err := List(client.ServiceClient()).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := ExtractDefaultRules(page)
		th.AssertNoErr(t, err)

		expected := []DefaultRule{
			DefaultRule{
				FromPort:   80,
				ID:         "f9a97fcf-3a97-47b0-b76f-919136afb7ed",
				IPProtocol: "TCP",
				IPRange:    secgroups.IPRange{CIDR: "10.10.10.0/24"},
				ToPort:     80,
			},
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, count)
}
