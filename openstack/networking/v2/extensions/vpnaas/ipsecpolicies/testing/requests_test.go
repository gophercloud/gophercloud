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
		Units: "seconds",
		Value: 7200,
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
	expectedLifetime := ipsecpolicies.Lifetime{
		Units: "seconds",
		Value: 7200,
	}
	expected := ipsecpolicies.Policy{
		TenantID:            "b4eedccc6fb74fa8a7ad6b08382b852b",
		Name:                "ipsecpolicy1",
		TransformProtocol:   "esp",
		AuthAlgorithm:       "sha1",
		EncapsulationMode:   "tunnel",
		EncryptionAlgorithm: "aes-128",
		PFS:                 "group5",
		Description:         "",
		Lifetime:            expectedLifetime,
	}
	th.AssertDeepEquals(t, expected, *actual)
}
