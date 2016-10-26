package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
	"github.com/gophercloud/gophercloud/pagination"
)

func TestShareCreate(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create a sharedfs client: %v", err)
	}

	share, err := CreateShare(t, client, "my test share")
	if err != nil {
		t.Fatalf("Unable to create a share: %v", err)
	}

	defer DeleteShare(t, client, share)

	created, err := shares.Get(client, share.ID).Extract()
	if err != nil {
		t.Errorf("Unable to retrieve share: %v", err)
	}
	PrintShare(t, created)
}

func TestShareListShort(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create a sharedfs client: %v", err)
	}

	names := []string{"share_one", "share_two", "share_three"}

	created, err := CreateShares(t, client, names)
	if err != nil {
		t.Fatalf("Unable to create a share: %v", err)
	}

	defer DeleteShares(t, client, created)

	pages, err := shares.List(client, &shares.ListOpts{}, false).AllPages()
	if err != nil {
		t.Fatalf("Failed to list shares: %v", err)
	}

	shares, err := shares.ExtractShares(pages)
	if err != nil {
		t.Fatalf("Unable to extract shares: %v", err)
	}

	for _, share := range shares {
		PrintShare(t, &share)
	}
}

func TestShareListDetail(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create a sharedfs client: %v", err)
	}

	names := []string{"share_one", "share_two", "share_three"}

	created, err := CreateShares(t, client, names)
	if err != nil {
		t.Fatalf("Unable to create a share: %v", err)
	}

	defer DeleteShares(t, client, created)

	pages, err := shares.List(client, &shares.ListOpts{}, true).AllPages()
	if err != nil {
		t.Fatalf("Failed to list shares: %v", err)
	}

	shares, err := shares.ExtractShares(pages)
	if err != nil {
		t.Fatalf("Unable to extract shares: %v", err)
	}

	for _, share := range shares {
		PrintShare(t, &share)
	}
}

func TestShareListDetailPaginateLimit(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create a sharedfs client: %v", err)
	}

	names := []string{"share_one", "share_two", "share_three",
		"share_four", "share_five", "share_six"}

	created, err := CreateShares(t, client, names)
	if err != nil {
		t.Fatalf("Unable to create a share: %v", err)
	}

	defer DeleteShares(t, client, created)
	limit := 2

	// Lists two shares per page
	err = shares.List(client, &shares.ListOpts{Offset: 0, Limit: limit}, true).EachPage(
		func(page pagination.Page) (bool, error) {
			l, err := shares.ExtractShares(page)
			if err != nil {
				t.Fatalf("Failed to extract shares: %v", err)
			}
			if len(l) != limit {
				t.Fatalf("Unexpected number of shares: %d", len(l))
			}
			t.Logf("Got %d shares for this page", len(l))
			for _, share := range l {
				PrintShare(t, &share)
			}
			return true, nil
		})
	if err != nil {
		t.Fatalf("Unable to retrieve shares: %v", err)
	}
}

// Lists two shares per page, starting at offset 2. This
// should result in a total of four shares
func TestShareListDetailPaginateLimitOffset(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create a sharedfs client: %v", err)
	}

	names := []string{"share_one", "share_two", "share_three",
		"share_four", "share_five", "share_six"}

	created, err := CreateShares(t, client, names)
	if err != nil {
		t.Fatalf("Unable to create a share: %v", err)
	}

	defer DeleteShares(t, client, created)
	limit := 2
	count := 0

	// Lists two shares per page
	err = shares.List(client, &shares.ListOpts{Offset: limit, Limit: limit}, true).EachPage(
		func(page pagination.Page) (bool, error) {
			l, err := shares.ExtractShares(page)
			if err != nil {
				t.Fatalf("Failed to extract shares: %v", err)
			}
			if len(l) != limit {
				t.Fatalf("Unexpected number of shares: %d", len(l))
			}

			count++

			t.Logf("Got %d shares for this page", len(l))
			for _, share := range l {
				PrintShare(t, &share)
			}
			return true, nil
		})
	if err != nil {
		t.Fatalf("Unable to retrieve shares: %v", err)
	}
	if count != 4 {
		t.Fatalf("Unexpected share count: %d", count)
	}
}
