package rules

import (
	"encoding/json"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/qos/ruletypes"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestSupportsRuleDirections(t *testing.T) {
	for _, tc := range []struct {
		name      string
		drivers   string
		supported bool
	}{
		{
			name:      "one driver supports both",
			drivers:   `[{"supported_parameters":[{"parameter_name":"direction","parameter_type":"choices","parameter_values":["ingress","egress"]}]}]`,
			supported: true,
		},
		{
			name:      "later driver supports both",
			drivers:   `[{"supported_parameters":[]},{"supported_parameters":[{"parameter_name":"direction","parameter_type":"choices","parameter_values":["ingress","egress"]}]}]`,
			supported: true,
		},
		{
			name:    "directions split across drivers",
			drivers: `[{"supported_parameters":[{"parameter_name":"direction","parameter_type":"choices","parameter_values":["ingress"]}]},{"supported_parameters":[{"parameter_name":"direction","parameter_type":"choices","parameter_values":["egress"]}]}]`,
		},
		{
			name:    "egress only",
			drivers: `[{"supported_parameters":[{"parameter_name":"direction","parameter_type":"choices","parameter_values":["egress"]}]}]`,
		},
		{
			name:    "wrong parameter",
			drivers: `[{"supported_parameters":[{"parameter_name":"other","parameter_type":"choices","parameter_values":["ingress","egress"]}]}]`,
		},
		{
			name:    "wrong parameter type",
			drivers: `[{"supported_parameters":[{"parameter_name":"direction","parameter_type":"range","parameter_values":["ingress","egress"]}]}]`,
		},
		{
			name:    "malformed choices",
			drivers: `[{"supported_parameters":[{"parameter_name":"direction","parameter_type":"choices","parameter_values":"ingress,egress"}]}]`,
		},
		{
			name:    "empty choices",
			drivers: `[{"supported_parameters":[{"parameter_name":"direction","parameter_type":"choices","parameter_values":[]}]}]`,
		},
		{name: "no drivers", drivers: `[]`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var ruleType ruletypes.RuleType
			th.AssertNoErr(t, json.Unmarshal([]byte(`{"type":"minimum_bandwidth","drivers":`+tc.drivers+`}`), &ruleType))
			th.AssertEquals(t, tc.supported, supportsRuleDirections(ruleType, []string{"ingress", "egress"}))
		})
	}
}
