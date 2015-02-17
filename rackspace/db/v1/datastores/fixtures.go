package datastores

import (
	"fmt"

	"github.com/rackspace/gophercloud"
)

const version1JSON = `
{
	"id": "b00000b0-00b0-0b00-00b0-000b000000bb",
	"links": [
		{
			"href": "https://api.staging.ord1.clouddb.rackspace.net/v1.0/1234/datastores/versions/b00000b0-00b0-0b00-00b0-000b000000bb",
			"rel": "self"
		},
		{
			"href": "https://api.staging.ord1.clouddb.rackspace.net/datastores/versions/b00000b0-00b0-0b00-00b0-000b000000bb",
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
			"href": "https://api.staging.ord1.clouddb.rackspace.net/v1.0/1234/datastores/versions/c00000b0-00c0-0c00-00c0-000b000000cc",
			"rel": "self"
		},
		{
			"href": "https://api.staging.ord1.clouddb.rackspace.net/datastores/versions/c00000b0-00c0-0c00-00c0-000b000000cc",
			"rel": "bookmark"
		}
	],
	"name": "5.2"
}
`

var versionsJSON = fmt.Sprintf(`"versions": [%s, %s]`, version1JSON, version2JSON)

var singleDSJSON = fmt.Sprintf(`
{
  "default_version": "c00000b0-00c0-0c00-00c0-000b000000cc",
  "id": "10000000-0000-0000-0000-000000000001",
  "links": [
    {
      "href": "https://api.staging.ord1.clouddb.rackspace.net/v1.0/1234/datastores/10000000-0000-0000-0000-000000000001",
      "rel": "self"
    },
    {
      "href": "https://api.staging.ord1.clouddb.rackspace.net/datastores/10000000-0000-0000-0000-000000000001",
      "rel": "bookmark"
    }
  ],
  "name": "mysql",
  %s
}
`, versionsJSON)

var (
	listDSResp       = fmt.Sprintf(`{"datastores":[%s]}`, singleDSJSON)
	getDSResp        = fmt.Sprintf(`{"datastore":%s}`, singleDSJSON)
	listVersionsResp = fmt.Sprintf(`{%s}`, versionsJSON)
	getVersionResp   = fmt.Sprintf(`{"version":%s}`, version1JSON)
)

var exampleVersion1 = Version{
	ID: "b00000b0-00b0-0b00-00b0-000b000000bb",
	Links: []gophercloud.Link{
		gophercloud.Link{Rel: "self", Href: "https://api.staging.ord1.clouddb.rackspace.net/v1.0/1234/datastores/versions/b00000b0-00b0-0b00-00b0-000b000000bb"},
		gophercloud.Link{Rel: "bookmark", Href: "https://api.staging.ord1.clouddb.rackspace.net/datastores/versions/b00000b0-00b0-0b00-00b0-000b000000bb"},
	},
	Name: "5.1",
}

var exampleVersion2 = Version{
	ID: "c00000b0-00c0-0c00-00c0-000b000000cc",
	Links: []gophercloud.Link{
		gophercloud.Link{Rel: "self", Href: "https://api.staging.ord1.clouddb.rackspace.net/v1.0/1234/datastores/versions/c00000b0-00c0-0c00-00c0-000b000000cc"},
		gophercloud.Link{Rel: "bookmark", Href: "https://api.staging.ord1.clouddb.rackspace.net/datastores/versions/c00000b0-00c0-0c00-00c0-000b000000cc"},
	},
	Name: "5.2",
}

var exampleVersions = []Version{exampleVersion1, exampleVersion2}

var exampleDatastore = Datastore{
	DefaultVersion: "c00000b0-00c0-0c00-00c0-000b000000cc",
	ID:             "10000000-0000-0000-0000-000000000001",
	Links: []gophercloud.Link{
		gophercloud.Link{Rel: "self", Href: "https://api.staging.ord1.clouddb.rackspace.net/v1.0/1234/datastores/10000000-0000-0000-0000-000000000001"},
		gophercloud.Link{Rel: "bookmark", Href: "https://api.staging.ord1.clouddb.rackspace.net/datastores/10000000-0000-0000-0000-000000000001"},
	},
	Name:     "mysql",
	Versions: exampleVersions,
}
