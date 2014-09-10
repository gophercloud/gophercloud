package images

import (
	"fmt"
	"github.com/racker/perigee"
)

var ErrNotImplemented = fmt.Errorf("Images functionality not implemented.")

type ListResults map[string]interface{}
type ImageResults map[string]interface{}

func List(c *Client) (ListResults, error) {
	var lr ListResults

	h, err := c.getListHeaders()
	if err != nil {
		return nil, err
	}

	err = perigee.Get(c.getListUrl(), perigee.Options{
		Results:     &lr,
		MoreHeaders: h,
	})
	return lr, err
}
