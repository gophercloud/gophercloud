package certificates

import (
	"encoding/pem"

	"github.com/gophercloud/gophercloud"
)

// CertificateResult temporarily contains the response from a GenerateCertificate or ImportCertificate call.
type CertificateResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a CertificateResult and extracts a certificate resource.
func (r CertificateResult) Extract() (*BayCertificate, error) {
	var s struct {
		BayID       string `json:"bay_uuid"`
		Certificate string `json:"pem"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return nil, err
	}

	pemBlock, _ := pem.Decode([]byte(s.Certificate))
	certificate := &BayCertificate{
		BayID:       s.BayID,
		Certificate: *pemBlock,
	}
	return certificate, nil
}

// Certificate represents a certificate associated with a bay
type BayCertificate struct {
	BayID       string
	Certificate pem.Block
}

// String returns a PEM encoded string representation of the certificate
func (c BayCertificate) String() string {
	return string(pem.EncodeToMemory(&c.Certificate))
}

// CreateCredentialsBundleResult temporarily contains the response from a CreateCredentialsBundle call.
type CreateCredentialsBundleResult struct {
	gophercloud.Result
}

// CredentialsBundle is a collection of certificates and supporting files necessary to communicate with a bay.
type CredentialsBundle struct {
	BayID         string
	COEEndpoint   string
	Certificate   pem.Block
	PrivateKey    pem.Block
	CACertificate pem.Block
	Scripts       map[string][]byte
}
