// +build acceptance rackspace objectstorage v1

package v1

import (
	"testing"

	osContainers "github.com/rackspace/gophercloud/openstack/objectstorage/v1/containers"
	"github.com/rackspace/gophercloud/pagination"
	raxCDNContainers "github.com/rackspace/gophercloud/rackspace/objectstorage/v1/cdncontainers"
	raxContainers "github.com/rackspace/gophercloud/rackspace/objectstorage/v1/containers"
	th "github.com/rackspace/gophercloud/testhelper"
)

func TestCDNContainers(t *testing.T) {
	raxClient, err := createClient(t, false)
	th.AssertNoErr(t, err)

	headers, err := raxContainers.Create(raxClient, "gophercloud-test", nil).ExtractHeaders()
	th.AssertNoErr(t, err)
	t.Logf("Headers from Create Container request: %+v\n", headers)
	defer func() {
		_, err := raxContainers.Delete(raxClient, "gophercloud-test").ExtractHeaders()
		th.AssertNoErr(t, err)
	}()

	raxCDNClient, err := createClient(t, true)
	th.AssertNoErr(t, err)

	headers, err = raxCDNContainers.Enable(raxCDNClient, "gophercloud-test", raxCDNContainers.EnableOpts{CDNEnabled: true, TTL: 900}).ExtractHeaders()
	th.AssertNoErr(t, err)
	t.Logf("Headers from Enable CDN Container request: %+v\n", headers)

	t.Logf("Container Names available to the currently issued token:")
	count := 0
	err = raxCDNContainers.List(raxCDNClient, &osContainers.ListOpts{Full: false}).EachPage(func(page pagination.Page) (bool, error) {
		t.Logf("--- Page %02d ---", count)

		names, err := raxCDNContainers.ExtractNames(page)
		th.AssertNoErr(t, err)

		for i, name := range names {
			t.Logf("[%02d] %s", i, name)
		}

		count++
		return true, nil
	})
	th.AssertNoErr(t, err)
	if count == 0 {
		t.Errorf("No containers listed for your current token.")
	}

	headers, err = raxCDNContainers.Update(raxCDNClient, "gophercloud-test", raxCDNContainers.UpdateOpts{CDNEnabled: false}).ExtractHeaders()
	th.AssertNoErr(t, err)
	t.Logf("Headers from Update CDN Container request: %+v\n", headers)

	headers, err = raxCDNContainers.Get(raxCDNClient, "gophercloud-test").ExtractHeaders()
	th.AssertNoErr(t, err)
	t.Logf("Headers from Get CDN Container request (after update): %+v\n", headers)
}
