package clouds

import (
	"fmt"
)

type ErrFileNotFound struct {
	file            string
	searchLocations []string
}

func (e ErrFileNotFound) Error() string {
	return fmt.Sprintf(
		"%s file not found. Search locations were: %v",
		e.file, e.searchLocations)
}
