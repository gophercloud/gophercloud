// +build acceptance rackspace objectstorage v1

package v1

import (
	"bytes"
	"testing"

	raxCDNContainers "github.com/rackspace/gophercloud/rackspace/objectstorage/v1/cdncontainers"
	raxCDNObjects "github.com/rackspace/gophercloud/rackspace/objectstorage/v1/cdnobjects"
	raxContainers "github.com/rackspace/gophercloud/rackspace/objectstorage/v1/containers"
	raxObjects "github.com/rackspace/gophercloud/rackspace/objectstorage/v1/objects"
	th "github.com/rackspace/gophercloud/testhelper"
)

func TestCDNObjects(t *testing.T) {
	raxClient, err := createClient(t, false)
	th.AssertNoErr(t, err)

	headers, err := raxContainers.Create(raxClient, "gophercloud-test", nil).ExtractHeaders()
	th.AssertNoErr(t, err)
	t.Logf("Headers from Create Container request: %+v\n", headers)
	defer func() {
		_, err := raxContainers.Delete(raxClient, "gophercloud-test").ExtractHeaders()
		th.AssertNoErr(t, err)
	}()

	headers, err = raxObjects.Create(raxClient, "gophercloud-test", "test-object", bytes.NewBufferString("gophercloud cdn test"), nil).ExtractHeaders()
	th.AssertNoErr(t, err)
	t.Logf("Headers from Create Object request: %+v\n", headers)
	defer func() {
		_, err := raxObjects.Delete(raxClient, "gophercloud-test", "test-object", nil).ExtractHeaders()
		th.AssertNoErr(t, err)
	}()

	raxCDNClient, err := createClient(t, true)
	th.AssertNoErr(t, err)

	headers, err = raxCDNContainers.Enable(raxCDNClient, "gophercloud-test", raxCDNContainers.EnableOpts{CDNEnabled: true, TTL: 900}).ExtractHeaders()
	th.AssertNoErr(t, err)
	t.Logf("Headers from Enable CDN Container request: %+v\n", headers)

	headers, err = raxCDNObjects.Delete(raxCDNClient, "gophercloud-test", "test-object", nil).ExtractHeaders()
	th.AssertNoErr(t, err)
	t.Logf("Headers from Delete CDN Object request: %+v\n", headers)
}
