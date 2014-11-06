package throttle

import (
	"errors"

	"github.com/racker/perigee"

	"github.com/rackspace/gophercloud"
)

// UpdateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Create operation in this package.
type UpdateOptsBuilder interface {
	ToSSLUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts is the common options struct used in this package's Update
// operation.
type UpdateOpts struct {
	// Required
	SecurePort int

	// Required
	PrivateKey string

	// Required
	Certificate string

	// Required
	IntCertificate string

	// Optional
	Enabled *bool

	// Optional
	SecureTrafficOnly *bool
}

// ToSSLUpdateMap casts a CreateOpts struct to a map.
func (opts UpdateOpts) ToSSLUpdateMap() (map[string]interface{}, error) {
	ssl := make(map[string]interface{})

	if opts.SecurePort == 0 {
		return ssl, errors.New("SecurePort needs to be an integer greater than 0")
	}
	if opts.PrivateKey == "" {
		return ssl, errors.New("PrivateKey is a required field")
	}
	if opts.Certificate == "" {
		return ssl, errors.New("Certificate is a required field")
	}
	if opts.IntCertificate == "" {
		return ssl, errors.New("IntCertificate is a required field")
	}

	ssl["securePort"] = opts.SecurePort
	ssl["privateKey"] = opts.PrivateKey
	ssl["certificate"] = opts.Certificate
	ssl["intermediatecertificate"] = opts.IntCertificate

	if opts.Enabled != nil {
		ssl["enabled"] = &opts.Enabled
	}

	if opts.SecureTrafficOnly != nil {
		ssl["secureTrafficOnly"] = &opts.SecureTrafficOnly
	}

	return map[string]interface{}{"sslTermination": ssl}, nil
}

// Update is the operation responsible for updating the SSL Termination
// configuration for a load balancer.
func Update(c *gophercloud.ServiceClient, lbID int, opts UpdateOptsBuilder) UpdateResult {
	var res UpdateResult

	reqBody, err := opts.ToSSLUpdateMap()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = perigee.Request("PUT", rootURL(c, lbID), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res.Body,
		OkCodes:     []int{200},
	})

	return res
}

// Get is the operation responsible for showing the details of the SSL
// Termination configuration for a load balancer.
func Get(c *gophercloud.ServiceClient, lbID int) GetResult {
	var res GetResult

	_, res.Err = perigee.Request("GET", rootURL(c, lbID), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		Results:     &res.Body,
		OkCodes:     []int{200},
	})

	return res
}

// Delete is the operation responsible for deleting the SSL Termination
// configuration for a load balancer.
func Delete(c *gophercloud.ServiceClient, lbID int) DeleteResult {
	var res DeleteResult

	_, res.Err = perigee.Request("DELETE", rootURL(c, lbID), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		OkCodes:     []int{200},
	})

	return res
}
