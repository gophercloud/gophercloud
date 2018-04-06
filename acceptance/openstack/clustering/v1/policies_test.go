// +build acceptance clustering policies

package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/clustering/v1/policies"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestPolicyList(t *testing.T) {
	client, err := clients.NewClusteringV1Client()
	th.AssertNoErr(t, err)

	allPages, err := policies.List(client, nil).AllPages()
	th.AssertNoErr(t, err)

	allPolicies, err := policies.ExtractPolicies(allPages)
	th.AssertNoErr(t, err)

	for _, v := range allPolicies {
		tools.PrintResource(t, v)

		if v.CreatedAt.IsZero() {
			t.Fatalf("CreatedAt value should not be zero")
		}
		t.Log("Created at: " + v.CreatedAt.String())

		if !v.UpdatedAt.IsZero() {
			t.Log("Updated at: " + v.UpdatedAt.String())
		}
	}
}

func TestPolicyCreateAndDelete(t *testing.T) {
	client, err := clients.NewClusteringV1Client()
	th.AssertNoErr(t, err)

	opts := policies.CreateOpts{
		Name: "new_policy2",
		Spec: policies.Spec{
			Description: "new policy description",
			Properties: map[string]interface{}{
				"destroy_after_deletion":  true,
				"grace_period":            60,
				"reduce_desired_capacity": false,
				"criteria":                "OLDEST_FIRST",
			},
			Type:    "senlin.policy.deletion",
			Version: "1.0",
		},
	}

	createdPolicy, err := policies.Create(client, opts).Extract()
	th.AssertNoErr(t, err)

	defer policies.Delete(client, createdPolicy.ID)

	tools.PrintResource(t, createdPolicy)

	if createdPolicy.CreatedAt.IsZero() {
		t.Fatalf("CreatedAt value should not be zero")
	}
	t.Log("Created at: " + createdPolicy.CreatedAt.String())

	if !createdPolicy.UpdatedAt.IsZero() {
		t.Log("Updated at: " + createdPolicy.UpdatedAt.String())
	}
}
