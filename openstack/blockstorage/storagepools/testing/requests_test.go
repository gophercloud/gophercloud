package testing

import (
	"github.com/gophercloud/gophercloud/openstack/blockstorage/storagepools"
	"github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
	"testing"
)

func TestListStoragePoolsDetail(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()
	HandleStoragePoolsListSuccessfully(t)

	r := storagepools.List(client.ServiceClient(), storagepools.ListOpts{Detail: true})
	actual, err := storagepools.ExtractStoragePools(r)
	testhelper.AssertNoErr(t, err)

	if len(actual) != 2 {
		t.Fatalf("Expected 2 backends,  got %d", len(actual))
	}
	testhelper.CheckDeepEquals(t, StoragePoolFake1, actual[0])
	testhelper.CheckDeepEquals(t, StoragePoolFake2, actual[1])
}
