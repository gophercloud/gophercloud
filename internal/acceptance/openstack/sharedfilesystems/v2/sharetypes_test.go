//go:build acceptance || sharedfilesystems || sharetypes

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/sharetypes"
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

	tools.PrintResource(t, shareType)

	defer DeleteShareType(t, client, shareType)
}

func TestShareTypeList(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create a shared file system client: %v", err)
	}

	allPages, err := sharetypes.List(client, sharetypes.ListOpts{}).AllPages(context.TODO())
	if err != nil {
		t.Fatalf("Unable to retrieve share types: %v", err)
	}

	allShareTypes, err := sharetypes.ExtractShareTypes(allPages)
	if err != nil {
		t.Fatalf("Unable to extract share types: %v", err)
	}

	for _, shareType := range allShareTypes {
		tools.PrintResource(t, &shareType)
	}
}

func TestShareTypeGetDefault(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create a shared file system client: %v", err)
	}

	shareType, err := sharetypes.GetDefault(context.TODO(), client).Extract()
	if err != nil {
		t.Fatalf("Unable to retrieve the default share type: %v", err)
	}

	tools.PrintResource(t, shareType)
}

func TestShareTypeExtraSpecs(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create shared file system client: %v", err)
	}

	shareType, err := CreateShareType(t, client)
	if err != nil {
		t.Fatalf("Unable to create share type: %v", err)
	}

	options := sharetypes.SetExtraSpecsOpts{
		ExtraSpecs: map[string]any{"my_new_key": "my_value"},
	}

	_, err = sharetypes.SetExtraSpecs(context.TODO(), client, shareType.ID, options).Extract()
	if err != nil {
		t.Fatalf("Unable to set extra specs for Share type: %s", shareType.Name)
	}

	extraSpecs, err := sharetypes.GetExtraSpecs(context.TODO(), client, shareType.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to retrieve share type: %s", shareType.Name)
	}

	if extraSpecs["driver_handles_share_servers"] != "True" {
		t.Fatal("driver_handles_share_servers was expected to be true")
	}

	if extraSpecs["my_new_key"] != "my_value" {
		t.Fatal("my_new_key was expected to be equal to my_value")
	}

	err = sharetypes.UnsetExtraSpecs(context.TODO(), client, shareType.ID, "my_new_key").ExtractErr()
	if err != nil {
		t.Fatalf("Unable to unset extra specs for Share type: %s", shareType.Name)
	}

	extraSpecs, err = sharetypes.GetExtraSpecs(context.TODO(), client, shareType.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to retrieve share type: %s", shareType.Name)
	}

	if _, ok := extraSpecs["my_new_key"]; ok {
		t.Fatalf("my_new_key was expected to be unset for Share type: %s", shareType.Name)
	}

	tools.PrintResource(t, shareType)

	defer DeleteShareType(t, client, shareType)
}

func TestShareTypeAccess(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create shared file system client: %v", err)
	}

	shareType, err := CreateShareType(t, client)
	if err != nil {
		t.Fatalf("Unable to create share type: %v", err)
	}

	options := sharetypes.AccessOpts{
		Project: "9e3a5a44e0134445867776ef53a37605",
	}

	err = sharetypes.AddAccess(context.TODO(), client, shareType.ID, options).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to add a new access to a share type: %v", err)
	}

	access, err := sharetypes.ShowAccess(context.TODO(), client, shareType.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to retrieve the access details for a share type: %v", err)
	}

	expected := []sharetypes.ShareTypeAccess{{ShareTypeID: shareType.ID, ProjectID: options.Project}}

	if access[0] != expected[0] {
		t.Fatal("Share type access is not the same than expected")
	}

	err = sharetypes.RemoveAccess(context.TODO(), client, shareType.ID, options).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to remove an access from a share type: %v", err)
	}

	access, err = sharetypes.ShowAccess(context.TODO(), client, shareType.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to retrieve the access details for a share type: %v", err)
	}

	if len(access) > 0 {
		t.Fatalf("No access should be left for the share type: %s", shareType.Name)
	}

	tools.PrintResource(t, shareType)

	defer DeleteShareType(t, client, shareType)

}
