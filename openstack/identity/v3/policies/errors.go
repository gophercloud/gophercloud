package policies

import "fmt"

// StringFieldLengthExceedsLimit is returned by the
// ToPolicyCreateMap/ToPolicyUpdateMap methods when validation of
// a type does not pass
type StringFieldLengthExceedsLimit struct {
	Field string
	Limit int
}

func (e StringFieldLengthExceedsLimit) Error() string {
	return fmt.Sprintf("String length of field [%s] exceeds limit (%d)",
		e.Field, e.Limit,
	)
}
