package testing

import (
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/datastores"
)

const version1JSON = `
{
	"id": "b00000b0-00b0-0b00-00b0-000b000000bb",
	"links": [
		{
			"href": "https://10.240.28.38:8779/v1.0/1234/datastores/versions/b00000b0-00b0-0b00-00b0-000b000000bb",
			"rel": "self"
		},
		{
			"href": "https://10.240.28.38:8779/datastores/versions/b00000b0-00b0-0b00-00b0-000b000000bb",
			"rel": "bookmark"
		}
	],
	"name": "5.1"
}
`

const version2JSON = `
{
	"id": "c00000b0-00c0-0c00-00c0-000b000000cc",
	"links": [
		{
			"href": "https://10.240.28.38:8779/v1.0/1234/datastores/versions/c00000b0-00c0-0c00-00c0-000b000000cc",
			"rel": "self"
		},
		{
			"href": "https://10.240.28.38:8779/datastores/versions/c00000b0-00c0-0c00-00c0-000b000000cc",
			"rel": "bookmark"
		}
	],
	"name": "5.2"
}
`

const createVersionReq = `
{
	"version": {
		"datastore_name": "mysql",
		"datastore_manager": "mysql",
		"name": "mysql-5.7",
		"image_tags": ["trove", "mysql"],
		"active": true,
		"default": false,
		"packages": ["mysql-server"],
		"version": "5.7.30"
	}
}
`

const createdVersionJSON = `
{
	"active": true,
	"datastore_id": "10000000-0000-0000-0000-000000000001",
	"datastore_manager": "mysql",
	"datastore_name": "mysql",
	"default": false,
	"id": "b00000b0-00b0-0b00-00b0-000b000000bb",
	"image": "",
	"image_tags": ["trove", "mysql"],
	"name": "mysql-5.7",
	"packages": ["mysql-server"],
	"version": "5.7.30"
}
`

const updateVersionReq = `
{
	"name": "mysql-5.7-updated",
	"active": true
}
`

const updatedVersionJSON = `
{
	"active": true,
	"datastore_id": "10000000-0000-0000-0000-000000000001",
	"datastore_manager": "mysql",
	"datastore_name": "mysql",
	"default": false,
	"id": "b00000b0-00b0-0b00-00b0-000b000000bb",
	"image": "",
	"image_tags": ["trove", "mysql"],
	"name": "mysql-5.7-updated",
	"packages": ["mysql-server"],
	"version": "5.7.30"
}
`

var versionsJSON = fmt.Sprintf(`"versions": [%s, %s]`, version1JSON, version2JSON)

var singleDSJSON = fmt.Sprintf(`
{
  "default_version": "c00000b0-00c0-0c00-00c0-000b000000cc",
  "id": "10000000-0000-0000-0000-000000000001",
  "links": [
    {
      "href": "https://10.240.28.38:8779/v1.0/1234/datastores/10000000-0000-0000-0000-000000000001",
      "rel": "self"
    },
    {
      "href": "https://10.240.28.38:8779/datastores/10000000-0000-0000-0000-000000000001",
      "rel": "bookmark"
    }
  ],
  "name": "mysql",
  %s
}
`, versionsJSON)

var (
	ListDSResp          = fmt.Sprintf(`{"datastores":[%s]}`, singleDSJSON)
	GetDSResp           = fmt.Sprintf(`{"datastore":%s}`, singleDSJSON)
	ListVersionsResp    = fmt.Sprintf(`{%s}`, versionsJSON)
	GetVersionResp      = fmt.Sprintf(`{"version":%s}`, version1JSON)
	CreateVersionResp   = fmt.Sprintf(`{"version":%s}`, createdVersionJSON)
	ListAllVersionsResp = fmt.Sprintf(`{"versions":[%s]}`, createdVersionJSON)
	UpdateVersionResp   = fmt.Sprintf(`{"version":%s}`, updatedVersionJSON)
)

var ExampleVersion1 = datastores.Version{
	ID: "b00000b0-00b0-0b00-00b0-000b000000bb",
	Links: []gophercloud.Link{
		{Rel: "self", Href: "https://10.240.28.38:8779/v1.0/1234/datastores/versions/b00000b0-00b0-0b00-00b0-000b000000bb"},
		{Rel: "bookmark", Href: "https://10.240.28.38:8779/datastores/versions/b00000b0-00b0-0b00-00b0-000b000000bb"},
	},
	Name: "5.1",
}

var exampleVersion2 = datastores.Version{
	ID: "c00000b0-00c0-0c00-00c0-000b000000cc",
	Links: []gophercloud.Link{
		{Rel: "self", Href: "https://10.240.28.38:8779/v1.0/1234/datastores/versions/c00000b0-00c0-0c00-00c0-000b000000cc"},
		{Rel: "bookmark", Href: "https://10.240.28.38:8779/datastores/versions/c00000b0-00c0-0c00-00c0-000b000000cc"},
	},
	Name: "5.2",
}

var ExampleVersions = []datastores.Version{ExampleVersion1, exampleVersion2}

var ExampleCreatedVersion = datastores.Version{
	Active:           true,
	DatastoreID:      "10000000-0000-0000-0000-000000000001",
	DatastoreManager: "mysql",
	DatastoreName:    "mysql",
	ID:               "b00000b0-00b0-0b00-00b0-000b000000bb",
	ImageTags:        []string{"trove", "mysql"},
	Name:             "mysql-5.7",
	Packages:         datastores.PackageList{"mysql-server"},
	Version:          "5.7.30",
}

var ExampleUpdatedVersion = datastores.Version{
	Active:           true,
	DatastoreID:      "10000000-0000-0000-0000-000000000001",
	DatastoreManager: "mysql",
	DatastoreName:    "mysql",
	ID:               "b00000b0-00b0-0b00-00b0-000b000000bb",
	ImageTags:        []string{"trove", "mysql"},
	Name:             "mysql-5.7-updated",
	Packages:         datastores.PackageList{"mysql-server"},
	Version:          "5.7.30",
}

var ExampleDatastore = datastores.Datastore{
	DefaultVersion: "c00000b0-00c0-0c00-00c0-000b000000cc",
	ID:             "10000000-0000-0000-0000-000000000001",
	Links: []gophercloud.Link{
		{Rel: "self", Href: "https://10.240.28.38:8779/v1.0/1234/datastores/10000000-0000-0000-0000-000000000001"},
		{Rel: "bookmark", Href: "https://10.240.28.38:8779/datastores/10000000-0000-0000-0000-000000000001"},
	},
	Name:     "mysql",
	Versions: ExampleVersions,
}
