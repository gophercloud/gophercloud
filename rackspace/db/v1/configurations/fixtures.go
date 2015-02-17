package configurations

import "fmt"

const singleConfigJSON = `
{
  "created": "2014-07-31T18:56:09",
  "datastore_name": "mysql",
  "datastore_version_id": "b00000b0-00b0-0b00-00b0-000b000000bb",
  "datastore_version_name": "5.6",
  "description": "example_description",
  "id": "005a8bb7-a8df-40ee-b0b7-fc144641abc2",
  "name": "example-configuration-name",
  "updated": "2014-07-31T18:56:09"
}
`

const singleConfigWithValuesJSON = `
{
  "created": "2014-07-31T15:02:52",
  "datastore_name": "mysql",
  "datastore_version_id": "b00000b0-00b0-0b00-00b0-000b000000bb",
  "datastore_version_name": "5.6",
  "description": "example description",
  "id": "005a8bb7-a8df-40ee-b0b7-fc144641abc2",
  "instance_count": 0,
  "name": "example-configuration-name",
  "updated": "2014-07-31T15:02:52",
  "values": {
    "collation_server": "latin1_swedish_ci",
    "connect_timeout": 120
  }
}
`

var (
	listConfigsJSON  = fmt.Sprintf(`{"configurations": [%s]}`, singleConfigJSON)
	getConfigJSON    = fmt.Sprintf(`{"configuration": %s}`, singleConfigJSON)
	createConfigJSON = fmt.Sprintf(`{"configuration": %s}`, singleConfigWithValuesJSON)
)

var createReq = `
{
  "configuration": {
    "datastore": {
      "type": "a00000a0-00a0-0a00-00a0-000a000000aa",
      "version": "b00000b0-00b0-0b00-00b0-000b000000bb"
    },
    "description": "example description",
    "name": "example-configuration-name",
    "values": {
      "collation_server": "latin1_swedish_ci",
      "connect_timeout": 120
    }
  }
}
`

var updateReq = `
{
  "configuration": {
    "values": {
      "connect_timeout": 300
    }
  }
}
`

var listInstancesJSON = `
{
  "instances": [
    {
      "id": "d4603f69-ec7e-4e9b-803f-600b9205576f",
      "name": "json_rack_instance"
    }
  ]
}
`

var exampleConfig = Config{
	Created:              "2014-07-31T18:56:09",
	DatastoreName:        "mysql",
	DatastoreVersionID:   "b00000b0-00b0-0b00-00b0-000b000000bb",
	DatastoreVersionName: "5.6",
	Description:          "example_description",
	ID:                   "005a8bb7-a8df-40ee-b0b7-fc144641abc2",
	Name:                 "example-configuration-name",
	Updated:              "2014-07-31T18:56:09",
}

var exampleConfigWithValues = Config{
	Created:              "2014-07-31T15:02:52",
	DatastoreName:        "mysql",
	DatastoreVersionID:   "b00000b0-00b0-0b00-00b0-000b000000bb",
	DatastoreVersionName: "5.6",
	Description:          "example description",
	ID:                   "005a8bb7-a8df-40ee-b0b7-fc144641abc2",
	Name:                 "example-configuration-name",
	Updated:              "2014-07-31T15:02:52",
	Values: map[string]interface{}{
		"collation_server": "latin1_swedish_ci",
		"connect_timeout":  120,
	},
}
