package firewall_groups

import "fmt"

func err(str string) error {
	return fmt.Errorf("%s", str)
}

var (
	errPolicyRequired = err("A policy ID is required")
)
