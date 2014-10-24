// +build acceptance rackspace

package v1

import (
	"testing"

	raxAccounts "github.com/rackspace/gophercloud/rackspace/objectstorage/v1/accounts"
	th "github.com/rackspace/gophercloud/testhelper"
)

func TestAccounts(t *testing.T) {
	c, err := createClient(t, false)
	th.AssertNoErr(t, err)

	headers, err := raxAccounts.Update(c, raxAccounts.UpdateOpts{Metadata: map[string]string{"white": "mountains"}}).Extract()
	th.AssertNoErr(t, err)
	t.Logf("Headers from Update Account request: %+v\n", headers)
	defer func() {
		_, err := raxAccounts.Update(c, raxAccounts.UpdateOpts{Metadata: map[string]string{"white": ""}}).Extract()
		th.AssertNoErr(t, err)
		metadata, err := raxAccounts.Get(c).ExtractMetadata()
		th.AssertNoErr(t, err)
		t.Logf("Metadata from Get Account request (after update reverted): %+v\n", metadata)
		th.CheckEquals(t, metadata["White"], "")
	}()

	metadata, err := raxAccounts.Get(c).ExtractMetadata()
	th.AssertNoErr(t, err)
	t.Logf("Metadata from Get Account request (after update): %+v\n", metadata)

	th.CheckEquals(t, metadata["White"], "mountains")
}
