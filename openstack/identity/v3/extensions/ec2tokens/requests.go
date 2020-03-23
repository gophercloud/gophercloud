package ec2tokens

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
)

const (
	// EC2CredentialsAwsRequestV4 is a constant, used to generate AWS
	// Credential V4.
	EC2CredentialsAwsRequestV4 = "aws4_request"
	// EC2CredentialsHmacSha1V2 is a HMAC SHA1 signature method. Used to
	// generate AWS Credential V2.
	EC2CredentialsHmacSha1V2 = "HmacSHA1"
	// EC2CredentialsHmacSha256V2 is a HMAC SHA256 signature method. Used
	// to generate AWS Credential V2.
	EC2CredentialsHmacSha256V2 = "HmacSHA256"
	// EC2CredentialsAwsHmacV4 is an AWS signature V4 signing method.
	// More details:
	// https://docs.aws.amazon.com/general/latest/gr/signature-version-4.html
	EC2CredentialsAwsHmacV4 = "AWS4-HMAC-SHA256"
	// EC2CredentialsTimestampFormatV4 is an AWS signature V4 timestamp
	// format.
	EC2CredentialsTimestampFormatV4 = "20060102T150405Z"
	// EC2CredentialsDateFormatV4 is an AWS signature V4 date format.
	EC2CredentialsDateFormatV4 = "20060102"
)

// AuthOptions represents options for authenticating a user using EC2 credentials.
type AuthOptions struct {
	// Access is the EC2 Credential Access ID.
	Access string `json:"access" required:"true"`
	// Secret is the EC2 Credential Secret, used to calculate signature.
	// Not used, when a Signature is is.
	Secret string `json:"-"`
	// Host is a HTTP request Host header. Used to calculate an AWS
	// signature V2. For signature V4 set the Host inside Headers map.
	// Optional.
	Host string `json:"host"`
	// Path is a HTTP request path. Optional.
	Path string `json:"path"`
	// Verb is a HTTP request method. Optional.
	Verb string `json:"verb"`
	// Headers is a map of HTTP request headers. Optional.
	Headers map[string]string `json:"headers"`
	// Region is a region name to calculate an AWS signature V4. Optional.
	Region string `json:"-"`
	// Service is a service name to calculate an AWS signature V4. Optional.
	Service string `json:"-"`
	// Params is a map of GET method parameters. Optional.
	Params map[string]string `json:"params"`
	// AllowReauth allows Gophercloud to re-authenticate automatically
	// if/when your token expires.
	AllowReauth bool `json:"-"`
	// Signature can be either a []byte (encoded to base64 automatically) or
	// a string. You can set the singature explicitly, when you already know
	// it. In this case default Params won't be automatically set. Optional.
	Signature interface{} `json:"signature"`
	// BodyHash is a HTTP request body sha256 hash. When nil and Signature
	// is not set, a random hash is generated. Optional.
	BodyHash *string `json:"body_hash"`
	// Timestamp is a timestamp to calculate a V4 signature. Optional.
	Timestamp *time.Time `json:"-"`
}

// EC2CredentialsBuildCanonicalQueryStringV2 builds a canonical query string
// for an AWS signature V2.
// https://github.com/openstack/python-keystoneclient/blob/stable/train/keystoneclient/contrib/ec2/utils.py#L133
func EC2CredentialsBuildCanonicalQueryStringV2(params map[string]string) string {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var pairs []string
	for _, k := range keys {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, url.QueryEscape(params[k])))
	}

	return strings.Join(pairs, "&")
}

// EC2CredentialsBuildStringToSignV2 builds a string to sign an AWS signature
// V2.
// https://github.com/openstack/python-keystoneclient/blob/stable/train/keystoneclient/contrib/ec2/utils.py#L148
func EC2CredentialsBuildStringToSignV2(opts AuthOptions) string {
	stringToSign := strings.Join([]string{
		opts.Verb,
		opts.Host,
		opts.Path,
	}, "\n")

	return strings.Join([]string{
		stringToSign,
		EC2CredentialsBuildCanonicalQueryStringV2(opts.Params),
	}, "\n")
}

// EC2CredentialsBuildCanonicalQueryStringV2 builds a canonical query string
// for an AWS signature V4.
// https://github.com/openstack/python-keystoneclient/blob/stable/train/keystoneclient/contrib/ec2/utils.py#L244
func EC2CredentialsBuildCanonicalQueryStringV4(verb string, params map[string]string) string {
	if verb == "POST" {
		return ""
	}
	return EC2CredentialsBuildCanonicalQueryStringV2(params)
}

// EC2CredentialsBuildCanonicalHeadersV4 builds a canonical string based on
// "headers" map and "signedHeaders" string parameters.
// https://github.com/openstack/python-keystoneclient/blob/stable/train/keystoneclient/contrib/ec2/utils.py#L216
func EC2CredentialsBuildCanonicalHeadersV4(headers map[string]string, signedHeaders string) string {
	headersLower := make(map[string]string, len(headers))
	for k, v := range headers {
		headersLower[strings.ToLower(k)] = v
	}

	var headersList []string
	for _, h := range strings.Split(signedHeaders, ";") {
		if v, ok := headersLower[h]; ok {
			headersList = append(headersList, h+":"+v)
		}
	}

	return strings.Join(headersList, "\n") + "\n"
}

// EC2CredentialsBuildSignatureKeyV4 builds a HMAC 256 signature key based on
// input parameters.
// https://github.com/openstack/python-keystoneclient/blob/stable/train/keystoneclient/contrib/ec2/utils.py#L169
func EC2CredentialsBuildSignatureKeyV4(secret, date, region, service string) []byte {
	kDate := sumHMAC256([]byte("AWS4"+secret), []byte(date))
	kRegion := sumHMAC256(kDate, []byte(region))
	kService := sumHMAC256(kRegion, []byte(service))
	return sumHMAC256(kService, []byte(EC2CredentialsAwsRequestV4))
}

// EC2CredentialsBuildSignatureV4 builds an AWS v4 signature based on input
// parameters.
// https://github.com/openstack/python-keystoneclient/blob/stable/train/keystoneclient/contrib/ec2/utils.py#L251
func EC2CredentialsBuildSignatureV4(opts AuthOptions, signedHeaders string, date time.Time, bodyHash string) string {
	scope := strings.Join([]string{
		date.Format(EC2CredentialsDateFormatV4),
		opts.Region,
		opts.Service,
		EC2CredentialsAwsRequestV4,
	}, "/")

	canonicalRequest := strings.Join([]string{
		opts.Verb,
		opts.Path,
		EC2CredentialsBuildCanonicalQueryStringV4(opts.Verb, opts.Params),
		EC2CredentialsBuildCanonicalHeadersV4(opts.Headers, signedHeaders),
		signedHeaders,
		bodyHash,
	}, "\n")
	hash := sha256.Sum256([]byte(canonicalRequest))

	strToSign := strings.Join([]string{
		EC2CredentialsAwsHmacV4,
		date.Format(EC2CredentialsTimestampFormatV4),
		scope,
		hex.EncodeToString(hash[:]),
	}, "\n")

	key := EC2CredentialsBuildSignatureKeyV4(opts.Secret, date.Format(EC2CredentialsDateFormatV4), opts.Region, opts.Service)

	return hex.EncodeToString(sumHMAC256(key, []byte(strToSign)))
}

// ToTokenV3ScopeMap is a dummy method to satisfy tokens.AuthOptionsBuilder interface
func (opts *AuthOptions) ToTokenV3ScopeMap() (map[string]interface{}, error) {
	return nil, nil
}

// CanReauth is a method method to satisfy tokens.AuthOptionsBuilder interface
func (opts *AuthOptions) CanReauth() bool {
	return opts.AllowReauth
}

// ToTokenV3CreateMap formats an AuthOptions into a create request.
func (opts *AuthOptions) ToTokenV3CreateMap(map[string]interface{}) (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "credentials")
	if err != nil {
		return nil, err
	}

	if opts.Signature != nil {
		return b, nil
	}

	// calculate signature, when it is not set
	c, _ := b["credentials"].(map[string]interface{})
	h := interfaceToMap(c, "headers")
	p := interfaceToMap(c, "params")

	// detect and process a signature v2
	if v, ok := p["SignatureVersion"]; ok && v == "2" {
		if _, ok := c["body_hash"]; ok {
			delete(c, "body_hash")
		}
		if _, ok := c["headers"]; ok {
			delete(c, "headers")
		}
		if v, ok := p["SignatureMethod"]; ok {
			// params is a map of strings
			strToSign := EC2CredentialsBuildStringToSignV2(*opts)
			switch v {
			case EC2CredentialsHmacSha1V2:
				// keystone uses this method only when HmacSHA256 is not available on the server side
				// https://github.com/openstack/python-keystoneclient/blob/stable/train/keystoneclient/contrib/ec2/utils.py#L151..L156
				c["signature"] = sumHMAC1([]byte(opts.Secret), []byte(strToSign))
				return b, nil
			case EC2CredentialsHmacSha256V2:
				c["signature"] = sumHMAC256([]byte(opts.Secret), []byte(strToSign))
				return b, nil
			}
			return nil, fmt.Errorf("unsupported signature method: %s", v)
		}
		return nil, fmt.Errorf("signature method must be provided")
	} else if ok {
		return nil, fmt.Errorf("unsupported signature version: %s", v)
	}

	// it is not a signature v2, but a signature v4
	date := time.Now().UTC()
	if opts.Timestamp != nil {
		date = *opts.Timestamp
	}
	if v, _ := c["body_hash"]; v == nil {
		// when body_hash is not set, generate a random one
		c["body_hash"] = randomBodyHash()
	}

	signedHeaders, _ := h["X-Amz-SignedHeaders"]
	c["signature"] = EC2CredentialsBuildSignatureV4(*opts, signedHeaders, date, c["body_hash"].(string))

	h["X-Amz-Date"] = date.Format(EC2CredentialsTimestampFormatV4)
	h["Authorization"] = fmt.Sprintf("%s Credential=%s/%s/%s/%s/%s, SignedHeaders=%s, Signature=%s",
		EC2CredentialsAwsHmacV4,
		opts.Access,
		date.Format(EC2CredentialsDateFormatV4),
		opts.Region,
		opts.Service,
		EC2CredentialsAwsRequestV4,
		signedHeaders,
		c["signature"])

	return b, nil
}

// Create authenticates and either generates a new token from EC2 credentials
func Create(c *gophercloud.ServiceClient, opts tokens.AuthOptionsBuilder) (r tokens.CreateResult) {
	b, err := opts.ToTokenV3CreateMap(nil)
	if err != nil {
		r.Err = err
		return
	}

	resp, err := c.Post(ec2tokensURL(c), b, &r.Body, &gophercloud.RequestOpts{
		MoreHeaders: map[string]string{"X-Auth-Token": ""},
		OkCodes:     []int{200},
	})
	r.Err = err
	if resp != nil {
		r.Header = resp.Header
	}
	return
}

// The following are small helper functions used to help build the signature.

// sumHMAC1 is a func to implement the HMAC SHA1 signature method.
func sumHMAC1(key []byte, data []byte) []byte {
	hash := hmac.New(sha1.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}

// sumHMAC256 is a func to implement the HMAC SHA256 signature method.
func sumHMAC256(key []byte, data []byte) []byte {
	hash := hmac.New(sha256.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}

// randomBodyHash is a func to generate a random sha256 hexdigest.
func randomBodyHash() string {
	h := make([]byte, 64)
	rand.Read(h)
	return hex.EncodeToString(h)
}

// interfaceToMap is a func used to represent a "credentials" map element as a
// "map[string]string"
func interfaceToMap(c map[string]interface{}, key string) map[string]string {
	// convert map[string]interface{} to map[string]string
	m := make(map[string]string)
	if v, _ := c[key].(map[string]interface{}); v != nil {
		for k, v := range v {
			m[k] = v.(string)
		}
	}

	c[key] = m

	return m
}
