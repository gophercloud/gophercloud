// +build acceptance quotasets

package v3

import (
	"os"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/extensions/quotasets"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestQuotasetGet(t *testing.T) {
	client, projectID := getClientAndProject(t)

	quotaSet, err := quotasets.Get(client, projectID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, quotaSet)
}

func TestQuotasetGetDefaults(t *testing.T) {
	client, projectID := getClientAndProject(t)

	quotaSet, err := quotasets.GetDefaults(client, projectID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, quotaSet)
}

// getClientAndProject reduces boilerplate by returning a new blockstorage v3
// ServiceClient and a project ID obtained from the OS_PROJECT_NAME envvar.
func getClientAndProject(t *testing.T) (*gophercloud.ServiceClient, string) {
	client, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	projectID := os.Getenv("OS_PROJECT_NAME")
	th.AssertNoErr(t, err)
	return client, projectID
}
