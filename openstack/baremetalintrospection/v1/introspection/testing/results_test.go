package testing

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/baremetalintrospection/v1/introspection"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestLLDPTLVErrors(t *testing.T) {
	badInputs := []string{
		"[1]",
		"[1, 2]",
		"[\"foo\", \"bar\"]",
	}

	for _, input := range badInputs {
		var output introspection.LLDPTLVType
		err := json.Unmarshal([]byte(input), &output)
		if err == nil {
			t.Fatalf("No JSON parse error for invalid LLDP TLV %s", input)
		}

		if !strings.Contains(err.Error(), "LLDP TLV") {
			t.Fatalf("Unexpected JSON parse error \"%s\" for invalid LLDP TLV %s", err, input)
		}
	}
}

func TestExtraHardware(t *testing.T) {
	var output introspection.ExtraHardwareDataType
	err := json.Unmarshal([]byte(IntrospectionExtraHardwareJSONSample), &output)
	if err != nil {
		t.Fatalf("Failed to unmarshal ExtraHardware data: %s", err)
	}

	th.CheckDeepEquals(t, IntrospectionExtraHardware, output)
}

func TestIntrospectionNUMA(t *testing.T) {
	var output introspection.Data
	err := json.Unmarshal([]byte(IntrospectionNUMADataJSONSample), &output)
	if err != nil {
		t.Fatalf("Failed to unmarshal NUMA Data: %s", err)
	}

	th.CheckDeepEquals(t, IntrospectionNUMA, output.NUMATopology)
}

func TestHostnameInInventory(t *testing.T) {
	var output introspection.Data
	err := json.Unmarshal([]byte(IntrospectionDataJSONSample), &output)
	if err != nil {
		t.Fatalf("Failed to unmarshal Inventory data: %s", err)
	}

	th.CheckDeepEquals(t, IntrospectionDataRes.Inventory.Hostname, "myawesomehost")
}
