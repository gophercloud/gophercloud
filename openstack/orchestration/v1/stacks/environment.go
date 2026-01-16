package stacks

import (
	"fmt"
	"strings"

	yaml "go.yaml.in/yaml/v3"
)

// Environment is a structure that represents stack environments
type Environment struct {
	TE
}

// EnvironmentSections is a map containing allowed sections in a stack environment file
var EnvironmentSections = map[string]bool{
	"parameters":         true,
	"parameter_defaults": true,
	"resource_registry":  true,
}

// Validate validates the contents of the Environment
func (e *Environment) Validate() error {
	if e.Parsed == nil {
		if err := e.Parse(); err != nil {
			return err
		}
	}
	for key := range e.Parsed {
		if _, ok := EnvironmentSections[key]; !ok {
			return ErrInvalidEnvironment{Section: key}
		}
	}
	return nil
}

// Parse environment file to resolve the URL's of the resources. This is done by
// reading from the `Resource Registry` section, which is why the function is
// named GetRRFileContents.
func (e *Environment) getRRFileContents(ignoreIf igFunc) error {
	// initialize environment if empty
	if e.Files == nil {
		e.Files = make(map[string]string)
	}
	if e.fileMaps == nil {
		e.fileMaps = make(map[string]string)
	}

	// get the resource registry, only process further if it is a map
	rr, ok := e.Parsed["resource_registry"].(map[string]any)
	if !ok {
		return nil
	}

	// search the resource registry for URLs
	//
	// the resource registry might contain a base URL for the resource. If
	// such a field is present, use it. Otherwise, use the default base URL.
	var baseURL string
	if val, ok := rr["base_url"]; ok {
		baseURL = val.(string)
	} else {
		baseURL = e.baseURL
	}

	// The contents of the resource may be located in a remote file, which
	// will be a template. Instantiate a temporary template to manage the
	// contents.
	tempTemplate := new(Template)
	tempTemplate.baseURL = baseURL
	tempTemplate.client = e.client

	// Fetch the contents of remote resource URL's
	if err := tempTemplate.getFileContents(rr, ignoreIf, false); err != nil {
		return err
	}
	// check the `resources` section (if it exists) for more URL's. Note that
	// the previous call to GetFileContents was (deliberately) not recursive
	// as we want more control over where to look for URL's
	if resourcesMap, ok := rr["resources"].(map[string]any); ok {
		for _, v := range resourcesMap {
			if resourceMap, ok := v.(map[string]any); ok {
				var resourceBaseURL string
				// if base_url for the resource type is defined, use it
				if val, ok := resourceMap["base_url"]; ok {
					resourceBaseURL = val.(string)
				} else {
					resourceBaseURL = baseURL
				}
				tempTemplate.baseURL = resourceBaseURL
				if err := tempTemplate.getFileContents(v, ignoreIf, false); err != nil {
					return err
				}
			}
		}
	}
	// if the resource registry contained any URL's, store them. This can
	// then be passed as parameter to api calls to Heat api.
	e.Files = tempTemplate.Files

	// In case some element was updated, regenerate the string representation
	if len(e.Files) > 0 {
		var err error
		e.Bin, err = yaml.Marshal(&e.Parsed)
		if err != nil {
			return fmt.Errorf("failed to marshal updated environment: %w", err)
		}
	}

	return nil
}

// function to choose keys whose values are other environment files
func ignoreIfEnvironment(key string, value any) bool {
	// base_url and hooks refer to components which cannot have urls
	if key == "base_url" || key == "hooks" {
		return true
	}
	// if value is not string, it cannot be a URL
	valueString, ok := value.(string)
	if !ok {
		return true
	}
	// if value contains `::`, it must be a reference to another resource type
	// e.g. OS::Nova::Server : Rackspace::Cloud::Server
	if strings.Contains(valueString, "::") {
		return true
	}
	return false
}
