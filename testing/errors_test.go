package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	th "github.com/gophercloud/gophercloud/testhelper"
)

var body = []byte("error message")

func returnsUnexpectedResp(code int) gophercloud.ErrUnexpectedResponseCode {
	return gophercloud.ErrUnexpectedResponseCode{
		URL:            "http://example.com",
		Method:         "GET",
		Expected:       []int{200},
		Actual:         code,
		Body:           body,
		ResponseHeader: nil,
	}
}

func TestGetResponseCode404(t *testing.T) {
	var err404 error = gophercloud.ErrDefault404{ErrUnexpectedResponseCode: returnsUnexpectedResp(404)}

	err, ok := err404.(gophercloud.StatusCodeError)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, err.GetStatusCode(), 404)
}

func TestGetResponseCode502(t *testing.T) {
	var err502 error = gophercloud.ErrDefault502{ErrUnexpectedResponseCode: returnsUnexpectedResp(502)}

	err, ok := err502.(gophercloud.StatusCodeError)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, err.GetStatusCode(), 502)
}

func TestGetResponseCode504(t *testing.T) {
	var err504 error = gophercloud.ErrDefault504{ErrUnexpectedResponseCode: returnsUnexpectedResp(504)}

	err, ok := err504.(gophercloud.StatusCodeError)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, err.GetStatusCode(), 504)
}

func TestGetResponseBody(t *testing.T) {
	var err404 error = gophercloud.ErrDefault404{ErrUnexpectedResponseCode: returnsUnexpectedResp(404)}

	err, ok := err404.(gophercloud.BodyError)
	th.AssertEquals(t, true, ok)
	th.AssertByteArrayEquals(t, err.GetBody(), body)
}
