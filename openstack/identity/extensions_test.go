package identity

import (
	"encoding/json"
	"testing"
)

// Taken from: http://docs.openstack.org/api/openstack-identity-service/2.0/content/GET_listExtensions_v2.0_extensions_.html#GET_listExtensions_v2.0_extensions_-Request
const queryResults = `{
    "extensions":[{
            "name": "Reset Password Extension",
            "namespace": "http://docs.rackspacecloud.com/identity/api/ext/rpe/v2.0",
            "alias": "RS-RPE",
            "updated": "2011-01-22T13:25:27-06:00",
            "description": "Adds the capability to reset a user's password. The user is emailed when the password has been reset.",
            "links":[{
                    "rel": "describedby",
                    "type": "application/pdf",
                    "href": "http://docs.rackspacecloud.com/identity/api/ext/identity-rpe-20111111.pdf"
                },
                {
                    "rel": "describedby",
                    "type": "application/vnd.sun.wadl+xml",
                    "href": "http://docs.rackspacecloud.com/identity/api/ext/identity-rpe.wadl"
                }
            ]
        },
        {
            "name": "User Metadata Extension",
            "namespace": "http://docs.rackspacecloud.com/identity/api/ext/meta/v2.0",
            "alias": "RS-META",
            "updated": "2011-01-12T11:22:33-06:00",
            "description": "Allows associating arbritrary metadata with a user.",
            "links":[{
                    "rel": "describedby",
                    "type": "application/pdf",
                    "href": "http://docs.rackspacecloud.com/identity/api/ext/identity-meta-20111201.pdf"
                },
                {
                    "rel": "describedby",
                    "type": "application/vnd.sun.wadl+xml",
                    "href": "http://docs.rackspacecloud.com/identity/api/ext/identity-meta.wadl"
                }
            ]
        }
    ],
    "extensions_links":[]
}`

func TestIsExtensionAvailable(t *testing.T) {
	// Make a response as we'd expect from the IdentityService.GetExtensions() call.
	getExtensionsResults := make(map[string]interface{})
	err := json.Unmarshal([]byte(queryResults), &getExtensionsResults)
	if err != nil {
		t.Error(err)
		return
	}

	e := Extensions(getExtensionsResults)
	for _, alias := range []string{"RS-RPE", "RS-META"} {
		if !e.IsExtensionAvailable(alias) {
			t.Errorf("Expected extension %s present.", alias)
			return
		}
	}
	if e.IsExtensionAvailable("blort") {
		t.Errorf("Input JSON doesn't list blort as an extension")
		return
	}
}

func TestGetExtensionDetails(t *testing.T) {
	// Make a response as we'd expect from the IdentityService.GetExtensions() call.
	getExtensionsResults := make(map[string]interface{})
	err := json.Unmarshal([]byte(queryResults), &getExtensionsResults)
	if err != nil {
		t.Error(err)
		return
	}

	e := Extensions(getExtensionsResults)
	ed, err := e.ExtensionDetailsByAlias("RS-META")
	if err != nil {
		t.Error(err)
		return
	}

	actuals := map[string]string{
		"name":        ed.Name,
		"namespace":   ed.Namespace,
		"updated":     ed.Updated,
		"description": ed.Description,
	}

	expecteds := map[string]string{
		"name":        "User Metadata Extension",
		"namespace":   "http://docs.rackspacecloud.com/identity/api/ext/meta/v2.0",
		"updated":     "2011-01-12T11:22:33-06:00",
		"description": "Allows associating arbritrary metadata with a user.",
	}

	for k, v := range expecteds {
		if actuals[k] != v {
			t.Errorf("Expected %s \"%s\", got \"%s\" instead", k, v, actuals[k])
			return
		}
	}
}
