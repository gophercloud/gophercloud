package certificates

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
)

// CreateOptsBuilder allows extensions to add additional parameters
// to the Create request.
type CreateOptsBuilder interface {
	ToCertificateCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create a certificate.
type CreateOpts struct {
	BayUUID string `json:"bay_uuid" required:"true"`
	CSR     string `json:"csr" required:"true"`
}

// ToCertificateCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToCertificateCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Get makes a request against the API to get details for a certificate.
func Get(client *gophercloud.ServiceClient, clusterID string) (r GetResult) {
	url := getURL(client, clusterID)

	_, r.Err = client.Get(url, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})

	return
}

// Create requests the creation of a new certificate.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToCertificateCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	var result *http.Response
	result, r.Err = client.Post(createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})

	if r.Err == nil {
		r.Header = result.Header
	}

	return
}
