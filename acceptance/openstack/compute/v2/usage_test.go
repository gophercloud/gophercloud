// +build acceptance compute usage

package v2

import (
	"strings"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/usage"
)

func TestUsageSingleTenant(t *testing.T) {
	client, err := clients.NewComputeV2Client()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	server, err := CreateServer(t, client)
	if err != nil {
		t.Fatalf("Unable to create server: %v", err)
	}
	DeleteServer(t, client, server)

	endpointParts := strings.Split(client.Endpoint, "/")
	tenantID := endpointParts[4]

	end := time.Now()
	start := end.AddDate(0, -1, 0)
	opts := usage.SingleTenantOpts{
		Start: &start,
		End:   &end,
	}

	page, err := usage.SingleTenant(client, tenantID, opts).AllPages()
	if err != nil {
		t.Fatal(err)
	}

	tenantUsage, err := usage.ExtractSingleTenant(page)
	if err != nil {
		t.Fatal(err)
	}

	tools.PrintResource(t, tenantUsage)

	if tenantUsage.TotalHours == 0 {
		t.Fatalf("TotalHours should not be 0")
	}
}
