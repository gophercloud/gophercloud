package ssl

import (
	"errors"

	"github.com/racker/perigee"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
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

var (
	errPrivateKey     = errors.New("PrivateKey is a required field")
	errCertificate    = errors.New("Certificate is a required field")
	errIntCertificate = errors.New("IntCertificate is a required field")
)

// ToSSLUpdateMap casts a CreateOpts struct to a map.
func (opts UpdateOpts) ToSSLUpdateMap() (map[string]interface{}, error) {
	ssl := make(map[string]interface{})

	if opts.SecurePort == 0 {
		return ssl, errors.New("SecurePort needs to be an integer greater than 0")
	}
	if opts.PrivateKey == "" {
		return ssl, errPrivateKey
	}
	if opts.Certificate == "" {
		return ssl, errCertificate
	}
	if opts.IntCertificate == "" {
		return ssl, errIntCertificate
	}

	ssl["securePort"] = opts.SecurePort
	ssl["privateKey"] = opts.PrivateKey
	ssl["certificate"] = opts.Certificate
	ssl["intermediateCertificate"] = opts.IntCertificate

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

func ListCerts(c *gophercloud.ServiceClient, lbID int) pagination.Pager {
	url := certURL(c, lbID)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return CertPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

type AddCertOptsBuilder interface {
	ToCertAddMap() (map[string]interface{}, error)
}

type AddCertOpts struct {
	HostName       string
	PrivateKey     string
	Certificate    string
	IntCertificate string
}

func (opts AddCertOpts) ToCertAddMap() (map[string]interface{}, error) {
	cm := make(map[string]interface{})

	if opts.HostName == "" {
		return cm, errors.New("HostName is a required option")
	}
	if opts.PrivateKey == "" {
		return cm, errPrivateKey
	}
	if opts.Certificate == "" {
		return cm, errCertificate
	}

	cm["hostName"] = opts.HostName
	cm["privateKey"] = opts.PrivateKey
	cm["certificate"] = opts.Certificate

	if opts.IntCertificate != "" {
		cm["intermediateCertificate"] = opts.IntCertificate
	}

	return map[string]interface{}{"certificateMapping": cm}, nil
}

func AddCert(c *gophercloud.ServiceClient, lbID int, opts AddCertOptsBuilder) AddCertResult {
	var res AddCertResult

	reqBody, err := opts.ToCertAddMap()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = perigee.Request("POST", certURL(c, lbID), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res.Body,
		OkCodes:     []int{200},
	})

	return res
}

func GetCert(c *gophercloud.ServiceClient, lbID, certID int) GetCertResult {
	var res GetCertResult

	_, res.Err = perigee.Request("GET", certResourceURL(c, lbID, certID), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		Results:     &res.Body,
		OkCodes:     []int{200},
	})

	return res
}

type UpdateCertOptsBuilder interface {
	ToCertUpdateMap() (map[string]interface{}, error)
}

type UpdateCertOpts struct {
	HostName       string
	PrivateKey     string
	Certificate    string
	IntCertificate string
}

func (opts UpdateCertOpts) ToCertUpdateMap() (map[string]interface{}, error) {
	cm := make(map[string]interface{})

	if opts.HostName != "" {
		cm["hostName"] = opts.HostName
	}
	if opts.PrivateKey != "" {
		cm["privateKey"] = opts.PrivateKey
	}
	if opts.Certificate != "" {
		cm["certificate"] = opts.Certificate
	}
	if opts.IntCertificate != "" {
		cm["intermediateCertificate"] = opts.IntCertificate
	}

	return map[string]interface{}{"certificateMapping": cm}, nil
}

func UpdateCert(c *gophercloud.ServiceClient, lbID, certID int, opts UpdateCertOptsBuilder) UpdateCertResult {
	var res UpdateCertResult

	reqBody, err := opts.ToCertUpdateMap()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = perigee.Request("PUT", certResourceURL(c, lbID, certID), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res.Body,
		OkCodes:     []int{202},
	})

	return res
}

func DeleteCert(c *gophercloud.ServiceClient, lbID, certID int) DeleteResult {
	var res DeleteResult

	_, res.Err = perigee.Request("DELETE", certResourceURL(c, lbID, certID), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		OkCodes:     []int{200},
	})

	return res
}
