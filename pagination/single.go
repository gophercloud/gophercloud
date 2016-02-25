package pagination

import (
	"fmt"
	"reflect"
)

// SinglePageBase may be embedded in a Page that contains all of the results from an operation at once.
type SinglePageBase PageResult

// NextPageURL always returns "" to indicate that there are no more pages to return.
func (current SinglePageBase) NextPageURL() (string, error) {
	return "", nil
}

func (current SinglePageBase) IsEmpty() (bool, error) {
	if b, ok := current.Body.([]interface{}); ok {
		return len(b) == 0, nil
	}
	return true, fmt.Errorf("Error while checking if SinglePageBase was empty: expected []interface type for Body bot got %+v", reflect.TypeOf(current.Body))
}

// GetBody returns the single page's body. This method is needed to satisfy the
// Page interface.
func (current SinglePageBase) GetBody() interface{} {
	return current.Body
}
