package stacks

import (
	"errors"
	"fmt"
	"strings"
)

// an interface to represent stack environments
type Environment struct {
	TE
}

// allowed sections in a stack environment file
var EnvironmentSections = map[string]bool{
	"parameters":         true,
	"parameter_defaults": true,
	"resource_registry":  true,
}

func (e *Environment) Validate() error {
	if e.Parsed == nil {
		if err := e.Parse(); err != nil {
			return err
		}
	}
	for key, _ := range e.Parsed {
		if _, ok := EnvironmentSections[key]; !ok {
			return errors.New(fmt.Sprintf("Environment has wrong section: %s", key))
		}
	}
	return nil
}

// Parse environment file to resolve the urls of the resources
func GetRRFileContents(e *Environment, ignoreIf igFunc) error {
	if e.Files == nil {
		e.Files = make(map[string]string)
	}
	if e.fileMaps == nil {
		e.fileMaps = make(map[string]string)
	}
	rr := e.Parsed["resource_registry"]
	// search the resource registry for URLs
	switch rr.(type) {
	case map[string]interface{}, map[interface{}]interface{}:
		rr_map, err := toStringKeys(rr)
		if err != nil {
			return err
		}
		var baseURL string
		if val, ok := rr_map["base_url"]; ok {
			baseURL = val.(string)
		} else {
			baseURL = e.baseURL
		}
		// use a fake template to fetch contents from URLs
		tempTemplate := new(Template)
		tempTemplate.baseURL = baseURL
		tempTemplate.client = e.client

		if err = GetFileContents(tempTemplate, rr, ignoreIf, false); err != nil {
			return err
		}
		// check the `resources` section (if it exists) for more URLs
		if val, ok := rr_map["resources"]; ok {
			switch val.(type) {
			case map[string]interface{}, map[interface{}]interface{}:
				resources_map, err := toStringKeys(val)
				if err != nil {
					return err
				}
				for _, v := range resources_map {
					switch v.(type) {
					case map[string]interface{}, map[interface{}]interface{}:
						resource_map, err := toStringKeys(v)
						if err != nil {
							return err
						}
						var resourceBaseURL string
						// if base_url for the resource type is defined, use it
						if val, ok := resource_map["base_url"]; ok {
							resourceBaseURL = val.(string)
						} else {
							resourceBaseURL = baseURL
						}
						tempTemplate.baseURL = resourceBaseURL
						if err := GetFileContents(tempTemplate, v, ignoreIf, false); err != nil {
							return err
						}
					}

				}

			}
		}
		e.Files = tempTemplate.Files
		return nil
	default:
		return nil
	}
}

// function to choose keys whose values are other environment files
func ignoreIfEnvironment(key string, value interface{}) bool {
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
