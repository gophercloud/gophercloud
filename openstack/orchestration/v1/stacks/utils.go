package stacks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"reflect"

	"github.com/gophercloud/gophercloud/v2"
	yaml "gopkg.in/yaml.v2"
)

// Client is an interface that expects a Get method similar to http.Get. This
// is needed for unit testing, since we can mock an http client. Thus, the
// client will usually be an http.Client EXCEPT in unit tests.
type Client interface {
	Get(string) (*http.Response, error)
}

// TE is a base structure for both Template and Environment
type TE struct {
	// Bin stores the contents of the template or environment.
	Bin []byte
	// URL stores the URL of the template. This is allowed to be a 'file://'
	// for local files.
	URL string
	// Parsed contains a parsed version of Bin. Since there are 2 different
	// fields referring to the same value, you must be careful when accessing
	// this filed.
	Parsed map[string]any
	// Files contains a mapping between the urls in templates to their contents.
	Files map[string]string
	// fileMaps is a map used internally when determining Files.
	fileMaps map[string]string
	// baseURL represents the location of the template or environment file.
	baseURL string
	// client is an interface which allows TE to fetch contents from URLS
	client Client
}

// Fetch fetches the contents of a TE from its URL. Once a TE structure has a
// URL, call the fetch method to fetch the contents.
func (t *TE) Fetch() error {
	// if the baseURL is not provided, use the current directors as the base URL
	if t.baseURL == "" {
		u, err := getBasePath()
		if err != nil {
			return err
		}
		t.baseURL = u
	}

	// if the contents are already present, do nothing.
	if t.Bin != nil {
		return nil
	}

	// get a fqdn from the URL using the baseURL of the TE. For local files,
	// the URL's will have the `file` scheme.
	u, err := gophercloud.NormalizePathURL(t.baseURL, t.URL)
	if err != nil {
		return err
	}
	t.URL = u

	// get an HTTP client if none present
	if t.client == nil {
		t.client = getHTTPClient()
	}

	// use the client to fetch the contents of the TE
	resp, err := t.client.Get(t.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return fmt.Errorf("error fetching %s: %s", t.URL, resp.Status)
	}
	t.Bin = body
	return nil
}

// get the basepath of the TE
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

// get a an HTTP client to retrieve URL's. This client allows the use of `file`
// scheme since we may need to fetch files from users filesystem
func getHTTPClient() Client {
	transport := &http.Transport{}
	transport.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
	return &http.Client{Transport: transport}
}

// Parse will parse the contents and then validate. The contents MUST be either JSON or YAML.
func (t *TE) Parse() error {
	if err := t.Fetch(); err != nil {
		return err
	}
	if jerr := json.Unmarshal(t.Bin, &t.Parsed); jerr != nil {
		if yerr := yaml.Unmarshal(t.Bin, &t.Parsed); yerr != nil {
			return ErrInvalidDataFormat{}
		}
	}
	return nil
}

// igfunc is a parameter used by GetFileContents and GetRRFileContents to check
// for valid URL's.
type igFunc func(string, any) bool

// convert map[any]any to map[string]any
func toStringKeys(m any) (map[string]any, error) {
	switch m.(type) {
	case map[string]any, map[any]any:
		typedMap := make(map[string]any)
		if _, ok := m.(map[any]any); ok {
			for k, v := range m.(map[any]any) {
				typedMap[k.(string)] = v
			}
		} else {
			typedMap = m.(map[string]any)
		}
		return typedMap, nil
	default:
		return nil, gophercloud.ErrUnexpectedType{Expected: "map[string]any/map[any]any", Actual: fmt.Sprintf("%v", reflect.TypeOf(m))}
	}
}
