package gophercloud

import (
	"fmt"
)

type AuthOptions struct {
	Username, Password, TenantId string
}

func Authenticate(provider string, options AuthOptions) (*int, error) {
	return nil, fmt.Errorf("Not implemented.")
}
