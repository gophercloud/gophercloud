package stacks

import (
	"strings"
	"testing"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
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
		value any
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

var myNovaExpected = map[string]any{
	"heat_template_version": "2014-10-16",
	"parameters": map[string]any{
		"flavor": map[string]any{
			"type":        "string",
			"description": "Flavor for the server to be created",
			"default":     4353,
			"hidden":      true,
		},
	},
	"resources": map[string]any{
		"test_server": map[string]any{
			"type": "OS::Nova::Server",
			"properties": map[string]any{
				"name":   "test-server",
				"flavor": "2 GB General Purpose v1",
				"image":  "Debian 7 (Wheezy) (PVHVM)",
				"networks": []any{
					map[string]any{
						"uuid": "11111111-1111-1111-1111-111111111111",
					},
				},
			},
		},
	},
}

func TestGetFileContentsWithType(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	baseurl, err := getBasePath()
	th.AssertNoErr(t, err)

	fakeURL := th.ServeFile(t, baseurl, "my_nova.yaml", "application/json", myNovaContent)

	client := fakeClient{BaseClient: getHTTPClient()}
	te := new(Template)
	te.Bin = []byte(`heat_template_version: 2015-04-30
resources:
  my_server:
    type: my_nova.yaml`)
	te.client = client

	th.AssertNoErr(t, te.Parse())
	th.AssertNoErr(t, te.getFileContents(te.Parsed, ignoreIfTemplate, true))

	// Now check template and referenced file
	expectedParsed := map[string]any{
		"heat_template_version": "2015-04-30",
		"resources": map[string]any{
			"my_server": map[string]any{
				"type": fakeURL,
			},
		},
	}
	th.AssertNoErr(t, te.Parse())
	th.AssertDeepEquals(t, expectedParsed, te.Parsed)

	novaTe := new(Template)
	novaTe.Bin = []byte(te.Files[fakeURL])
	th.AssertNoErr(t, novaTe.Parse())
	th.AssertDeepEquals(t, myNovaExpected, novaTe.Parsed)
}

func TestGetFileContentsWithFile(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	baseurl, err := getBasePath()
	th.AssertNoErr(t, err)

	somefile := `Welcome!`
	fakeURL := th.ServeFile(t, baseurl, "somefile", "text/plain", somefile)

	client := fakeClient{BaseClient: getHTTPClient()}
	te := new(Template)
	// Note: We include the path that should be replaced also as a not-to-be-replaced
	// keyword ("path: somefile" below) to validate that no updates happen outside of
	// the real local URLs (child templates (type:) or included files (get_file:)).
	te.Bin = []byte(`heat_template_version: 2015-04-30
resources:
  test_resource:
    type: OS::Heat::TestResource
    properties:
      path: somefile
      value: { get_file: somefile }`)
	te.client = client

	th.AssertNoErr(t, te.Parse())
	th.AssertNoErr(t, te.getFileContents(te.Parsed, ignoreIfTemplate, true))
	expectedFiles := map[string]string{
		"somefile": "Welcome!",
	}
	th.AssertEquals(t, expectedFiles["somefile"], te.Files[fakeURL])
	expectedParsed := map[string]any{
		"heat_template_version": "2015-04-30",
		"resources": map[string]any{
			"test_resource": map[string]any{
				"type": "OS::Heat::TestResource",
				"properties": map[string]any{
					"path": "somefile",
					"value": map[string]any{
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
	novaURL := th.ServeFile(t, baseurl, novaPath, "application/json", myNovaContent)

	mySubStackContent := `heat_template_version: 2015-04-30
resources:
  my_server:
    type: ../templates/my_nova.yaml
  my_backend:
    type: "OS::Nova::Server"
    properties:
      name: test-backend
      flavor: 4 GB General Purpose v1
      image: Debian 7 (Wheezy) (PVHVM)
      networks:
      - {uuid: 11111111-1111-1111-1111-111111111111}`
	mySubstrackExpected := map[string]any{
		"heat_template_version": "2015-04-30",
		"resources": map[string]any{
			"my_server": map[string]any{
				"type": novaURL,
			},
			"my_backend": map[string]any{
				"type": "OS::Nova::Server",
				"properties": map[string]any{
					"name":   "test-backend",
					"flavor": "4 GB General Purpose v1",
					"image":  "Debian 7 (Wheezy) (PVHVM)",
					"networks": []any{
						map[string]any{
							"uuid": "11111111-1111-1111-1111-111111111111",
						},
					},
				},
			},
		},
	}
	subStacksPath := strings.Join([]string{"substacks", "my_substack.yaml"}, "/")
	subStackURL := th.ServeFile(t, baseurl, subStacksPath, "application/json", mySubStackContent)

	client := fakeClient{BaseClient: getHTTPClient()}
	te := new(Template)
	te.Bin = []byte(`heat_template_version: 2015-04-30
resources:
  my_stack:
    type: substacks/my_substack.yaml`)
	te.client = client

	th.AssertNoErr(t, te.Parse())
	th.AssertNoErr(t, te.getFileContents(te.Parsed, ignoreIfTemplate, true))

	expectedParsed := map[string]any{
		"heat_template_version": "2015-04-30",
		"resources": map[string]any{
			"my_stack": map[string]any{
				"type": subStackURL,
			},
		},
	}
	th.AssertNoErr(t, te.Parse())
	th.AssertDeepEquals(t, expectedParsed, te.Parsed)

	expectedFiles := map[string]any{
		novaURL:     myNovaExpected,
		subStackURL: mySubstrackExpected,
	}
	for path, expected := range expectedFiles {
		checkTe := new(Template)
		checkTe.Bin = []byte(te.Files[path])
		th.AssertNoErr(t, checkTe.Parse())
		th.AssertDeepEquals(t, expected, checkTe.Parsed)
	}
}
