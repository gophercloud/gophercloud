package stacks

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/rackspace/gophercloud"
	"gopkg.in/yaml.v2"
)

type Client interface {
	Get(string) (*http.Response, error)
}

type TE struct {
	// Bin stores the contents of the template or environment.
	Bin []byte
	// URL stores the URL of the template. This is allowed to be a 'file://'
	// for local files.
	URL string
	// Parsed contains a parsed version of Bin. Since there are 2 different
	// fields referring to the same value, you must be careful when accessing
	// this filed.
	Parsed map[string]interface{}
	// Files contains a mapping between the urls in templates to their contents.
	Files map[string]string
	// fileMaps is a map used internally when determining Files.
	fileMaps map[string]string
	// baseURL represents the location of the template or environment file.
	baseURL string
	// client is an interface which allows TE to fetch contents from URLS
	client Client
}

func (t *TE) Fetch() error {
	// get baseURL if not already defined
	if t.baseURL == "" {
		u, err := getBasePath()
		if err != nil {
			return err
		}
		t.baseURL = u
	}
	if t.Bin != nil {
		// already have contents
		return nil
	}
	u, err := gophercloud.NormalizePathURL(t.baseURL, t.URL)
	if err != nil {
		return err
	}
	t.URL = u
	// get an HTTP client if none present
	if t.client == nil {
		t.client = getHTTPClient()
	}
	resp, err := t.client.Get(t.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	t.Bin = body
	return nil
}

// get the basepath of the template
func getBasePath() (string, error) {
	basePath, err := filepath.Abs(".")
	if err != nil {
		return "", err
	}
	u, err := gophercloud.NormalizePathURL("", basePath)
	if err != nil {
		return "", err
	}
	return u, nil
}

// get a an HTTP client to retrieve URLs
func getHTTPClient() Client {
	transport := &http.Transport{}
	transport.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
	return &http.Client{Transport: transport}
}

// parse the contents and validate
func (t *TE) Parse() error {
	if err := t.Fetch(); err != nil {
		return err
	}
	if jerr := json.Unmarshal(t.Bin, &t.Parsed); jerr != nil {
		if yerr := yaml.Unmarshal(t.Bin, &t.Parsed); yerr != nil {
			return errors.New(fmt.Sprintf("Data in neither json nor yaml format."))
		}
	}
	return t.Validate()
}

// base Validate method, always returns true
func (t *TE) Validate() error {
	return nil
}

type igFunc func(string, interface{}) bool

// convert map[interface{}]interface{} to map[string]interface{}
func toStringKeys(m interface{}) (map[string]interface{}, error) {
	switch m.(type) {
	case map[string]interface{}, map[interface{}]interface{}:
		typed_map := make(map[string]interface{})
		if _, ok := m.(map[interface{}]interface{}); ok {
			for k, v := range m.(map[interface{}]interface{}) {
				typed_map[k.(string)] = v
			}
		} else {
			typed_map = m.(map[string]interface{})
		}
		return typed_map, nil
	default:
		return nil, errors.New(fmt.Sprintf("Expected a map of type map[string]interface{} or map[interface{}]interface{}, actual type: %v", reflect.TypeOf(m)))

	}
}

// fix the template reference to files
func (t *TE) FixFileRefs() {
	t_str := string(t.Bin)
	if t.fileMaps == nil {
		return
	}
	for k, v := range t.fileMaps {
		t_str = strings.Replace(t_str, k, v, -1)
	}
	t.Bin = []byte(t_str)
}
