package tokens

import (
	"context"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
)

// Create authenticates to the identity service and attempts to acquire a Token.
// Generally, rather than interact with this call directly, end users should
// call openstack.AuthenticatedClient(), which abstracts all of the gory details
// about navigating service catalogs and such.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts auth.AuthOptionsBuilderV2) (r CreateResult) {
	request, err := auth.NewRequestV2(opts)
	if err != nil {
		r.Err = err
		return
	}

	request.JSONResponse = &r.Body
	resp, err := client.Request(ctx, http.MethodPost, CreateURL(client), request)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get validates and retrieves information for user's token.
func Get(ctx context.Context, client *gophercloud.ServiceClient, token string) (r GetResult) {
	resp, err := client.Get(ctx, GetURL(client, token), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 203},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
