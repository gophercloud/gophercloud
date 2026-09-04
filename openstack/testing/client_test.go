package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
	"github.com/gophercloud/gophercloud/v2/openstack"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

const ID = "0123456789"

func TestAuthenticatedClientV3(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Subject-Token", ID)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{ "token": { "expires_at": "2013-02-02T18:30:59.000000Z" } }`)
	})

	options := auth.AuthOptionsV3{
		AuthURL: fakeServer.Endpoint(),
		Auth: auth.V3PasswordOpts{
			Username:       "me",
			Password:       "secret",
			UserDomainName: "default",
			Scope:          &auth.Scope{ProjectName: "project", ProjectDomainName: "default"},
		},
	}
	client, err := openstack.AuthenticatedClient(context.TODO(), options)
	th.AssertNoErr(t, err)
	th.CheckEquals(t, ID, client.TokenID)
}

func TestAuthenticatedClientV2(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/tokens", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `
			{
				"access": {
					"token": {
						"id": "01234567890",
						"expires": "2014-10-01T10:00:00.000000Z"
					},
					"serviceCatalog": []
				}
			}
		`)
	})

	options := auth.AuthOptionsV2{
		AuthURL: fakeServer.Endpoint(),
		Auth: auth.V2PasswordOpts{
			Username: "me",
			Password: "secret",
		},
	}
	client, err := openstack.AuthenticatedClient(context.TODO(), options)
	th.AssertNoErr(t, err)
	th.CheckEquals(t, "01234567890", client.TokenID)
}

func TestAuthenticatedClientCanReauthWiresReauthFunc(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Subject-Token", ID)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{ "token": { "expires_at": "2013-02-02T18:30:59.000000Z" } }`)
	})

	options := auth.AuthOptionsV3{
		AuthURL: fakeServer.Endpoint(),
		Auth: auth.V3PasswordOpts{
			Username:       "me",
			Password:       "secret",
			UserDomainName: "default",
			AllowReauth:    true,
		},
	}
	client, err := openstack.AuthenticatedClient(context.TODO(), options)
	th.AssertNoErr(t, err)
	if client.ReauthFunc == nil {
		t.Fatal("expected ReauthFunc to be wired for a CanReauth()==true mechanism")
	}
}

func TestAuthenticatedClientCanReauthFalseLeavesReauthFuncNil(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Subject-Token", ID)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{ "token": { "expires_at": "2013-02-02T18:30:59.000000Z" } }`)
	})

	options := auth.AuthOptionsV3{
		AuthURL: fakeServer.Endpoint(),
		Auth:    auth.V3TokenOpts{Token: "sometoken"},
	}
	client, err := openstack.AuthenticatedClient(context.TODO(), options)
	th.AssertNoErr(t, err)
	if client.ReauthFunc != nil {
		t.Fatal("expected ReauthFunc to stay nil for a CanReauth()==false mechanism")
	}
}

func TestAuthenticatedClientReauthenticatesOn401AndRefreshesEndpointLocator(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	tokenCount := 0
	fakeServer.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		tokenCount++
		w.Header().Add("X-Subject-Token", fmt.Sprintf("token-%d", tokenCount))
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{
			"token": {
				"expires_at": "2013-02-02T18:30:59.000000Z",
				"catalog": [{
					"id": "1", "type": "compute", "name": "nova",
					"endpoints": [{"id": "1", "interface": "public", "region": "RegionOne", "url": "%s"}]
				}]
			}
		}`, fmt.Sprintf("http://example.com/compute-%d", tokenCount))
	})

	protectedCalls := 0
	fakeServer.Mux.HandleFunc("/protected", func(w http.ResponseWriter, r *http.Request) {
		protectedCalls++
		if protectedCalls == 1 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	options := auth.AuthOptionsV3{
		AuthURL: fakeServer.Endpoint(),
		Auth: auth.V3PasswordOpts{
			Username:       "me",
			Password:       "secret",
			UserDomainName: "default",
			AllowReauth:    true,
		},
	}
	client, err := openstack.AuthenticatedClient(context.TODO(), options)
	th.AssertNoErr(t, err)
	th.CheckEquals(t, "token-1", client.TokenID)

	endpointBefore, err := client.EndpointLocator(context.TODO(), gophercloud.EndpointOpts{Type: "compute"})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, "http://example.com/compute-1/", endpointBefore)

	resp, err := client.Request(context.TODO(), "GET", fakeServer.Endpoint()+"protected", &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	th.AssertNoErr(t, err)
	defer resp.Body.Close()

	th.CheckEquals(t, "token-2", client.TokenID)
	endpointAfter, err := client.EndpointLocator(context.TODO(), gophercloud.EndpointOpts{Type: "compute"})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, "http://example.com/compute-2/", endpointAfter)
}

func testAuthenticatedClientFails(t *testing.T, endpoint string) {
	options := auth.AuthOptionsV3{
		AuthURL: endpoint,
		Auth: auth.V3PasswordOpts{
			Username:       "me",
			Password:       "secret",
			UserDomainName: "default",
		},
	}
	_, err := openstack.AuthenticatedClient(context.TODO(), options)
	if err == nil {
		t.Fatal("expected error but call succeeded")
	}
}

func TestAuthenticatedClientFails(t *testing.T) {
	testAuthenticatedClientFails(t, "http://bad-address.example.com/v3")
}
