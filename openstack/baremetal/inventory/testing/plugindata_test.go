package testing

import (
	"encoding/json"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/inventory"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
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
	var output inventory.StandardPluginData
	err := json.Unmarshal([]byte(NUMADataJSONSample), &output)
	if err != nil {
		t.Fatalf("Failed to unmarshal NUMA Data: %s", err)
	}

	th.CheckDeepEquals(t, NUMATopology, output.NUMATopology)
}

func TestStandardPluginData(t *testing.T) {
	var output inventory.StandardPluginData

	err := json.Unmarshal([]byte(StandardPluginDataSample), &output)
	if err != nil {
		t.Fatalf("Failed to unmarshal plugin data: %s", err)
	}

	th.CheckDeepEquals(t, StandardPluginData, output)
}
