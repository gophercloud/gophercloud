package tokens

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/identity/v2/tenants"
	os "github.com/rackspace/gophercloud/openstack/identity/v2/tokens"
	th "github.com/rackspace/gophercloud/testhelper"
)

var expectedToken = &os.Token{
	ID:        "aaaabbbbccccdddd",
	ExpiresAt: time.Date(2014, time.January, 31, 15, 30, 58, 0, time.UTC),
	Tenant: tenants.Tenant{
		ID:          "fc394f2ab2df4114bde39905f800dc57",
		Name:        "test",
		Description: "There are many tenants. This one is yours.",
		Enabled:     true,
	},
}

var expectedServiceCatalog = &os.ServiceCatalog{
	Entries: []os.CatalogEntry{
		os.CatalogEntry{
			Name: "inscrutablewalrus",
			Type: "something",
			Endpoints: []os.Endpoint{
				os.Endpoint{
					PublicURL: "http://something0:1234/v2/",
					Region:    "region0",
				},
				os.Endpoint{
					PublicURL: "http://something1:1234/v2/",
					Region:    "region1",
				},
			},
		},
		os.CatalogEntry{
			Name: "arbitrarypenguin",
			Type: "else",
			Endpoints: []os.Endpoint{
				os.Endpoint{
					PublicURL: "http://else0:4321/v3/",
					Region:    "region0",
				},
			},
		},
	},
}

func tokenPost(t *testing.T, options gophercloud.AuthOptions, requestJSON string) os.CreateResult {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	client := gophercloud.ServiceClient{Endpoint: th.Endpoint()}

	th.Mux.HandleFunc("/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, requestJSON)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `
{
  "access": {
    "token": {
      "issued_at": "2014-01-30T15:30:58.000000Z",
      "expires": "2014-01-31T15:30:58Z",
      "id": "aaaabbbbccccdddd",
      "tenant": {
        "description": "There are many tenants. This one is yours.",
        "enabled": true,
        "id": "fc394f2ab2df4114bde39905f800dc57",
        "name": "test"
      }
    },
    "serviceCatalog": [
      {
        "endpoints": [
          {
            "publicURL": "http://something0:1234/v2/",
            "region": "region0"
          },
          {
            "publicURL": "http://something1:1234/v2/",
            "region": "region1"
          }
        ],
        "type": "something",
        "name": "inscrutablewalrus"
      },
      {
        "endpoints": [
          {
            "publicURL": "http://else0:4321/v3/",
            "region": "region0"
          }
        ],
        "type": "else",
        "name": "arbitrarypenguin"
      }
    ]
  }
}
    `)
	})

	return Create(&client, WrapOptions(options))
}

func tokenPostErr(t *testing.T, options gophercloud.AuthOptions, expectedErr error) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	client := gophercloud.ServiceClient{Endpoint: th.Endpoint()}

	th.Mux.HandleFunc("/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{}`)
	})

	actualErr := Create(&client, WrapOptions(options)).Err
	th.CheckEquals(t, expectedErr, actualErr)
}

func isSuccessful(t *testing.T, result os.CreateResult) {
	token, err := result.ExtractToken()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, expectedToken, token)

	serviceCatalog, err := result.ExtractServiceCatalog()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, expectedServiceCatalog, serviceCatalog)
}

func TestCreateTokenWithAPIKey(t *testing.T) {
	options := gophercloud.AuthOptions{
		Username: "me",
		APIKey:   "1234567890abcdef",
	}

	isSuccessful(t, tokenPost(t, options, `
    {
      "auth": {
        "RAX-KSKEY:apiKeyCredentials": {
          "username": "me",
          "apiKey": "1234567890abcdef"
        }
      }
    }
  `))
}
