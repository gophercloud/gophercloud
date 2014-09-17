package images

import (
	"encoding/json"
	"testing"
)

const (
	// This example was taken from: http://docs.openstack.org/api/openstack-compute/2/content/Rebuild_Server-d1e3538.html

	simpleImageJSON = `{
		"id": "52415800-8b69-11e0-9b19-734f6f006e54",
		"name": "CentOS 5.2",
		"links": [{
			"rel": "self",
			"href": "http://servers.api.openstack.org/v2/1234/images/52415800-8b69-11e0-9b19-734f6f006e54"
		},{
			"rel": "bookmark",
			"href": "http://servers.api.openstack.org/1234/images/52415800-8b69-11e0-9b19-734f6f006e54"
		}]
	}`
)

func TestExtractImage(t *testing.T) {
	var simpleImageMap map[string]interface{}
	err := json.Unmarshal([]byte(simpleImageJSON), &simpleImageMap)
	if err != nil {
		t.Fatal(err)
	}

	image, err := ExtractImage(simpleImageMap)
	if err != nil {
		t.Fatal(err)
	}

	if image.Id != "52415800-8b69-11e0-9b19-734f6f006e54" {
		t.Fatal("I expected an image ID of 52415800-8b69-11e0-9b19-734f6f006e54; got " + image.Id)
	}

	if image.Name != "CentOS 5.2" {
		t.Fatal("I expected an image name of CentOS 5.2; got " + image.Name)
	}
}
