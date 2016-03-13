package tokens

import "github.com/gophercloud/gophercloud"

// AuthOptionsBuilder describes any argument that may be passed to the Create call.
type AuthOptionsBuilder interface {
	// ToTokenCreateMap assembles the Create request body, returning an error if parameters are
	// missing or inconsistent.
	ToTokenV2CreateMap() (map[string]interface{}, error)
}

// Create authenticates to the identity service and attempts to acquire a Token.
// If successful, the CreateResult
// Generally, rather than interact with this call directly, end users should call openstack.AuthenticatedClient(),
// which abstracts all of the gory details about navigating service catalogs and such.
func Create(client *gophercloud.ServiceClient, auth AuthOptionsBuilder) CreateResult {
	var r CreateResult
	b, err := auth.ToTokenV2CreateMap()
	if err != nil {
		r.Err = err
		return r
	}
	_, r.Err = client.Post(CreateURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 203},
	})
	return r
}

// Get validates and retrieves information for user's token.
func Get(client *gophercloud.ServiceClient, token string) GetResult {
	var r GetResult
	_, r.Err = client.Get(GetURL(client, token), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 203},
	})
	return r
}
