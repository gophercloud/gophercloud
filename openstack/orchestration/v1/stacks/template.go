package stacks

import (
	"errors"
	"fmt"
	"github.com/rackspace/gophercloud"
	"reflect"
	"strings"
)

type Template struct {
	TE
}

var TemplateFormatVersions = map[string]bool{
	"HeatTemplateFormatVersion": true,
	"heat_template_version":     true,
	"AWSTemplateFormatVersion":  true,
}

func (t *Template) Validate() error {
	if t.Parsed == nil {
		if err := t.Parse(); err != nil {
			return err
		}
	}
	for key, _ := range t.Parsed {
		if _, ok := TemplateFormatVersions[key]; ok {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Template format version not found."))
}

func GetFileContents(t *Template, te interface{}, ignoreIf igFunc, recurse bool) error {
	if t.Files == nil {
		t.Files = make(map[string]string)
	}
	if t.fileMaps == nil {
		t.fileMaps = make(map[string]string)
	}
	switch te.(type) {
	case map[string]interface{}, map[interface{}]interface{}:
		te_map, err := toStringKeys(te)
		if err != nil {
			return err
		}
		for k, v := range te_map {
			value, ok := v.(string)
			if !ok {
				if err := GetFileContents(t, v, ignoreIf, recurse); err != nil {
					return err
				}
			} else if !ignoreIf(k, value) {
				// at this point, the k, v pair has a reference to an external template.
				// The assumption of heatclient is that value v is a relative reference
				// to a file in the users environment
				childTemplate := new(Template)
				baseURL, err := gophercloud.NormalizePathURL(t.baseURL, value)
				if err != nil {
					return err
				}
				childTemplate.baseURL = baseURL
				childTemplate.client = t.client
				if err := childTemplate.Parse(); err != nil {
					return err
				}
				// process child template recursively if required
				if recurse {
					if err := GetFileContents(childTemplate, childTemplate.Parsed, ignoreIf, recurse); err != nil {
						return err
					}
				}
				// update parent template with current child templates' content
				t.fileMaps[value] = childTemplate.URL
				t.Files[childTemplate.URL] = string(childTemplate.Bin)

			}
		}
		return nil
	case []interface{}:
		te_slice := te.([]interface{})
		for i := range te_slice {
			if err := GetFileContents(t, te_slice[i], ignoreIf, recurse); err != nil {
				return err
			}
		}
	case string, bool, float64, nil, int:
		return nil
	default:
		return errors.New(fmt.Sprintf("%v: Unrecognized type", reflect.TypeOf(te)))

	}
	return nil
}

// function to choose keys whose values are other template files
func ignoreIfTemplate(key string, value interface{}) bool {
	// key must be either `get_file` or `type` for value to be a URL
	if key != "get_file" && key != "type" {
		return true
	}
	// value must be a string
	valueString, ok := value.(string)
	if !ok {
		return true
	}
	// `.template` and `.yaml` are allowed suffixes for template URLs when referred to by `type`
	if key == "type" && !(strings.HasSuffix(valueString, ".template") || strings.HasSuffix(valueString, ".yaml")) {
		return true
	}
	return false
}
