package testing

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/inventory"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestInventory(t *testing.T) {
	var output inventory.InventoryType
	err := json.Unmarshal([]byte(InventorySample), &output)
	if err != nil {
		t.Fatalf("Failed to unmarshal inventory: %s", err)
	}

	th.CheckDeepEquals(t, Inventory, output)
}

func TestLLDPTLVErrors(t *testing.T) {
	badInputs := []string{
		"[1]",
		"[1, 2]",
		"[\"foo\", \"bar\"]",
	}

	for _, input := range badInputs {
		var output inventory.LLDPTLVType
		err := json.Unmarshal([]byte(input), &output)
		if err == nil {
			t.Fatalf("No JSON parse error for invalid LLDP TLV %s", input)
		}

		if !strings.Contains(err.Error(), "LLDP TLV") {
			t.Fatalf("Unexpected JSON parse error \"%s\" for invalid LLDP TLV %s", err, input)
		}
	}
}
