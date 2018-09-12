package gophercloud

import (
	"testing"
	"time"

	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestBuildRequestBody(t *testing.T) {
	type Params struct {
		ID   string    `json:"id"`
		Date time.Time `json:"date"`
	}

	params := &Params{
		ID:   "Foo",
		Date: time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	body, err := BuildRequestBody(params, "params")
	th.AssertNoErr(t, err)

	th.AssertJSONEquals(t, `{"params": {"id": "Foo", "date": "1970-01-01T00:00:00Z"}}`, body)
}
