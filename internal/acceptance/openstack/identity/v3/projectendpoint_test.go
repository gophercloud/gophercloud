//go:build acceptance || identity || projectendpoints

package v3

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/endpoints"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/projectendpoints"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
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
	allEndpointsPages, err := endpoints.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allEndpoints, err := endpoints.ExtractEndpoints(allEndpointsPages)
	th.AssertNoErr(t, err)
	th.AssertIntGreaterOrEqual(t, len(allEndpoints), 1)
	endpoint := allEndpoints[0]

	// Attach endpoint
	err = projectendpoints.Create(context.TODO(), client, project.ID, endpoint.ID).Err
	th.AssertNoErr(t, err)

	// List endpoints
	allProjectEndpointsPages, err := projectendpoints.List(client, project.ID).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allProjectEndpoints, err := projectendpoints.ExtractEndpoints(allProjectEndpointsPages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, len(allProjectEndpoints))

	tools.PrintResource(t, allProjectEndpoints[0])

	// Detach endpoint
	err = projectendpoints.Delete(context.TODO(), client, project.ID, endpoint.ID).Err
	th.AssertNoErr(t, err)

}
