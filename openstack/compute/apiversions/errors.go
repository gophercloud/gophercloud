package apiversions

import (
	"fmt"
)

// ErrVersionNotFound is the error when the requested API version
// could not be found.
type ErrVersionNotFound struct{}

func (e ErrVersionNotFound) Error() string {
	return fmt.Sprintf("Unable to find requested API version")
}
