// +build acceptance

package monitoring

import (
	"fmt"
	identity "github.com/rackspace/gophercloud/openstack/identity/v2"
	"github.com/rackspace/gophercloud/openstack/utils"
	"github.com/rackspace/gophercloud/rackspace/monitoring"
	"github.com/rackspace/gophercloud/rackspace/monitoring/notificationPlans"
	"testing"
)

func TestRBACPermissions(t *testing.T) {
	ao, err := utils.AuthOptions()
	if err != nil {
		t.Fatal(err)
	}

	ao.AllowReauth = true
	r, err := identity.Authenticate(ao)
	if err != nil {
		t.Fatal(err)
	}

	// Find the cloud monitoring API

	sc, err := identity.GetServiceCatalog(r)
	if err != nil {
		t.Fatal(err)
	}

	ces, err := sc.CatalogEntries()
	if err != nil {
		t.Fatal(err)
	}

	monUrl, err := findMonitoringEndpoint(ces)
	if err != nil {
		t.Fatal(err)
	}

	// Build ourselves an interface to cloud monitoring!

	np := notificationPlans.NewClient(monitoring.Options{
		Endpoint:       monUrl,
		AuthOptions:    ao,
		Authentication: r,
	})

	// Try to delete a bogus notification plan

	dr, err := np.Delete("ajkhdlkajhdflkajshdf")
	if err != nil {
		fmt.Printf("%#v\n", err)
	}
	fmt.Printf("%#v\n", dr)
}

func findMonitoringEndpoint(ces []identity.CatalogEntry) (string, error) {
	for _, ce := range ces {
		if ce.Type == "rax:monitor" {
			return ce.Endpoints[0].PublicURL, nil
		}
	}
	return "", fmt.Errorf("No monitoring API in the service catalog")
}
