package identity

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
)

var (
	ErrNotFound = fmt.Errorf("Identity extension not found")
)

type ExtensionDetails struct {
	Name        string
	Namespace   string
	Updated     string
	Description string
}

// ExtensionsDesc structures are returned by the Extensions() function for valid input.
// This structure implements the ExtensionInquisitor interface.
type ExtensionsDesc struct {
	extensions []interface{}
}

func Extensions(m map[string]interface{}) *ExtensionsDesc {
	return &ExtensionsDesc{extensions: m["extensions"].([]interface{})}
}

func extensionIndexByAlias(e *ExtensionsDesc, alias string) (int, error) {
	for i, ee := range e.extensions {
		extensionRecord := ee.(map[string]interface{})
		if extensionRecord["alias"] == alias {
			return i, nil
		}
	}
	return 0, ErrNotFound
}

func (e *ExtensionsDesc) IsExtensionAvailable(alias string) bool {
	_, err := extensionIndexByAlias(e, alias)
	return err == nil
}

func (e *ExtensionsDesc) ExtensionDetailsByAlias(alias string) (*ExtensionDetails, error) {
	i, err := extensionIndexByAlias(e, alias)
	if err != nil {
		return nil, err
	}
	ed := &ExtensionDetails{}
	err = mapstructure.Decode(e.extensions[i], ed)
	return ed, err
}
