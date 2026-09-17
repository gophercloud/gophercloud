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

	output.Compat()
	th.CheckDeepEquals(t, Inventory, output)
}

func TestFrequency(t *testing.T) {
	cases := []struct {
		Value    string
		Expected string
	}{
		{
			Value:    `"frequency": "42"`,
			Expected: "42",
		},
		{
			Value:    `"missing frequency": 42`,
			Expected: "",
		},
		{
			Value:    `"frequency": 2100.084`,
			Expected: "2100.084",
		},
		{
			Value:    `"frequency": ""`,
			Expected: "",
		},
		{
			Value:    `"frequency": null`,
			Expected: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.Value, func(t *testing.T) {
			var output inventory.InventoryType
			sample := strings.Replace(InventorySample,
				`"frequency": "2100.084"`, tc.Value, 1)
			if sample == InventorySample {
				t.Fatal("TestFrequency needs updating")
			}

			err := json.Unmarshal([]byte(sample), &output)
			if err != nil {
				t.Fatalf("Failed to unmarshal inventory: %s", err)
			}

			expected := Inventory
			expected.CPU.Frequency = tc.Expected
			expected.CPU.RealFrequency = json.Number(tc.Expected)
			th.CheckDeepEquals(t, expected, output)
		})
	}
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
