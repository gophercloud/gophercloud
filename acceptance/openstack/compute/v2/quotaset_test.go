// +build acceptance compute quotasets

package v2

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/quotasets"
	"github.com/gophercloud/gophercloud/openstack/identity/v2/tenants"
)

func TestQuotasetGet(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	identityClient, err := newIdentityClient()
	if err != nil {
		t.Fatalf("Unable to get a new identity client: %v", err)
	}

	tenantID, err := getTenantID(t, identityClient)
	if err != nil {
		t.Fatal(err)
	}

	quotaSet, err := quotasets.Get(client, tenantID).Extract()
	if err != nil {
		t.Fatal(err)
	}

	printQuotaSet(t, quotaSet)
}

func getTenantID(t *testing.T, client *gophercloud.ServiceClient) (string, error) {
	allPages, err := tenants.List(client, nil).AllPages()
	if err != nil {
		t.Fatalf("Unable to get list of tenants: %v", err)
	}

	allTenants, err := tenants.ExtractTenants(allPages)
	if err != nil {
		t.Fatalf("Unable to extract tenants: %v", err)
	}

	for _, tenant := range allTenants {
		return tenant.ID, nil
	}

	return "", fmt.Errorf("Unable to get tenant ID")
}

func printQuotaSet(t *testing.T, quotaSet *quotasets.QuotaSet) {
	t.Logf("instances: %d\n", quotaSet.Instances)
	t.Logf("cores: %d\n", quotaSet.Cores)
	t.Logf("ram: %d\n", quotaSet.Ram)
	t.Logf("key_pairs: %d\n", quotaSet.KeyPairs)
	t.Logf("metadata_items: %d\n", quotaSet.MetadataItems)
	t.Logf("security_groups: %d\n", quotaSet.SecurityGroups)
	t.Logf("security_group_rules: %d\n", quotaSet.SecurityGroupRules)
	t.Logf("fixed_ips: %d\n", quotaSet.FixedIps)
	t.Logf("floating_ips: %d\n", quotaSet.FloatingIps)
	t.Logf("injected_file_content_bytes: %d\n", quotaSet.InjectedFileContentBytes)
	t.Logf("injected_file_path_bytes: %d\n", quotaSet.InjectedFilePathBytes)
	t.Logf("injected_files: %d\n", quotaSet.InjectedFiles)
}
