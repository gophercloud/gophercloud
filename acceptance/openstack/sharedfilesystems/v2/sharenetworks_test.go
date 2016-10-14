package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
)

func TestShareNetworkCreate(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create shared file system client: %v", err)
	}

	shareNetwork, err := CreateShareNetwork(t, client)
	if err != nil {
		t.Fatalf("Unable to create share network: %v", err)
	}

	// TODO: delete the share network when the delete is implemented

	PrintShareNetwork(t, shareNetwork)
}
