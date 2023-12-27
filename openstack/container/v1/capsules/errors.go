package capsules

import (
	"github.com/gophercloud/gophercloud"
)

type ErrInvalidDataFormat struct {
	gophercloud.BaseError
}

func (e ErrInvalidDataFormat) Error() string {
	return "Data in neither json nor yaml format."
}
