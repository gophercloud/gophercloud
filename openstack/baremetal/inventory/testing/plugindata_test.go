package testing

import (
	"encoding/json"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/baremetal/inventory"
	"github.com/gophercloud/gophercloud/openstack/baremetalintrospection/v1/introspection"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestExtraHardware(t *testing.T) {
	var output inventory.ExtraDataType
	err := json.Unmarshal([]byte(ExtraDataJSONSample), &output)
	if err != nil {
		t.Fatalf("Failed to unmarshal ExtraHardware data: %s", err)
	}

	th.CheckDeepEquals(t, ExtraData, output)
}

func TestIntrospectionNUMA(t *testing.T) {
	var output introspection.Data
	err := json.Unmarshal([]byte(NUMADataJSONSample), &output)
	if err != nil {
		t.Fatalf("Failed to unmarshal NUMA Data: %s", err)
	}

	th.CheckDeepEquals(t, NUMATopology, output.NUMATopology)
}
