package gophercloud

import (
	"testing"
	"github.com/racker/perigee"
)

// This reauth-handler does nothing, and returns no error.
func doNothing() error {
	return nil
}

func TestOtherErrorsPropegate(t *testing.T) {
	calls := 0
	c := TestContext().WithReauthHandler(doNothing)

	myObj, err := c.WithReauth(func() (interface{}, error) {
		calls++
		return nil, &perigee.UnexpectedResponseCodeError{
			Expected: []int{204},
			Actual: 404,
		}
	})

	if myObj != nil {
		t.Errorf("Returned nil myObj; got %#v", myObj)
		return
	}
	if err == nil {
		t.Error("Expected MyError to be returned; got nil instead.")
		return
	}
	if _, ok := err.(*perigee.UnexpectedResponseCodeError); !ok {
		t.Error("Expected UnexpectedResponseCodeError; got %#v", err)
		return
	}
	if calls != 1 {
		t.Errorf("Expected the body to be invoked once; found %d calls instead", calls)
		return
	}
}

func Test401ErrorCausesBodyInvokation2ndTime(t *testing.T) {
	calls := 0
	c := TestContext().WithReauthHandler(doNothing)

	myObj, err := c.WithReauth(func() (interface{}, error) {
		calls++
		return nil, &perigee.UnexpectedResponseCodeError{
			Expected: []int{204},
			Actual: 401,
		}
	})

	if myObj != nil {
		t.Errorf("Returned nil myObj; got %#v", myObj)
		return
	}
	if err == nil {
		t.Error("Expected MyError to be returned; got nil instead.")
		return
	}
	if calls != 2 {
		t.Errorf("Expected the body to be invoked once; found %d calls instead", calls)
		return
	}
}

func TestReauthAttemptShouldHappen(t *testing.T) {
	calls := 0
	c := TestContext().WithReauthHandler(func() error {
		calls++
		return nil
	})
	c.WithReauth(func() (interface{}, error) {
		return nil, &perigee.UnexpectedResponseCodeError{
			Expected: []int{204},
			Actual: 401,
		}
	})

	if calls != 1 {
		t.Errorf("Expected Reauthenticator to be called once; found %d instead", calls)
		return
	}
}

type MyError struct {}
func (*MyError) Error() string {
	return "MyError instance"
}

func TestReauthErrorShouldPropegate(t *testing.T) {
	c := TestContext().WithReauthHandler(func() error {
		return &MyError{}
	})

	_, err := c.WithReauth(func() (interface{}, error) {
		return nil, &perigee.UnexpectedResponseCodeError{
			Expected: []int{204},
			Actual: 401,
		}
	})

	if _, ok := err.(*MyError); !ok {
		t.Errorf("Expected a MyError; got %#v", err)
		return
	}
}
