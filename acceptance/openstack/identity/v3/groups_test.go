// +build acceptance

package v3

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/groups"
)

func TestGroupsList(t *testing.T) {
	client, err := clients.NewIdentityV3Client()
	if err != nil {
		t.Fatalf("Unable to obtain an identity client: %v", err)
	}

	listOpts := groups.ListOpts{
		DomainID: "default",
	}

	allPages, err := groups.List(client, listOpts).AllPages()
	if err != nil {
		t.Fatalf("Unable to list groups: %v", err)
	}

	allGroups, err := groups.ExtractGroups(allPages)
	if err != nil {
		t.Fatalf("Unable to extract groups: %v", err)
	}

	for _, group := range allGroups {
		tools.PrintResource(t, group)
		tools.PrintResource(t, group.Extra)
	}
}

func TestGroupsGet(t *testing.T) {
	client, err := clients.NewIdentityV3Client()
	if err != nil {
		t.Fatalf("Unable to obtain an identity client: %v", err)
	}

	allPages, err := groups.List(client, nil).AllPages()
	if err != nil {
		t.Fatalf("Unable to list groups: %v", err)
	}

	allGroups, err := groups.ExtractGroups(allPages)
	if err != nil {
		t.Fatalf("Unable to extract groups: %v", err)
	}

	group := allGroups[0]
	p, err := groups.Get(client, group.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get group: %v", err)
	}

	tools.PrintResource(t, p)
}
