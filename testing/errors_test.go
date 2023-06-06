package testing

import (
	"errors"
	"testing"

	"github.com/gophercloud/gophercloud"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func returnsUnexpectedResp(code int) gophercloud.ErrUnexpectedResponseCode {
	return gophercloud.ErrUnexpectedResponseCode{
		URL:            "http://example.com",
		Method:         "GET",
		Expected:       []int{200},
		Actual:         code,
		Body:           []byte("the response body"),
		ResponseHeader: nil,
	}
}

func TestGetResponseCode404(t *testing.T) {
	var err404 error = gophercloud.ErrDefault404{ErrUnexpectedResponseCode: returnsUnexpectedResp(404)}

	err, ok := err404.(gophercloud.StatusCodeError)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, err.GetStatusCode(), 404)

	t.Run("wraps ErrUnexpectedResponseCode", func(t *testing.T) {
		var unexpectedResponseCode gophercloud.ErrUnexpectedResponseCode
		if errors.As(err, &unexpectedResponseCode) {
			if want, have := "the response body", string(unexpectedResponseCode.Body); want != have {
				t.Errorf("expected the wrapped error to contain the response body, found %q", have)
			}
		} else {
			t.Errorf("err.Unwrap() didn't return ErrUnexpectedResponseCode")
		}
	})
}

func TestGetResponseCode502(t *testing.T) {
	var err502 error = gophercloud.ErrDefault502{ErrUnexpectedResponseCode: returnsUnexpectedResp(502)}

	err, ok := err502.(gophercloud.StatusCodeError)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, err.GetStatusCode(), 502)

	t.Run("wraps ErrUnexpectedResponseCode", func(t *testing.T) {
		var unexpectedResponseCode gophercloud.ErrUnexpectedResponseCode
		if errors.As(err, &unexpectedResponseCode) {
			if want, have := "the response body", string(unexpectedResponseCode.Body); want != have {
				t.Errorf("expected the wrapped error to contain the response body, found %q", have)
			}
		} else {
			t.Errorf("err.Unwrap() didn't return ErrUnexpectedResponseCode")
		}
	})
}

func TestGetResponseCode504(t *testing.T) {
	var err504 error = gophercloud.ErrDefault504{ErrUnexpectedResponseCode: returnsUnexpectedResp(504)}

	err, ok := err504.(gophercloud.StatusCodeError)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, err.GetStatusCode(), 504)

	t.Run("wraps ErrUnexpectedResponseCode", func(t *testing.T) {
		var unexpectedResponseCode gophercloud.ErrUnexpectedResponseCode
		if errors.As(err, &unexpectedResponseCode) {
			if want, have := "the response body", string(unexpectedResponseCode.Body); want != have {
				t.Errorf("expected the wrapped error to contain the response body, found %q", have)
			}
		} else {
			t.Errorf("err.Unwrap() didn't return ErrUnexpectedResponseCode")
		}
	})
}

// Compile-time check that all response-code errors implement `Unwrap()`
type unwrapper interface {
	Unwrap() error
}

var (
	_ unwrapper = gophercloud.ErrDefault401{}
	_ unwrapper = gophercloud.ErrDefault403{}
	_ unwrapper = gophercloud.ErrDefault404{}
	_ unwrapper = gophercloud.ErrDefault405{}
	_ unwrapper = gophercloud.ErrDefault408{}
	_ unwrapper = gophercloud.ErrDefault409{}
	_ unwrapper = gophercloud.ErrDefault429{}
	_ unwrapper = gophercloud.ErrDefault500{}
	_ unwrapper = gophercloud.ErrDefault502{}
	_ unwrapper = gophercloud.ErrDefault503{}
	_ unwrapper = gophercloud.ErrDefault504{}
)
