package v2

import (
	"github.com/mitchellh/mapstructure"
)

// ExtensionDetails provides the details offered by the OpenStack Identity V2 extensions API
// for a named extension.
//
// Name provides the name, presumably the same as that used to query the API with.
//
// Updated provides, in a sense, the version of the extension supported.  It gives the timestamp
// of the most recent extension deployment.
//
// Description provides a more customer-oriented description of the extension.
type ExtensionDetails struct {
	Name        string
	Namespace   string
	Updated     string
	Description string
}

// ExtensionsResult encapsulates the raw data returned by a call to
// GetExtensions().  As OpenStack extensions may freely alter the response
// bodies of structures returned to the client, you may only safely access the
// data provided through separate, type-safe accessors or methods.
type ExtensionsResult map[string]interface{}

// IsExtensionAvailable returns true if and only if the provider supports the named extension.
func (er ExtensionsResult) IsExtensionAvailable(alias string) bool {
	e, err := extensions(er)
	if err != nil {
		return false
	}
	_, err = extensionIndexByAlias(e, alias)
	return err == nil
}

// ExtensionDetailsByAlias returns more detail than the mere presence of an extension by the provider.
// See the ExtensionDetails structure.
func (er ExtensionsResult) ExtensionDetailsByAlias(alias string) (*ExtensionDetails, error) {
	e, err := extensions(er)
	if err != nil {
		return nil, err
	}
	i, err := extensionIndexByAlias(e, alias)
	if err != nil {
		return nil, err
	}
	ed := &ExtensionDetails{}
	err = mapstructure.Decode(e[i], ed)
	return ed, err
}

func extensionIndexByAlias(records []interface{}, alias string) (int, error) {
	for i, er := range records {
		extensionRecord := er.(map[string]interface{})
		if extensionRecord["alias"] == alias {
			return i, nil
		}
	}
	return 0, ErrNotImplemented
}

func extensions(er ExtensionsResult) ([]interface{}, error) {
	ei, ok := er["extensions"]
	if !ok {
		return nil, ErrNotImplemented
	}
	e := ei.(map[string]interface{})
	vi, ok := e["values"]
	if !ok {
		return nil, ErrNotImplemented
	}
	v := vi.([]interface{})
	return v, nil
}

// Aliases returns the set of extension handles, or "aliases" as OpenStack calls them.
// These are not the names of the extensions, but rather opaque, symbolic monikers for their corresponding extension.
// Use the ExtensionDetailsByAlias() method to query more information for an extension if desired.
func (er ExtensionsResult) Aliases() ([]string, error) {
	e, err := extensions(er)
	if err != nil {
		return nil, err
	}
	aliases := make([]string, len(e))
	for i, ex := range e {
		ext := ex.(map[string]interface{})
		extn, ok := ext["alias"]
		if ok {
			aliases[i] = extn.(string)
		}
	}
	return aliases, nil
}
