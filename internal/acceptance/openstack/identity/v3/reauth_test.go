//go:build acceptance || identity || reauth

package v3

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/openstack"
	th "github.com/gophercloud/gophercloud/v2/testhelper"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/projects"
)

func TestReauthAuthResultDeadlock(t *testing.T) {
	clients.RequireAdmin(t)

	ao, err := openstack.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	ao.AllowReauth = true

	provider, err := openstack.AuthenticatedClient(context.TODO(), ao)
	th.AssertNoErr(t, err)

	provider.SetToken("this is not a valid token")

	client, err := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
	th.AssertNoErr(t, err)
	pages, err := projects.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	_, err = projects.ExtractProjects(pages)
	th.AssertNoErr(t, err)
}
