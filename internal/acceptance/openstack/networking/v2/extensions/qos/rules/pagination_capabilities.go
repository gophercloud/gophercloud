package rules

import "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/qos/ruletypes"

// supportsRuleDirections requires one driver to advertise all requested
// directions. Combining directions from different drivers would not establish
// that a single backend supports the fixture.
func supportsRuleDirections(ruleType ruletypes.RuleType, directions []string) bool {
	for _, driver := range ruleType.Drivers {
		for _, parameter := range driver.SupportedParameters {
			if parameter.ParameterName != "direction" || parameter.ParameterType != "choices" {
				continue
			}
			choices, ok := parameter.ParameterValues.([]any)
			if !ok {
				continue
			}
			available := make(map[string]bool)
			for _, choice := range choices {
				if direction, ok := choice.(string); ok {
					available[direction] = true
				}
			}
			supported := true
			for _, direction := range directions {
				if !available[direction] {
					supported = false
					break
				}
			}
			if supported {
				return true
			}
		}
	}
	return false
}
