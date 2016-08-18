package certificates

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/containerorchestration/v1/baymodels"
	"github.com/gophercloud/gophercloud/openstack/containerorchestration/v1/bays"
	"github.com/gophercloud/gophercloud/openstack/containerorchestration/v1/certificates/scripts"
	"github.com/gophercloud/gophercloud/openstack/containerorchestration/v1/common"
)

type createCertificateOpts struct {
	BayID                     string `json:"bay_uuid"`
	CertificateSigningRequest string `json:"csr"`
}

// GetCACertificate retrieves the CA certificate for a bay.
func GetCACertificate(c *gophercloud.ServiceClient, bayID string) (r CertificateResult) {
	ro := &gophercloud.RequestOpts{ErrorContext: &common.ErrorResponse{}}
	_, r.Err = c.Get(getCertificateAuthorityURL(c, bayID), &r.Body, ro)
	return
}

// CreateCertificate associates an existing private key with a bay and returns TLS certificate signed with the bay's CA certificate.
func CreateCertificate(c *gophercloud.ServiceClient, bayID string, privateKey *rsa.PrivateKey) (r CertificateResult) {
	if privateKey == nil || privateKey.D == nil {
		r.Err = errors.New("CreateCertificate requires privateKey to not be nil or empty")
		return
	}
	if privateKey.PublicKey.N == nil {
		r.Err = errors.New("CreateCertificate requires privateKey.PublicKey to not be empty")
		return
	}

	bay, err := bays.Get(c, bayID).Extract()
	if err != nil {
		r.Err = err
		return
	}

	csr, err := generateCertificateSigningRequest(privateKey)
	if err != nil {
		r.Err = err
		return
	}

	b := createCertificateOpts{BayID: bay.ID, CertificateSigningRequest: csr}
	ro := &gophercloud.RequestOpts{ErrorContext: &common.ErrorResponse{}}
	_, r.Err = c.Post(createURL(c), b, &r.Body, ro)
	return
}

func generateCertificateSigningRequest(privateKey *rsa.PrivateKey) (string, error) {
	// For some reason these aren't exported by crypto/x509, so I am duplicating here
	type basicConstraints struct {
		IsCA       bool `asn1:"optional"`
		MaxPathLen int  `asn1:"optional,default:-1"`
	}
	var extKeyUsageClientAuth = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 2}

	// Black magic voodoo necessary to set the key usage for the certificate to be clientAuth
	constraint := basicConstraints{IsCA: true, MaxPathLen: -1}
	constraintBytes, err := asn1.Marshal(constraint)
	if err != nil {
		return "", err
	}

	subject := pkix.Name{CommonName: "gophercloud"}
	rawSubject, err := asn1.Marshal(subject.ToRDNSequence())
	if err != nil {
		return "", err
	}

	// Create a certificate signing request
	template := x509.CertificateRequest{
		RawSubject:         rawSubject,
		Extensions:         []pkix.Extension{pkix.Extension{Id: extKeyUsageClientAuth, Value: constraintBytes}},
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	csr, err := x509.CreateCertificateRequest(rand.Reader, &template, privateKey)
	if err != nil {
		return "", err
	}

	csrPEM := pem.EncodeToMemory(&pem.Block{Bytes: csr, Type: "CERTIFICATE REQUEST"})
	return string(csrPEM), nil
}

func generateCertificate(c *gophercloud.ServiceClient, bayID string) (*rsa.PrivateKey, *BayCertificate, error) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 4096)

	r := CreateCertificate(c, bayID, privateKey)
	if r.Err != nil {
		return nil, nil, r.Err
	}

	cert, err := r.Extract()
	if err != nil {
		return nil, nil, err
	}

	return privateKey, cert, nil
}

// CreateCredentialsBundle generates credentials bundle for the specified bay.
func CreateCredentialsBundle(c *gophercloud.ServiceClient, bayID string) (*CredentialsBundle, error) {
	bay, err := bays.Get(c, bayID).Extract()
	if err != nil {
		return nil, err
	}
	baymodel, err := baymodels.Get(c, bay.BayModelID).Extract()
	if err != nil {
		return nil, err
	}

	key, cert, err := generateCertificate(c, bay.ID)
	if err != nil {
		return nil, err
	}

	caResult := GetCACertificate(c, bay.ID)
	if caResult.Err != nil {
		return nil, err
	}

	ca, err := caResult.Extract()
	if err != nil {
		return nil, err
	}

	bundle := &CredentialsBundle{
		BayID:       bay.ID,
		Certificate: cert.Certificate,
		PrivateKey: pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
		CACertificate: ca.Certificate,
		COEEndpoint:   bay.COEEndpoint,
		Scripts:       scripts.Generate(baymodel.COE, bay),
	}

	return bundle, nil
}
