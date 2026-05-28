package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestErrUnexpectedResponseCode(t *testing.T) {
	err := gophercloud.ErrUnexpectedResponseCode{
		URL:            "http://example.com",
		Method:         "GET",
		Expected:       []int{200},
		Actual:         404,
		Body:           []byte("the response body"),
		ResponseHeader: nil,
	}

	th.AssertEquals(t, 404, err.GetStatusCode())
	th.AssertTrue(t, gophercloud.ResponseCodeIs(err, http.StatusNotFound))
	th.AssertFalse(t, gophercloud.ResponseCodeIs(err, http.StatusInternalServerError))

	//even if application code wraps our error, ResponseCodeIs() should still work
	errWrapped := fmt.Errorf("could not frobnicate the foobar: %w", err)
	th.AssertTrue(t, gophercloud.ResponseCodeIs(errWrapped, http.StatusNotFound))
	th.AssertFalse(t, gophercloud.ResponseCodeIs(errWrapped, http.StatusInternalServerError))
}
