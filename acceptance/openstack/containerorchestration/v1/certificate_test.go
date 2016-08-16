// +build acceptance containerorchestration

package v1

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/containerorchestration/v1/bays"
	"github.com/gophercloud/gophercloud/openstack/containerorchestration/v1/certificates"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestGenerateCredentialsBundle(t *testing.T) {
	Setup(t)
	defer Teardown()

	// Use the most recently created bay
	var bay *bays.Bay
	pager := bays.List(Client, bays.ListOpts{})
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		bayList, _ := bays.ExtractBays(page)
		bay = &bayList[len(bayList)-1]
		return true, nil
	})
	bay, _ = bays.Get(Client, bay.ID).Extract()

	bundle, err := certificates.CreateCredentialsBundle(Client, bay.ID)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, bay.ID, bundle.BayID)
	th.AssertEquals(t, bay.COEEndpoint, bundle.COEEndpoint)
	th.AssertEquals(t, true, bundle.PrivateKey.Bytes != nil)
	th.AssertEquals(t, true, bundle.Certificate.Bytes != nil)
	th.AssertEquals(t, true, bundle.CACertificate.Bytes != nil)
}

func waitForStatus(client *gophercloud.ServiceClient, bay *bays.Bay, status string) error {
	return tools.WaitFor(func() (bool, error) {
		latest, err := bays.Get(client, bay.ID).Extract()
		if err != nil {
			return false, err
		}

		if latest.Status == status {
			// Success!
			return true, nil
		}

		if strings.HasSuffix(latest.Status, "FAILED") {
			return false, fmt.Errorf("The bay is in the failed status. %s", latest.StatusReason)
		}

		return false, nil
	})
}
