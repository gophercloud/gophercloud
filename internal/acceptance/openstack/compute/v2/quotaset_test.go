//go:build acceptance || compute || quotasets

package v2

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/quotasets"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/projects"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestQuotasetGet(t *testing.T) {
	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	identityClient, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	projectID, err := getProjectID(t, identityClient)
	th.AssertNoErr(t, err)

	quotaSet, err := quotasets.Get(context.TODO(), client, projectID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, quotaSet)

	th.AssertEquals(t, quotaSet.FixedIPs, -1)
}

func getProjectID(t *testing.T, client *gophercloud.ServiceClient) (string, error) {
	allPages, err := projects.ListAvailable(client).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allProjects, err := projects.ExtractProjects(allPages)
	th.AssertNoErr(t, err)

	for _, project := range allProjects {
		return project.ID, nil
	}

	return "", fmt.Errorf("Unable to get project ID")
}

func getProjectIDByName(t *testing.T, client *gophercloud.ServiceClient, name string) (string, error) {
	allPages, err := projects.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allProjects, err := projects.ExtractProjects(allPages)
	th.AssertNoErr(t, err)

	for _, project := range allProjects {
		if project.Name == name {
			return project.ID, nil
		}
	}

	return "", fmt.Errorf("Unable to get project ID")
}

// What will be sent as desired Quotas to the Server
var UpdateQuotaOpts = quotasets.UpdateOpts{
	FixedIPs:                 gophercloud.IntToPointer(10),
	FloatingIPs:              gophercloud.IntToPointer(10),
	InjectedFileContentBytes: gophercloud.IntToPointer(10240),
	InjectedFilePathBytes:    gophercloud.IntToPointer(255),
	InjectedFiles:            gophercloud.IntToPointer(5),
	KeyPairs:                 gophercloud.IntToPointer(10),
	MetadataItems:            gophercloud.IntToPointer(128),
	RAM:                      gophercloud.IntToPointer(20000),
	SecurityGroupRules:       gophercloud.IntToPointer(20),
	SecurityGroups:           gophercloud.IntToPointer(10),
	Cores:                    gophercloud.IntToPointer(10),
	Instances:                gophercloud.IntToPointer(4),
	ServerGroups:             gophercloud.IntToPointer(2),
	ServerGroupMembers:       gophercloud.IntToPointer(3),
}

// What the Server hopefully returns as the new Quotas
var UpdatedQuotas = quotasets.QuotaSet{
	FixedIPs:                 10,
	FloatingIPs:              10,
	InjectedFileContentBytes: 10240,
	InjectedFilePathBytes:    255,
	InjectedFiles:            5,
	KeyPairs:                 10,
	MetadataItems:            128,
	RAM:                      20000,
	SecurityGroupRules:       20,
	SecurityGroups:           10,
	Cores:                    10,
	Instances:                4,
	ServerGroups:             2,
	ServerGroupMembers:       3,
}

func TestQuotasetUpdateDelete(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	idclient, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	projectid, err := getProjectIDByName(t, idclient, os.Getenv("OS_PROJECT_NAME"))
	th.AssertNoErr(t, err)

	// save original quotas
	orig, err := quotasets.Get(context.TODO(), client, projectid).Extract()
	th.AssertNoErr(t, err)

	// Test Update
	res, err := quotasets.Update(context.TODO(), client, projectid, UpdateQuotaOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, UpdatedQuotas, *res)

	// Test Delete
	_, err = quotasets.Delete(context.TODO(), client, projectid).Extract()
	th.AssertNoErr(t, err)

	// We dont know the default quotas, so just check if the quotas are not the same as before
	newres, err := quotasets.Get(context.TODO(), client, projectid).Extract()
	th.AssertNoErr(t, err)
	if newres.RAM == res.RAM {
		t.Fatalf("Failed to update quotas")
	}

	restore := quotasets.UpdateOpts{}
	FillUpdateOptsFromQuotaSet(*orig, &restore)

	// restore original quotas
	res, err = quotasets.Update(context.TODO(), client, projectid, restore).Extract()
	th.AssertNoErr(t, err)

	orig.ID = ""
	th.AssertDeepEquals(t, orig, res)
}
