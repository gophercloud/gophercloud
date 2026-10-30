//go:build acceptance || identity || oauth2mtls

package v3

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/oauth2mtls"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestOAuth2MTLSAuthentication(t *testing.T) {
	required := []string{"OS_AUTH_URL", "OS_OAUTH2_ENDPOINT", "OS_OAUTH2_CLIENT_ID", "OS_CERT", "OS_KEY"}
	for _, name := range required {
		if os.Getenv(name) == "" {
			t.Skipf("%s must be set", name)
		}
	}

	cert, err := tls.LoadX509KeyPair(os.Getenv("OS_CERT"), os.Getenv("OS_KEY"))
	th.AssertNoErr(t, err)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}
	if caPath := os.Getenv("OS_CACERT"); caPath != "" {
		ca, err := os.ReadFile(caPath)
		th.AssertNoErr(t, err)
		tlsConfig.RootCAs = x509.NewCertPool()
		if !tlsConfig.RootCAs.AppendCertsFromPEM(ca) {
			t.Fatal("failed to parse OS_CACERT")
		}
	}

	provider, err := openstack.NewClient(os.Getenv("OS_AUTH_URL"))
	th.AssertNoErr(t, err)
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = tlsConfig
	provider.HTTPClient.Transport = transport

	opts := &oauth2mtls.AuthOptions{
		OAuth2Endpoint: os.Getenv("OS_OAUTH2_ENDPOINT"),
		ClientID:       os.Getenv("OS_OAUTH2_CLIENT_ID"),
		AllowReauth:    true,
	}
	ctx := context.Background()
	th.AssertNoErr(t, openstack.AuthenticateV3(ctx, provider, opts, gophercloud.EndpointOpts{}))
	th.AssertNoErr(t, provider.Reauthenticate(ctx, ""))

	compute, err := openstack.NewComputeV2(ctx, provider, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
	th.AssertNoErr(t, err)
	_, err = servers.List(compute, servers.ListOpts{Limit: 1}).AllPages(ctx)
	th.AssertNoErr(t, err)
}
