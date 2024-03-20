//go:build acceptance || networking || apiversion

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/apiversions"
)

func TestAPIVersionsList(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}

	allPages, err := apiversions.ListVersions(client).AllPages(context.TODO())
	if err != nil {
		t.Fatalf("Unable to list api versions: %v", err)
	}

	allAPIVersions, err := apiversions.ExtractAPIVersions(allPages)
	if err != nil {
		t.Fatalf("Unable to extract api versions: %v", err)
	}

	for _, apiVersion := range allAPIVersions {
		tools.PrintResource(t, apiVersion)
	}
}

func TestAPIResourcesList(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}

	allPages, err := apiversions.ListVersionResources(client, "v2.0").AllPages(context.TODO())
	if err != nil {
		t.Fatalf("Unable to list api version resources: %v", err)
	}

	allVersionResources, err := apiversions.ExtractVersionResources(allPages)
	if err != nil {
		t.Fatalf("Unable to extract version resources: %v", err)
	}

	for _, versionResource := range allVersionResources {
		tools.PrintResource(t, versionResource)
	}
}
