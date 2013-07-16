package gophercloud

import (
	"github.com/racker/perigee"
)

// WithReauth wraps a Perigee request fragment with logic to perform re-authentication
// if it's deemed necessary.
//
// Do not confuse this function with WithReauth()!  Although they work together to support reauthentication,
// WithReauth() actually contains the decision-making logic to determine when to perform a reauth,
// while WithReauthHandler() is used to configure what a reauth actually entails.
func (c *Context) WithReauth(f func() (interface{}, error)) (interface{}, error) {
	result, err := f()
	cause, ok := err.(*perigee.UnexpectedResponseCodeError)
	if ok && cause.Actual == 401 {
		err = c.reauthHandler()
		if err == nil {
			result, err = f()
		}
	}
	return result, err
}