package serviceassets

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
)

// DeleteOptsBuilder allows extensions to add additional parameters to the Delete
// request.
type DeleteOptsBuilder interface {
	ToCDNAssetDeleteParams() (string, error)
}

// DeleteOpts is a structure that holds options for deleting CDN service assets.
type DeleteOpts struct {
	// If all is set to true, specifies that the delete occurs against all of the
	// assets for the service.
	All bool `q:"all"`
	// Specifies the relative URL of the asset to be deleted.
	URL string `q:"url"`
}

// ToCDNAssetDeleteParams formats a DeleteOpts into a query string.
func (opts DeleteOpts) ToCDNAssetDeleteParams() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

// Delete accepts a unique ID and deletes the CDN service asset associated with
// it.
func Delete(c *gophercloud.ServiceClient, id string, opts DeleteOptsBuilder) DeleteResult {
	var res DeleteResult
	_, res.Err = perigee.Request("DELETE", deleteURL(c, id), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		OkCodes:     []int{202},
	})
	return res
}
