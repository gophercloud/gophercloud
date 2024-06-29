package httpbasic

import (
	"encoding/base64"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
)

// EndpointOpts specifies a "http_basic" Ironic Inspector Endpoint.
type EndpointOpts struct {
	IronicInspectorEndpoint     string
	IronicInspectorUser         string
	IronicInspectorUserPassword string
}

func initClientOpts(client *gophercloud.ProviderClient, eo EndpointOpts) (gophercloud.Client, error) {
	sc := new(gophercloud.Client)
	if eo.IronicInspectorEndpoint == "" {
		return nil, fmt.Errorf("IronicInspectorEndpoint is required")
	}
	if eo.IronicInspectorUser == "" || eo.IronicInspectorUserPassword == "" {
		return nil, fmt.Errorf("User and Password are required")
	}

	token := []byte(eo.IronicInspectorUser + ":" + eo.IronicInspectorUserPassword)
	encodedToken := base64.StdEncoding.EncodeToString(token)
	sc.MoreHeaders = map[string]string{"Authorization": "Basic " + encodedToken}
	sc.Endpoint = gophercloud.NormalizeURL(eo.IronicInspectorEndpoint)
	sc.ProviderClient = client
	return sc, nil
}

// NewBareMetalIntrospectionHTTPBasic creates a ServiceClient that may be used to access a
// "http_basic" bare metal introspection service.
func NewBareMetalIntrospectionHTTPBasic(eo EndpointOpts) (gophercloud.Client, error) {
	sc, err := initClientOpts(&gophercloud.ProviderClient{}, eo)
	if err != nil {
		return nil, err
	}

	sc.Type = "baremetal-introspection"

	return sc, nil
}
