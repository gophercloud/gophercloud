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

// GetFileContents recursively parses a template to search for urls. These urls
// are assumed to point to other templates (known in OpenStack Heat as child
// templates). The contents of these urls are fetched and stored in the `Files`
// parameter of the template structure. This is the only way that a user can
// use child templates that are located in their filesystem; urls located on the
// web (e.g. on github or swift) can be fetched directly by Heat engine.
func (t *Template) GetFileContents(te interface{}, ignoreIf igFunc, recurse bool) error {
	// initialize template if empty
	if t.Files == nil {
		t.Files = make(map[string]string)
	}
	if t.fileMaps == nil {
		t.fileMaps = make(map[string]string)
	}
	switch te.(type) {
	// if te is a map
	case map[string]interface{}, map[interface{}]interface{}:
		te_map, err := toStringKeys(te)
		if err != nil {
			return err
		}
		for k, v := range te_map {
			value, ok := v.(string)
			if !ok {
				// if the value is not a string, recursively parse that value
				if err := t.GetFileContents(v, ignoreIf, recurse); err != nil {
					return err
				}
			} else if !ignoreIf(k, value) {
				// at this point, the k, v pair has a reference to an external template.
				// The assumption of heatclient is that value v is a reference
				// to a file in the users environment

				// create a new child template
				childTemplate := new(Template)

				// initialize child template

				// get the base location of the child template
				baseURL, err := gophercloud.NormalizePathURL(t.baseURL, value)
				if err != nil {
					return err
				}
				childTemplate.baseURL = baseURL
				childTemplate.client = t.client

				// fetch the contents of the child template
				if err := childTemplate.Parse(); err != nil {
					return err
				}

				// process child template recursively if required. This is
				// required if the child template itself contains references to
				// other templates
				if recurse {
					if err := childTemplate.GetFileContents(childTemplate.Parsed, ignoreIf, recurse); err != nil {
						return err
					}
				}
				// update parent template with current child templates' content.
				// At this point, the child template has been parsed recursively.
				t.fileMaps[value] = childTemplate.URL
				t.Files[childTemplate.URL] = string(childTemplate.Bin)

			}
		}
		return nil
	// if te is a slice, call the function on each element of the slice.
	case []interface{}:
		te_slice := te.([]interface{})
		for i := range te_slice {
			if err := t.GetFileContents(te_slice[i], ignoreIf, recurse); err != nil {
				return err
			}
		}
	// if te is anything else, return
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
