package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/sharenetworks"
)

func TestShareNetworkCreateDestroy(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create shared file system client: %v", err)
	}

	shareNetwork, err := CreateShareNetwork(t, client)
	if err != nil {
		t.Fatalf("Unable to create share network: %v", err)
	}

	PrintShareNetwork(t, shareNetwork)

	defer DeleteShareNetwork(t, client, shareNetwork)
}

func TestShareNetworkList(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create a shared file system client: %v", err)
	}

	allPages, err := sharenetworks.List(client, sharenetworks.ListOpts{}).AllPages()
	if err != nil {
		t.Fatalf("Unable to retrieve share networks: %v", err)
	}

	allShareNetworks, err := sharenetworks.ExtractShareNetworks(allPages)
	if err != nil {
		t.Fatalf("Unable to extract share networks: %v", err)
	}

	for _, shareNetwork := range allShareNetworks {
		PrintShareNetwork(t, &shareNetwork)
	}
}

// The test creates 2 shared networks and verifies that only the one(s) with
// a particular name are being listed
func TestShareNetworkListFiltering(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create a shared file system client: %v", err)
	}

	shareNetwork, err := CreateShareNetwork(t, client)
	if err != nil {
		t.Fatalf("Unable to create share network: %v", err)
	}
	defer DeleteShareNetwork(t, client, shareNetwork)

	shareNetwork, err = CreateShareNetwork(t, client)
	if err != nil {
		t.Fatalf("Unable to create share network: %v", err)
	}
	defer DeleteShareNetwork(t, client, shareNetwork)

	options := sharenetworks.ListOpts{
		Name: shareNetwork.Name,
	}

	allPages, err := sharenetworks.List(client, options).AllPages()
	if err != nil {
		t.Fatalf("Unable to retrieve share networks: %v", err)
	}

	allShareNetworks, err := sharenetworks.ExtractShareNetworks(allPages)
	if err != nil {
		t.Fatalf("Unable to extract share networks: %v", err)
	}

	for _, listedShareNetwork := range allShareNetworks {
		if listedShareNetwork.Name != shareNetwork.Name {
			t.Fatalf("The name of the share network was expected to be %s", shareNetwork.Name)
		}
		PrintShareNetwork(t, &listedShareNetwork)
	}
}
