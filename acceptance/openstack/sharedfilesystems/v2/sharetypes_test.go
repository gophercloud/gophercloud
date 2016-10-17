package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
)

func TestShareTypeCreateDestroy(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create shared file system client: %v", err)
	}

	shareType, err := CreateShareType(t, client)
	if err != nil {
		t.Fatalf("Unable to create share type: %v", err)
	}

	PrintShareType(t, shareType)

	defer DeleteShareType(t, client, shareType)
}
