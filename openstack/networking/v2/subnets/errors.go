package subnets

import "fmt"

func err(str string) error {
	return fmt.Errorf("%s", str)
}

var (
	ErrNetworkIDRequired = err("A network ID is required")
	ErrCIDRRequired      = err("A valid CIDR is required")
	ErrInvalidIPType     = err("An IP type must either be 4 or 6")
)
