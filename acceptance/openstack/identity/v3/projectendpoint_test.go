//go:build acceptance
// +build acceptance

package v3

import (
	"testing"

	"github.com/bizflycloud/gophercloud/acceptance/clients"
	"github.com/bizflycloud/gophercloud/acceptance/tools"
	"github.com/bizflycloud/gophercloud/openstack/identity/v3/endpoints"
	"github.com/bizflycloud/gophercloud/openstack/identity/v3/extensions/projectendpoints"
	th "github.com/bizflycloud/gophercloud/testhelper"
)

func TestProjectEndpoints(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	// Create a project to assign endpoints.
	project, err := CreateProject(t, client, nil)
	th.AssertNoErr(t, err)
	defer DeleteProject(t, client, project.ID)

	tools.PrintResource(t, project)

	// Get an endpoint
	allEndpointsPages, err := endpoints.List(client, nil).AllPages()
	th.AssertNoErr(t, err)

	allEndpoints, err := endpoints.ExtractEndpoints(allEndpointsPages)
	th.AssertNoErr(t, err)
	th.AssertIntGreaterOrEqual(t, len(allEndpoints), 1)
	endpoint := allEndpoints[0]

	// Attach endpoint
	err = projectendpoints.Create(client, project.ID, endpoint.ID).Err
	th.AssertNoErr(t, err)

	// List endpoints
	allProjectEndpointsPages, err := projectendpoints.List(client, project.ID).AllPages()
	th.AssertNoErr(t, err)

	allProjectEndpoints, err := projectendpoints.ExtractEndpoints(allProjectEndpointsPages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, len(allProjectEndpoints))

	tools.PrintResource(t, allProjectEndpoints[0])

	// Detach endpoint
	err = projectendpoints.Delete(client, project.ID, endpoint.ID).Err
	th.AssertNoErr(t, err)

}
