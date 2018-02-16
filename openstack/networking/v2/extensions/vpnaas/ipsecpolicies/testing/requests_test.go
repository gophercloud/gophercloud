package testing

import (
	"fmt"
	"net/http"
	"testing"

	fake "github.com/gophercloud/gophercloud/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/ipsecpolicies"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/vpn/ipsecpolicies", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "ipsecpolicy": {
        "name": "ipsecpolicy1",
        "transform_protocol": "esp",
        "auth_algorithm": "sha1",
        "encapsulation_mode": "tunnel",
        "encryption_algorithm": "aes-128",
        "pfs": "group5",
        "lifetime": {
            "units": "seconds",
            "value": 7200
        },
        "tenant_id": "b4eedccc6fb74fa8a7ad6b08382b852b"
}
}      `)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, `
{
    "ipsecpolicy": {
        "name": "ipsecpolicy1",
        "transform_protocol": "esp",
        "auth_algorithm": "sha1",
        "encapsulation_mode": "tunnel",
        "encryption_algorithm": "aes-128",
        "pfs": "group5",
        "tenant_id": "b4eedccc6fb74fa8a7ad6b08382b852b",
        "lifetime": {
            "units": "seconds",
            "value": 7200
        },
        "id": "5291b189-fd84-46e5-84bd-78f40c05d69c",
        "description": ""
    }
}
    `)
	})

	lifetime := ipsecpolicies.LifetimeCreateOpts{
		LifetimeUnits: "seconds",
		LifetimeValue: 7200,
	}
	options := ipsecpolicies.CreateOpts{
		TenantID:            "b4eedccc6fb74fa8a7ad6b08382b852b",
		Name:                "ipsecpolicy1",
		TransformProtocol:   "esp",
		AuthAlgorithm:       "sha1",
		EncapsulationMode:   "tunnel",
		EncryptionAlgorithm: "aes-128",
		PFS:                 "group5",
		Lifetime:            &lifetime,
		Description:         "",
	}
	actual, err := ipsecpolicies.Create(fake.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "b4eedccc6fb74fa8a7ad6b08382b852b", actual.TenantID)
	th.AssertEquals(t, "ipsecpolicy1", actual.Name)
	th.AssertEquals(t, "esp", actual.TransformProtocol)
	th.AssertEquals(t, "sha1", actual.AuthAlgorithm)
	th.AssertEquals(t, "tunnel", actual.EncapsulationMode)
	th.AssertEquals(t, "aes-128", actual.EncryptionAlgorithm)
	th.AssertEquals(t, "group5", actual.PFS)
	th.AssertEquals(t, "", actual.Description)
	th.AssertEquals(t, "seconds", actual.Lifetime.LifetimeUnits)
	th.AssertEquals(t, 7200, actual.Lifetime.LifetimeValue)
}
