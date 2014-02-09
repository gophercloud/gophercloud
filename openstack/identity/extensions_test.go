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

	e := ExtensionsResult(getExtensionsResults)
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

func TestMalformedResponses(t *testing.T) {
	getExtensionsResults := make(map[string]interface{})
	err := json.Unmarshal([]byte(bogusExtensionsResults), &getExtensionsResults)
	if err != nil {
		t.Error(err)
		return
	}
	e := ExtensionsResult(getExtensionsResults)

	_, err = e.ExtensionDetailsByAlias("RS-META")
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
	if len(aliases) != len(e) {
		t.Error("Expected one alias name per extension")
		return
	}
}
