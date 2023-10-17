package stacks

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestTemplateValidation(t *testing.T) {
	templateJSON := new(Template)
	templateJSON.Bin = []byte(ValidJSONTemplate)
	th.AssertNoErr(t, templateJSON.Validate())

	templateYAML := new(Template)
	templateYAML.Bin = []byte(ValidYAMLTemplate)
	th.AssertNoErr(t, templateYAML.Validate())

	templateInvalid := new(Template)
	templateInvalid.Bin = []byte(InvalidTemplateNoVersion)
	if err := templateInvalid.Validate(); err == nil {
		t.Error("Template validation did not catch invalid template")
	}
}

func TestTemplateParsing(t *testing.T) {
	templateJSON := new(Template)
	templateJSON.Bin = []byte(ValidJSONTemplate)
	th.AssertNoErr(t, templateJSON.Parse())
	th.AssertDeepEquals(t, ValidJSONTemplateParsed, templateJSON.Parsed)

	templateYAML := new(Template)
	templateYAML.Bin = []byte(ValidJSONTemplate)
	th.AssertNoErr(t, templateYAML.Parse())
	th.AssertDeepEquals(t, ValidJSONTemplateParsed, templateYAML.Parsed)

	templateInvalid := new(Template)
	templateInvalid.Bin = []byte("Keep Austin Weird")
	err := templateInvalid.Parse()
	if err == nil {
		t.Error("Template parsing did not catch invalid template")
	}
}

func TestIgnoreIfTemplate(t *testing.T) {
	var keyValueTests = []struct {
		key   string
		value interface{}
		out   bool
	}{
		{"not_get_file", "afksdf", true},
		{"not_type", "sdfd", true},
		{"get_file", "shdfuisd", false},
		{"type", "dfsdfsd", true},
		{"type", "sdfubsduf.yaml", false},
		{"type", "sdfsdufs.template", false},
		{"type", "sdfsdf.file", true},
		{"type", map[string]string{"key": "value"}, true},
	}
	var result bool
	for _, kv := range keyValueTests {
		result = ignoreIfTemplate(kv.key, kv.value)
		if result != kv.out {
			t.Errorf("key: %v, value: %v expected: %v, actual: %v", kv.key, kv.value, result, kv.out)
		}
	}
}

const myNovaContent = `heat_template_version: 2014-10-16
parameters:
  flavor:
    type: string
    description: Flavor for the server to be created
    default: 4353
    hidden: true
resources:
  test_server:
    type: "OS::Nova::Server"
    properties:
      name: test-server
      flavor: 2 GB General Purpose v1
      image: Debian 7 (Wheezy) (PVHVM)
      networks:
      - {uuid: 11111111-1111-1111-1111-111111111111}`

func TestGetFileContentsWithType(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	baseurl, err := getBasePath()
	th.AssertNoErr(t, err)
	fakeURL := strings.Join([]string{baseurl, "my_nova.yaml"}, "/")
	urlparsed, err := url.Parse(fakeURL)
	th.AssertNoErr(t, err)
	th.Mux.HandleFunc(urlparsed.Path, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, myNovaContent)
	})

	client := fakeClient{BaseClient: getHTTPClient()}
	te := new(Template)
	te.Bin = []byte(`heat_template_version: 2015-04-30
resources:
  my_server:
    type: my_nova.yaml`)
	te.client = client

	th.AssertNoErr(t, te.Parse())
	th.AssertNoErr(t, te.getFileContents(te.Parsed, ignoreIfTemplate, true))

	expectedFiles := map[string]string{
		"my_nova.yaml": `heat_template_version: 2014-10-16
parameters:
  flavor:
    type: string
    description: Flavor for the server to be created
    default: 4353
    hidden: true
resources:
  test_server:
    type: "OS::Nova::Server"
    properties:
      name: test-server
      flavor: 2 GB General Purpose v1
      image: Debian 7 (Wheezy) (PVHVM)
      networks:
      - {uuid: 11111111-1111-1111-1111-111111111111}`}
	th.AssertEquals(t, expectedFiles["my_nova.yaml"], te.Files[fakeURL])
	te.fixFileRefs()
	expectedParsed := map[string]interface{}{
		"heat_template_version": "2015-04-30",
		"resources": map[string]interface{}{
			"my_server": map[string]interface{}{
				"type": fakeURL,
			},
		},
	}
	th.AssertNoErr(t, te.Parse())
	th.AssertDeepEquals(t, expectedParsed, te.Parsed)
}

func TestGetFileContentsWithFile(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	baseurl, err := getBasePath()
	th.AssertNoErr(t, err)
	fakeURL := strings.Join([]string{baseurl, "somefile"}, "/")
	urlparsed, err := url.Parse(fakeURL)
	th.AssertNoErr(t, err)
	somefile := `Welcome!`
	th.Mux.HandleFunc(urlparsed.Path, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, somefile)
	})

	client := fakeClient{BaseClient: getHTTPClient()}
	te := new(Template)
	te.Bin = []byte(`heat_template_version: 2015-04-30
resources:
  test_resource:
    type: OS::Heat::TestResource
    properties:
      value: {get_file: somefile }`)
	te.client = client

	th.AssertNoErr(t, te.Parse())
	th.AssertNoErr(t, te.getFileContents(te.Parsed, ignoreIfTemplate, true))
	expectedFiles := map[string]string{
		"somefile": "Welcome!",
	}
	th.AssertEquals(t, expectedFiles["somefile"], te.Files[fakeURL])
	te.fixFileRefs()
	expectedParsed := map[string]interface{}{
		"heat_template_version": "2015-04-30",
		"resources": map[string]interface{}{
			"test_resource": map[string]interface{}{
				"type": "OS::Heat::TestResource",
				"properties": map[string]interface{}{
					"value": map[string]interface{}{
						"get_file": fakeURL,
					},
				},
			},
		},
	}
	th.AssertNoErr(t, te.Parse())
	th.AssertDeepEquals(t, expectedParsed, te.Parsed)
}

func TestGetFileContentsComposeRelativePath(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	baseurl, err := getBasePath()
	th.AssertNoErr(t, err)

	novaPath := strings.Join([]string{"templates", "my_nova.yaml"}, "/")
	novaURL := strings.Join([]string{baseurl, novaPath}, "/")
	novaURLParse, err := url.Parse(novaURL)
	th.AssertNoErr(t, err)
	th.Mux.HandleFunc(novaURLParse.Path, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, myNovaContent)
	})

	subStacksPath := strings.Join([]string{"substacks", "my_substack.yaml"}, "/")
	subStackURL := strings.Join([]string{baseurl, subStacksPath}, "/")
	subStackURLParsed, err := url.Parse(subStackURL)
	th.AssertNoErr(t, err)
	mySubStackContentFmt := `heat_template_version: 2015-04-30
resources:
  my_server:
    type: %s
  my_backend:
    type: "OS::Nova::Server"
    properties:
      name: test-backend
      flavor: 4 GB General Purpose v1
      image: Debian 7 (Wheezy) (PVHVM)
      networks:
      - {uuid: 11111111-1111-1111-1111-111111111111}`
	th.Mux.HandleFunc(subStackURLParsed.Path, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, mySubStackContentFmt, "../templates/my_nova.yaml")
	})

	client := fakeClient{BaseClient: getHTTPClient()}
	te := new(Template)
	te.Bin = []byte(`heat_template_version: 2015-04-30
resources:
  my_stack:
    type: substacks/my_substack.yaml`)
	te.client = client

	th.AssertNoErr(t, te.Parse())
	th.AssertNoErr(t, te.getFileContents(te.Parsed, ignoreIfTemplate, true))

	expectedFiles := map[string]string{
		"templates/my_nova.yaml":     myNovaContent,
		"substacks/my_substack.yaml": fmt.Sprintf(mySubStackContentFmt, novaURLParse),
	}
	th.AssertEquals(t, expectedFiles[novaPath], te.Files[novaURL])
	th.AssertEquals(t, expectedFiles[subStacksPath], te.Files[subStackURL])

	te.fixFileRefs()
	expectedParsed := map[string]interface{}{
		"heat_template_version": "2015-04-30",
		"resources": map[string]interface{}{
			"my_stack": map[string]interface{}{
				"type": subStackURL,
			},
		},
	}
	th.AssertNoErr(t, te.Parse())
	th.AssertDeepEquals(t, expectedParsed, te.Parsed)
}
