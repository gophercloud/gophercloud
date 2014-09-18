package ports

import "fmt"

func err(str string) error {
	return fmt.Errorf("%s", str)
}

var ()
