//go:build acceptance || networking || extensions

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/common/extensions"
)

func TestExtensionsList(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}

	allPages, err := extensions.List(client).AllPages(context.TODO())
	if err != nil {
		t.Fatalf("Unable to list extensions: %v", err)
	}

	allExtensions, err := extensions.ExtractExtensions(allPages)
	if err != nil {
		t.Fatalf("Unable to extract extensions: %v", err)
	}

	for _, extension := range allExtensions {
		tools.PrintResource(t, extension)
	}
}

func TestExtensionGet(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}

	extension, err := extensions.Get(context.TODO(), client, "router").Extract()
	if err != nil {
		t.Fatalf("Unable to get extension port-security: %v", err)
	}

	tools.PrintResource(t, extension)
}
