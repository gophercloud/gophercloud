package networks

import "fmt"

func requiredAttr(attr string) error {
	return fmt.Errorf("You must specify %s for this resource", attr)
}

var (
	ErrNameRequired = requiredAttr("name")
)
