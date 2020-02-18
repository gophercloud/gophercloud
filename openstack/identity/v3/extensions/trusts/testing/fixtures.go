package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/extensions/trusts"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	"github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

const CreateRequest = `
{
    "trust": {
        "expires_at": "2019-12-01T14:00:00.999999Z",
        "impersonation": false,
        "allow_redelegation": true,
        "project_id": "9b71012f5a4a4aef9193f1995fe159b2",
        "roles": [
            {
                "name": "member"
            }
        ],
        "trustee_user_id": "ecb37e88cc86431c99d0332208cb6fbf",
        "trustor_user_id": "959ed913a32c4ec88c041c98e61cbbc3"
    }
}
`

const CreateResponse = `
{
    "trust": {
        "expires_at": "2019-12-01T14:00:00.999999Z",
        "id": "3422b7c113894f5d90665e1a79655e23",
        "impersonation": false,
        "redelegation_count": 10,
        "project_id": "9b71012f5a4a4aef9193f1995fe159b2",
        "remaining_uses": null,
        "roles": [
            {
                "id": "b627fca5-beb0-471a-9857-0e852b719e76",
                "links": {
                    "self": "http://example.com/identity/v3/roles/b627fca5-beb0-471a-9857-0e852b719e76"
                },
                "name": "member"
            }
        ],
        "trustee_user_id": "ecb37e88cc86431c99d0332208cb6fbf",
        "trustor_user_id": "959ed913a32c4ec88c041c98e61cbbc3"
    }
}
`

// GetOutput provides a Get result.
const GetResponse = `
{
    "trust": {
        "id": "987fe8",
        "expires_at": "2013-02-27T18:30:59.999999Z",
        "impersonation": true,
        "links": {
            "self": "http://example.com/identity/v3/OS-TRUST/trusts/987fe8"
        },
        "roles": [
            {
                "id": "ed7b78",
                "links": {
                    "self": "http://example.com/identity/v3/roles/ed7b78"
                },
                "name": "member"
            }
        ],
        "roles_links": {
            "next": null,
            "previous": null,
            "self": "http://example.com/identity/v3/OS-TRUST/trusts/1ff900/roles"
        },
        "project_id": "0f1233",
        "trustee_user_id": "be34d1",
        "trustor_user_id": "56ae32"
    }
}
`

// ListOutput provides a single page of Role results.
const ListResponse = `
{
    "trusts": [
        {
            "id": "1ff900",
            "expires_at": "2013-02-27T18:30:59.999999Z",
            "impersonation": true,
            "links": {
                "self": "http://example.com/identity/v3/OS-TRUST/trusts/1ff900"
            },
            "project_id": "0f1233",
            "trustee_user_id": "86c0d5",
            "trustor_user_id": "a0fdfd"
        },
        {
            "id": "f4513a",
            "impersonation": false,
            "links": {
                "self": "http://example.com/identity/v3/OS-TRUST/trusts/f45513a"
            },
            "project_id": "0f1233",
            "trustee_user_id": "86c0d5",
            "trustor_user_id": "3cd2ce"
        }
    ]
}
`

// HandleCreateTokenWithTrustID verifies that providing certain AuthOptions and Scope results in an expected JSON structure.
func HandleCreateTokenWithTrustID(t *testing.T, options tokens.AuthOptionsBuilder, requestJSON string) {
	testhelper.Mux.HandleFunc("/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "POST")
		testhelper.TestHeader(t, r, "Content-Type", "application/json")
		testhelper.TestHeader(t, r, "Accept", "application/json")
		testhelper.TestJSONRequest(t, r, requestJSON)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{
    "token": {
        "expires_at": "2013-02-27T18:30:59.999999Z",
        "issued_at": "2013-02-27T16:30:59.999999Z",
        "methods": [
            "password"
        ],
        "OS-TRUST:trust": {
            "id": "fe0aef",
            "impersonation": false,
						"redelegated_trust_id": "3ba234",
						"redelegation_count": 2,
            "links": {
                "self": "http://example.com/identity/v3/trusts/fe0aef"
            },
            "trustee_user": {
                "id": "0ca8f6",
                "links": {
                    "self": "http://example.com/identity/v3/users/0ca8f6"
                }
            },
            "trustor_user": {
                "id": "bd263c",
                "links": {
                    "self": "http://example.com/identity/v3/users/bd263c"
                }
            }
        },
        "user": {
            "domain": {
                "id": "1789d1",
                "links": {
                    "self": "http://example.com/identity/v3/domains/1789d1"
                },
                "name": "example.com"
            },
            "email": "joe@example.com",
            "id": "0ca8f6",
            "links": {
                "self": "http://example.com/identity/v3/users/0ca8f6"
            },
            "name": "Joe"
        }
    }
}`)
	})
}

// HandleCreateTrust creates an HTTP handler at `/OS-TRUST/trusts` on the
// test handler mux that tests trust creation.
func HandleCreateTrust(t *testing.T) {
	testhelper.Mux.HandleFunc("/OS-TRUST/trusts", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "POST")
		testhelper.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		testhelper.TestJSONRequest(t, r, CreateRequest)

		w.WriteHeader(http.StatusCreated)
		_, err := fmt.Fprintf(w, CreateResponse)
		testhelper.AssertNoErr(t, err)
	})
}

// HandleDeleteUserSuccessfully creates an HTTP handler at `/users` on the
// test handler mux that tests user deletion.
func HandleDeleteTrust(t *testing.T) {
	testhelper.Mux.HandleFunc("/OS-TRUST/trusts/3422b7c113894f5d90665e1a79655e23", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "DELETE")
		testhelper.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})
}

// HandleGetTrustSuccessfully creates an HTTP handler at `/OS-TRUST/trusts` on the
// test handler mux that responds with a single trusts.
func HandleGetTrustSuccessfully(t *testing.T) {
	testhelper.Mux.HandleFunc("/OS-TRUST/trusts/987fe8", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "GET")
		testhelper.TestHeader(t, r, "Accept", "application/json")
		testhelper.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, GetResponse)
	})
}

var FirstTrust = trusts.Trust{
	ID:            "1ff900",
	Impersonation: true,
	TrusteeUserID: "86c0d5",
	TrustorUserID: "a0fdfd",
	ProjectID:     "0f1233",
	ExpiresAt:     time.Date(2013, 02, 27, 18, 30, 59, 999999000, time.UTC),
	DeletedAt:     time.Time{},
}

var SecondTrust = trusts.Trust{
	ID:            "f4513a",
	Impersonation: false,
	TrusteeUserID: "86c0d5",
	TrustorUserID: "3cd2ce",
	ProjectID:     "0f1233",
	ExpiresAt:     time.Time{},
	DeletedAt:     time.Time{},
}

// ExpectedRolesSlice is the slice of roles expected to be returned from ListOutput.
var ExpectedTrustsSlice = []trusts.Trust{FirstTrust, SecondTrust}

// HandleListTrustsSuccessfully creates an HTTP handler at `/OS-TRUST/trusts` on the
// test handler mux that responds with a list of two trusts.
func HandleListTrustsSuccessfully(t *testing.T) {
	testhelper.Mux.HandleFunc("/OS-TRUST/trusts", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "GET")
		testhelper.TestHeader(t, r, "Accept", "application/json")
		testhelper.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListResponse)
	})
}
