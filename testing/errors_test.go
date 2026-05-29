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

	th.AssertEquals(t, err.GetStatusCode(), 404)
	th.AssertEquals(t, gophercloud.ResponseCodeIs(err, http.StatusNotFound), true)
	th.AssertEquals(t, gophercloud.ResponseCodeIs(err, http.StatusInternalServerError), false)

	//even if application code wraps our error, ResponseCodeIs() should still work
	errWrapped := fmt.Errorf("could not frobnicate the foobar: %w", err)
	th.AssertEquals(t, gophercloud.ResponseCodeIs(errWrapped, http.StatusNotFound), true)
	th.AssertEquals(t, gophercloud.ResponseCodeIs(errWrapped, http.StatusInternalServerError), false)
}
