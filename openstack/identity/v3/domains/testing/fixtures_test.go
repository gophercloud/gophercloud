package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/domains"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// ListAvailableOutput provides a single page of available domain results.
const ListAvailableOutput = `
{
    "domains": [
        {
            "id": "52af04aec5f84182b06959d2775d2000",
            "name": "TestDomain",
            "description": "Testing domain",
            "enabled": false,
            "links": {
                "self": "https://example.com/v3/domains/52af04aec5f84182b06959d2775d2000"
            }
        },
        {
            "id": "a720688fb87f4575a4c000d818061eae",
            "name": "ProdDomain",
            "description": "Production domain",
            "enabled": true,
            "links": {
                "self": "https://example.com/v3/domains/a720688fb87f4575a4c000d818061eae"
            }
        }
    ],
    "links": {
        "next": null,
        "self": "https://example.com/v3/auth/domains",
        "previous": null
    }
}
`

// ListOutput provides a single page of Domain results.
const ListOutput = `
{
    "links": {
        "next": null,
        "previous": null,
        "self": "http://example.com/identity/v3/domains"
    },
    "domains": [
        {
            "enabled": true,
            "id": "2844b2a08be147a08ef58317d6471f1f",
            "links": {
                "self": "http://example.com/identity/v3/domains/2844b2a08be147a08ef58317d6471f1f"
            },
            "name": "domain one",
            "description": "some description"
        },
        {
            "enabled": true,
            "id": "9fe1d3",
            "links": {
                "self": "https://example.com/identity/v3/domains/9fe1d3"
            },
            "name": "domain two"
        }
    ]
}
`

// GetOutput provides a Get result.
const GetOutput = `
{
    "domain": {
        "enabled": true,
        "id": "9fe1d3",
        "links": {
            "self": "https://example.com/identity/v3/domains/9fe1d3"
        },
        "name": "domain two"
    }
}
`

// CreateRequest provides the input to a Create request.
const CreateRequest = `
{
    "domain": {
        "name": "domain two"
    }
}
`

// UpdateRequest provides the input to as Update request.
const UpdateRequest = `
{
    "domain": {
        "description": "Staging Domain"
    }
}
`

// UpdateOutput provides an update result.
const UpdateOutput = `
{
    "domain": {
		"enabled": true,
        "id": "9fe1d3",
        "links": {
            "self": "https://example.com/identity/v3/domains/9fe1d3"
        },
        "name": "domain two",
        "description": "Staging Domain"
    }
}
`

// ProdDomain is a domain fixture.
var ProdDomain = domains.Domain{
	Enabled: true,
	ID:      "a720688fb87f4575a4c000d818061eae",
	Links: map[string]any{
		"self": "https://example.com/v3/domains/a720688fb87f4575a4c000d818061eae",
	},
	Name:        "ProdDomain",
	Description: "Production domain",
}

// TestDomain is a domain fixture.
var TestDomain = domains.Domain{
	Enabled: false,
	ID:      "52af04aec5f84182b06959d2775d2000",
	Links: map[string]any{
		"self": "https://example.com/v3/domains/52af04aec5f84182b06959d2775d2000",
	},
	Name:        "TestDomain",
	Description: "Testing domain",
}

// FirstDomain is the first domain in the List request.
var FirstDomain = domains.Domain{
	Enabled: true,
	ID:      "2844b2a08be147a08ef58317d6471f1f",
	Links: map[string]any{
		"self": "http://example.com/identity/v3/domains/2844b2a08be147a08ef58317d6471f1f",
	},
	Name:        "domain one",
	Description: "some description",
}

// SecondDomain is the second domain in the List request.
var SecondDomain = domains.Domain{
	Enabled: true,
	ID:      "9fe1d3",
	Links: map[string]any{
		"self": "https://example.com/identity/v3/domains/9fe1d3",
	},
	Name: "domain two",
}

// SecondDomainUpdated is how SecondDomain should look after an Update.
var SecondDomainUpdated = domains.Domain{
	Enabled: true,
	ID:      "9fe1d3",
	Links: map[string]any{
		"self": "https://example.com/identity/v3/domains/9fe1d3",
	},
	Name:        "domain two",
	Description: "Staging Domain",
}

// ExpectedAvailableDomainsSlice is the slice of domains expected to be returned
// from ListAvailableOutput.
var ExpectedAvailableDomainsSlice = []domains.Domain{TestDomain, ProdDomain}

// ExpectedDomainsSlice is the slice of domains expected to be returned from ListOutput.
var ExpectedDomainsSlice = []domains.Domain{FirstDomain, SecondDomain}

// HandleListAvailableDomainsSuccessfully creates an HTTP handler at `/auth/domains`
// on the test handler mux that responds with a list of two domains.
func HandleListAvailableDomainsSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/auth/domains", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListAvailableOutput)
	})
}

// HandleListDomainsSuccessfully creates an HTTP handler at `/domains` on the
// test handler mux that responds with a list of two domains.
func HandleListDomainsSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/domains", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListOutput)
	})
}

// HandleGetDomainSuccessfully creates an HTTP handler at `/domains` on the
// test handler mux that responds with a single domain.
func HandleGetDomainSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/domains/9fe1d3", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, GetOutput)
	})
}

// HandleCreateDomainSuccessfully creates an HTTP handler at `/domains` on the
// test handler mux that tests domain creation.
func HandleCreateDomainSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/domains", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, CreateRequest)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, GetOutput)
	})
}

// HandleDeleteDomainSuccessfully creates an HTTP handler at `/domains` on the
// test handler mux that tests domain deletion.
func HandleDeleteDomainSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/domains/9fe1d3", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})
}

// HandleUpdateDomainSuccessfully creates an HTTP handler at `/domains` on the
// test handler mux that tests domain update.
func HandleUpdateDomainSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/domains/9fe1d3", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, UpdateRequest)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, UpdateOutput)
	})
}
