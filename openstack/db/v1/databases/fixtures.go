package databases

var createDBsReq = `
{
	"databases": [
		{
			"character_set": "utf8",
			"collate": "utf8_general_ci",
			"name": "testingdb"
		},
		{
			"name": "sampledb"
		}
	]
}
`

var listDBsResp = `
{
	"databases": [
		{
			"name": "anotherexampledb"
		},
		{
			"name": "exampledb"
		},
		{
			"name": "nextround"
		},
		{
			"name": "sampledb"
		},
		{
			"name": "testingdb"
		}
	]
}
`
