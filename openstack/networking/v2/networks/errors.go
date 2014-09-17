package networks

import "fmt"

func requiredAttr(attr string) error {
	return fmt.Errorf("You must specify %s for this resource", attr)
}

func err(str string) error {
	return fmt.Errorf("%s", str)
}

var (
	ErrNameRequired = requiredAttr("name")
	ErrNoURLFound   = err("Next URL could not be extracted from collection")
)
