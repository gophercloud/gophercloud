package cdnobjects

import (
	"testing"

	os "github.com/rackspace/gophercloud/openstack/objectstorage/v1/objects"
	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

func TestDeleteCDNObject(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	os.HandleDeleteObjectSuccessfully(t)

	_, err := Delete(fake.ServiceClient(), "testContainer", "testObject", nil).ExtractHeaders()
	th.AssertNoErr(t, err)

}
