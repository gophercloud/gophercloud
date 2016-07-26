// +build acceptance compute extensions

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/common/extensions"
)

func TestExtensionsList(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	allPages, err := extensions.List(client).AllPages()
	if err != nil {
		t.Fatalf("Unable to list extensions: %v", err)
	}

	allExtensions, err := extensions.ExtractExtensions(allPages)
	if err != nil {
		t.Fatalf("Unable to extract extensions: %v", err)
	}

	for _, extension := range allExtensions {
		printExtension(t, &extension)
	}
}

func TestExtensionGet(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	extension, err := extensions.Get(client, "os-admin-actions").Extract()
	if err != nil {
		t.Fatalf("Unable to get extension os-admin-actions: %v", err)
	}

	printExtension(t, extension)
}

func printExtension(t *testing.T, extension *extensions.Extension) {
	t.Logf("Name: %s", extension.Name)
	t.Logf("Namespace: %s", extension.Namespace)
	t.Logf("Alias: %s", extension.Alias)
	t.Logf("Description: %s", extension.Description)
	t.Logf("Updated: %s", extension.Updated)
	t.Logf("Links: %v", extension.Links)
}
