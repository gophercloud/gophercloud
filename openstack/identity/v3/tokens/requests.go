package tokens

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
)

// AuthOptionsBuilder describes any argument that may be passed to the Create call.
type AuthOptionsBuilder interface {
	// ToTokenV3CreateMap assembles the Create request body, returning an error if parameters are
	// missing or inconsistent.
	ToTokenV3CreateMap(*gophercloud.ScopeOptsV3) (map[string]interface{}, error)
}

func subjectTokenHeaders(c *gophercloud.ServiceClient, subjectToken string) map[string]string {
	return map[string]string{
		"X-Subject-Token": subjectToken,
	}
}

// Create authenticates and either generates a new token, or changes the Scope of an existing token.
func Create(c *gophercloud.ServiceClient, opts AuthOptionsBuilder, scopeOpts *gophercloud.ScopeOptsV3) CreateResult {
	var r CreateResult
	b, err := opts.ToTokenV3CreateMap(scopeOpts)
	if err != nil {
		r.Err = err
		return r
	}
	var resp *http.Response
	resp, r.Err = c.Post(tokenURL(c), b, &r.Body, nil)
	if resp != nil {
		r.Header = resp.Header
	}
	return r
}

// Get validates and retrieves information about another token.
func Get(c *gophercloud.ServiceClient, token string) GetResult {
	var r GetResult
	var resp *http.Response
	resp, r.Err = c.Get(tokenURL(c), &r.Body, &gophercloud.RequestOpts{
		MoreHeaders: subjectTokenHeaders(c, token),
		OkCodes:     []int{200, 203},
	})
	if resp != nil {
		r.Header = resp.Header
	}
	return r
}

// Validate determines if a specified token is valid or not.
func Validate(c *gophercloud.ServiceClient, token string) (bool, error) {
	response, err := c.Request("HEAD", tokenURL(c), &gophercloud.RequestOpts{
		MoreHeaders: subjectTokenHeaders(c, token),
		OkCodes:     []int{204, 404},
	})
	if err != nil {
		return false, err
	}

	return response.StatusCode == 204, nil
}

// Revoke immediately makes specified token invalid.
func Revoke(c *gophercloud.ServiceClient, token string) RevokeResult {
	var r RevokeResult
	_, r.Err = c.Delete(tokenURL(c), &gophercloud.RequestOpts{
		MoreHeaders: subjectTokenHeaders(c, token),
	})
	return r
}
