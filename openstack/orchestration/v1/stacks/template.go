package stacks

import (
	"fmt"
	"net/url"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
	yaml "gopkg.in/yaml.v2"
)

// Template is a structure that represents OpenStack Heat templates
type Template struct {
	TE
}

// TemplateFormatVersions is a map containing allowed variations of the template format version
// Note that this contains the permitted variations of the _keys_ not the values.
var TemplateFormatVersions = map[string]bool{
	"HeatTemplateFormatVersion": true,
	"heat_template_version":     true,
	"AWSTemplateFormatVersion":  true,
}

// Validate validates the contents of the Template
func (t *Template) Validate() error {
	if t.Parsed == nil {
		if err := t.Parse(); err != nil {
			return err
		}
	}
	var invalid string
	for key := range t.Parsed {
		if _, ok := TemplateFormatVersions[key]; ok {
			return nil
		}
		invalid = key
	}
	return ErrInvalidTemplateFormatVersion{Version: invalid}
}

func (t *Template) makeChildTemplate(childURL string, ignoreIf igFunc, recurse bool) (*Template, error) {
	// create a new child template
	childTemplate := new(Template)

	// initialize child template

	// get the base location of the child template. Child path is relative
	// to its parent location so that templates can be composed
	if t.URL != "" {
		// Preserve all elements of the URL but take the directory part of the path
		u, err := url.Parse(t.URL)
		if err != nil {
			return nil, err
		}
		u.Path = filepath.Dir(u.Path)
		childTemplate.baseURL = u.String()
	}
	childTemplate.URL = childURL
	childTemplate.client = t.client

	// fetch the contents of the child template or file
	if err := childTemplate.Fetch(); err != nil {
		return nil, err
	}

	// process child template recursively if required. This is
	// required if the child template itself contains references to
	// other templates
	if recurse {
		if err := childTemplate.Parse(); err == nil {
			if err := childTemplate.Validate(); err == nil {
				if err := childTemplate.getFileContents(childTemplate.Parsed, ignoreIf, recurse); err != nil {
					return nil, err
				}
			}
		}
	}

	return childTemplate, nil
}

// Applies the transformation for getFileContents() to just one element of a map.
// In case the element requires transforming, the function returns its new value.
func (t *Template) mapElemFileContents(k any, v any, ignoreIf igFunc, recurse bool) (any, error) {
	key, ok := k.(string)
	if !ok {
		return nil, fmt.Errorf("can't convert map key to string: %v", k)
	}

	value, ok := v.(string)
	if !ok {
		// if the value is not a string, recursively parse that value
		if err := t.getFileContents(v, ignoreIf, recurse); err != nil {
			return nil, err
		}
	} else if !ignoreIf(key, value) {
		// at this point, the k, v pair has a reference to an external template
		// or file (for 'get_file' function).
		// The assumption of heatclient is that value v is a reference
		// to a file in the users environment, so we have to the path

		// create a new child template with the referenced contents
		childTemplate, err := t.makeChildTemplate(value, ignoreIf, recurse)
		if err != nil {
			return nil, err
		}

		// update parent template with current child templates' content.
		// At this point, the child template has been parsed recursively.
		t.fileMaps[value] = childTemplate.URL
		t.Files[childTemplate.URL] = string(childTemplate.Bin)

		// Also add child templates' own children (templates or get_file)!
		for k, v := range childTemplate.Files {
			t.Files[k] = v
		}

		return childTemplate.URL, nil
	}

	return nil, nil
}

// GetFileContents recursively parses a template to search for urls. These urls
// are assumed to point to other templates (known in OpenStack Heat as child
// templates). The contents of these urls are fetched and stored in the `Files`
// parameter of the template structure. This is the only way that a user can
// use child templates that are located in their filesystem; urls located on the
// web (e.g. on github or swift) can be fetched directly by Heat engine.
func (t *Template) getFileContents(te any, ignoreIf igFunc, recurse bool) error {
	// initialize template if empty
	if t.Files == nil {
		t.Files = make(map[string]string)
	}
	if t.fileMaps == nil {
		t.fileMaps = make(map[string]string)
	}

	updated := false

	switch teTyped := (te).(type) {
	// if te is a map[string], go check all elements for URLs to replace
	case map[string]any:
		for k, v := range teTyped {
			newVal, err := t.mapElemFileContents(k, v, ignoreIf, recurse)
			if err != nil {
				return err
			} else if newVal != nil {
				teTyped[k] = newVal
				updated = true
			}
		}
	// same if te is a map[non-string] (can't group with above case because we
	// can't range over and update 'te' without knowing its key type)
	case map[any]any:
		for k, v := range teTyped {
			newVal, err := t.mapElemFileContents(k, v, ignoreIf, recurse)
			if err != nil {
				return err
			} else if newVal != nil {
				teTyped[k] = newVal
				updated = true
			}
		}
	// if te is a slice, call the function on each element of the slice.
	case []any:
		for i := range teTyped {
			if err := t.getFileContents(teTyped[i], ignoreIf, recurse); err != nil {
				return err
			}
		}
	// if te is anything else, there is nothing to do.
	case string, bool, float64, nil, int:
		return nil
	default:
		return gophercloud.ErrUnexpectedType{Actual: fmt.Sprintf("%v", reflect.TypeOf(te))}
	}

	// In case some element was updated, we have to regenerate the string representation
	if updated {
		var err error
		t.Bin, err = yaml.Marshal(&t.Parsed)
		if err != nil {
			return fmt.Errorf("failed to marshal updated data: %w", err)
		}
	}
	return nil
}

// function to choose keys whose values are other template files
func ignoreIfTemplate(key string, value any) bool {
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
