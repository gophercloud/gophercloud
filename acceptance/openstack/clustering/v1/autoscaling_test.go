// +build acceptance clustering autoscaling profiles

package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/openstack/clustering/v1/profiles"

	"github.com/gophercloud/gophercloud/acceptance/tools"
	th "github.com/gophercloud/gophercloud/testhelper"
)

var testName string

func TestAutoScaling(t *testing.T) {
	testName = tools.RandomString("TESTACC-", 8)
	profileCreate(t)
}

func profileCreate(t *testing.T) {
	client, err := clients.NewClusteringV1Client()
	if err != nil {
		t.Fatalf("Unable to create clustering client: %v", err)
	}

	networks := []map[string]interface{}{
		{"network": "sandbox-internal-net"},
	}

	props := map[string]interface{}{
		"name":            "centos-server",
		"flavor":          "t2.micro",
		"image":           "centos7.3-latest",
		"networks":        networks,
		"security_groups": "",
	}

	profileName := testName
	optsProfile := &profiles.CreateOpts{
		Metadata: map[string]interface{}{
			"foo":  "bar",
			"test": "123",
		},
		Name: profileName,
		Spec: profiles.Spec{
			Type:       "os.nova.server",
			Version:    "1.0",
			Properties: props,
		},
	}

	createResult := profiles.Create(client, optsProfile)
	th.AssertNoErr(t, createResult.Err)

	requestID := createResult.Header.Get("X-OpenStack-Request-Id")
	th.AssertEquals(t, true, requestID != "")

	profile, err := createResult.Extract()
	if err != nil {
		t.Fatalf("Unable to create profile %s: %v", profileName, err)
	} else {
		t.Logf("Profile created %v", profile)
	}

	th.AssertEquals(t, profileName, profile.Name)
	th.AssertEquals(t, "os.nova.server", profile.Spec.Type)
	th.AssertEquals(t, "1.0", profile.Spec.Version)
}
