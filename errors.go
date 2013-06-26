package gophercloud

import (
	"fmt"
)

var ErrNotImplemented = fmt.Errorf("Not implemented")
var ErrProvider = fmt.Errorf("Missing or incorrect provider")
var ErrCredentials = fmt.Errorf("Missing or incomplete credentials")
var ErrConfiguration = fmt.Errorf("Missing or incomplete configuration")
