package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/clustering/v1/actions"
)

func TestActionsList(t *testing.T) {
	client, err := clients.NewClusteringV1Client()
	if err != nil {
		t.Fatalf("Unable to create a clustering client: %v", err)
	}

	allPages, err := actions.List(client, nil).AllPages()
	if err != nil {
		t.Fatalf("Unable to list actions info: %v", err)
	}

	actionInfos, err := actions.ExtractActions(allPages)
	if err != nil {
		t.Fatalf("Unable to extract actions info: %v", err)
	}

	for _, actionInfo := range actionInfos {
		tools.PrintResource(t, actionInfo)
	}
}
