package identity

import "fmt"

var ErrNotImplemented = fmt.Errorf("Identity feature not yet implemented")
var ErrEndpoint = fmt.Errorf("Improper or missing Identity endpoint")
var ErrCredentials = fmt.Errorf("Improper or missing Identity credentials")
