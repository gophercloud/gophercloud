package identity

import (
	"encoding/json"
	"testing"
)

func TestIsExtensionAvailable(t *testing.T) {
	// Make a response as we'd expect from the IdentityService.GetExtensions() call.
	getExtensionsResults := make(map[string]interface{})
	err := json.Unmarshal([]byte(queryResults), &getExtensionsResults)
	if err != nil {
		t.Error(err)
		return
	}

	e := ExtensionsResult(getExtensionsResults)
	for _, alias := range []string{"OS-KSADM", "OS-FEDERATION"} {
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

	e := ExtensionsResult(getExtensionsResults)
	ed, err := e.ExtensionDetailsByAlias("OS-KSADM")
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
		"name":        "OpenStack Keystone Admin",
		"namespace":   "http://docs.openstack.org/identity/api/ext/OS-KSADM/v1.0",
		"updated":     "2013-07-11T17:14:00-00:00",
		"description": "OpenStack extensions to Keystone v2.0 API enabling Administrative Operations.",
	}

	for k, v := range expecteds {
		if actuals[k] != v {
			t.Errorf("Expected %s \"%s\", got \"%s\" instead", k, v, actuals[k])
			return
		}
	}
}

func TestMalformedResponses(t *testing.T) {
	getExtensionsResults := make(map[string]interface{})
	err := json.Unmarshal([]byte(bogusExtensionsResults), &getExtensionsResults)
	if err != nil {
		t.Error(err)
		return
	}
	e := ExtensionsResult(getExtensionsResults)

	_, err = e.ExtensionDetailsByAlias("OS-KSADM")
	if err == nil {
		t.Error("Expected ErrNotImplemented at least")
		return
	}
	if err != ErrNotImplemented {
		t.Error("Expected ErrNotImplemented")
		return
	}

	if e.IsExtensionAvailable("anything at all") {
		t.Error("No extensions are available with a bogus result.")
		return
	}
}

func TestAliases(t *testing.T) {
	getExtensionsResults := make(map[string]interface{})
	err := json.Unmarshal([]byte(queryResults), &getExtensionsResults)
	if err != nil {
		t.Error(err)
		return
	}

	e := ExtensionsResult(getExtensionsResults)
	aliases, err := e.Aliases()
	if err != nil {
		t.Error(err)
		return
	}
	extensions := (((e["extensions"]).(map[string]interface{}))["values"]).([]interface{})
	if len(aliases) != len(extensions) {
		t.Error("Expected one alias name per extension")
		return
	}
}
