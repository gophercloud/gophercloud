package certificates

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
)

// CreateOptsBuilder allows extensions to add additional parameters
// to the Create request.
type CreateOptsBuilder interface {
	ToCertificateCreateMap() (map[string]any, error)
}

// CreateOpts represents options used to create a certificate.
type CreateOpts struct {
	ClusterUUID string `json:"cluster_uuid,omitempty" xor:"BayUUID"`
	BayUUID     string `json:"bay_uuid,omitempty" xor:"ClusterUUID"`
	CSR         string `json:"csr" required:"true"`
}

// ToCertificateCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToCertificateCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Get makes a request against the API to get details for a certificate.
func Get(ctx context.Context, client *gophercloud.ServiceClient, clusterID string) (r GetResult) {
	url := getURL(client, clusterID)
	resp, err := client.Get(ctx, url, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Create requests the creation of a new certificate.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToCertificateCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Update will rotate the CA certificate for a cluster
func Update(ctx context.Context, client *gophercloud.ServiceClient, clusterID string) (r UpdateResult) {
	resp, err := client.Patch(ctx, updateURL(client, clusterID), nil, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
