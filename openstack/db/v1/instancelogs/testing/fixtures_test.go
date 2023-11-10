package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/db/v1/instancelogs"
	"github.com/gophercloud/gophercloud/testhelper/fixture"
)

var (
	showReq    = `{"name": "general"}`
	enableReq  = `{"name": "general", "enable": 1}`
	disableReq = `{"name": "general", "disable": 1}`
	publishReq = `{"name": "general", "publish": 1}`
	discardReq = `{"name": "general", "discard": 1}`
)

var (
	logName    = "general"
	instanceID = "{instanceID}"
	rootURL    = "/instances/" + instanceID + "/log"
)

var listResp = `
{
    "logs": [
        {
            "name": "general",
            "type": "USER",
            "status": "Partial",
            "published": 128,
            "pending": 4096,
            "container": "data_logs",
            "prefix": "5e9e616c-1827-45f5-a487-679084d82f7e/mysql-general/",
            "metafile": "5e9e616c-1827-45f5-a487-679084d82f7e/mysql-general_metafile"
        }
    ]
}
`

var showResp = `
{
    "log": {
        "name": "general",
        "type": "USER",
        "status": "Partial",
        "published": 128,
        "pending": 4096,
        "container": "data_logs",
        "prefix": "5e9e616c-1827-45f5-a487-679084d82f7e/mysql-general/",
        "metafile": "5e9e616c-1827-45f5-a487-679084d82f7e/mysql-general_metafile"
    }
}
`

var enableResp = `
{
    "log": {
        "name": "general",
        "type": "USER",
        "status": "Enabled",
        "published": 0,
        "pending": 0,
        "container": "None",
        "prefix": "None",
        "metafile": "5e9e616c-1827-45f5-a487-679084d82f7e/mysql-general_metafile"
    }
}
`

var disableResp = `
{
    "log": {
        "name": "general",
        "type": "USER",
        "status": "Disabled",
        "published": 4096,
        "pending": 0,
        "container": "data_logs",
        "prefix": "5e9e616c-1827-45f5-a487-679084d82f7e/mysql-general/",
        "metafile": "5e9e616c-1827-45f5-a487-679084d82f7e/mysql-general_metafile"
    }
}
`

var publishResp = `
{
    "log": {
        "name": "general",
        "type": "USER",
        "status": "Published",
        "published": 128,
        "pending": 0,
        "container": "data_logs",
        "prefix": "5e9e616c-1827-45f5-a487-679084d82f7e/mysql-general/",
        "metafile": "5e9e616c-1827-45f5-a487-679084d82f7e/mysql-general_metafile"
    }
}
`

var discardResp = `
{
    "log": {
        "name": "general",
        "type": "USER",
        "status": "Ready",
        "published": 0,
        "pending": 128,
        "container": "None",
        "prefix": "None",
        "metafile": "5e9e616c-1827-45f5-a487-679084d82f7e/mysql-general_metafile"
    }
}
`

var expectedLog = instancelogs.Log{
	Name:      logName,
	Type:      "USER",
	Status:    "Partial",
	Published: 128,
	Pending:   4096,
	Container: "data_logs",
	Prefix:    "5e9e616c-1827-45f5-a487-679084d82f7e/mysql-general/",
	Metafile:  "5e9e616c-1827-45f5-a487-679084d82f7e/mysql-general_metafile",
}

var expectedShowLog = instancelogs.Log{
	Name:      logName,
	Type:      "USER",
	Status:    "Partial",
	Published: 128,
	Pending:   4096,
	Container: "data_logs",
	Prefix:    "5e9e616c-1827-45f5-a487-679084d82f7e/mysql-general/",
	Metafile:  "5e9e616c-1827-45f5-a487-679084d82f7e/mysql-general_metafile",
}

var expectedEnableLog = instancelogs.Log{
	Name:      logName,
	Type:      "USER",
	Status:    "Enabled",
	Published: 0,
	Pending:   0,
	Container: "None",
	Prefix:    "None",
	Metafile:  "5e9e616c-1827-45f5-a487-679084d82f7e/mysql-general_metafile",
}

var expectedDisableLog = instancelogs.Log{
	Name:      logName,
	Type:      "USER",
	Status:    "Disabled",
	Published: 4096,
	Pending:   0,
	Container: "data_logs",
	Prefix:    "5e9e616c-1827-45f5-a487-679084d82f7e/mysql-general/",
	Metafile:  "5e9e616c-1827-45f5-a487-679084d82f7e/mysql-general_metafile",
}

var expectedPublishLog = instancelogs.Log{
	Name:      logName,
	Type:      "USER",
	Status:    "Published",
	Published: 128,
	Pending:   0,
	Container: "data_logs",
	Prefix:    "5e9e616c-1827-45f5-a487-679084d82f7e/mysql-general/",
	Metafile:  "5e9e616c-1827-45f5-a487-679084d82f7e/mysql-general_metafile",
}

var expectedDiscardLog = instancelogs.Log{
	Name:      logName,
	Type:      "USER",
	Status:    "Ready",
	Published: 0,
	Pending:   128,
	Container: "None",
	Prefix:    "None",
	Metafile:  "5e9e616c-1827-45f5-a487-679084d82f7e/mysql-general_metafile",
}

func HandleList(t *testing.T) {
	fixture.SetupHandler(t, rootURL, "GET", "", listResp, 200)
}

func HandleShow(t *testing.T) {
	fixture.SetupHandler(t, rootURL, "POST", showReq, showResp, 200)
}

func HandleEnable(t *testing.T) {
	fixture.SetupHandler(t, rootURL, "POST", enableReq, enableResp, 200)
}

func HandleDisable(t *testing.T) {
	fixture.SetupHandler(t, rootURL, "POST", disableReq, disableResp, 200)
}

func HandlePublish(t *testing.T) {
	fixture.SetupHandler(t, rootURL, "POST", publishReq, publishResp, 200)
}

func HandleDiscard(t *testing.T) {
	fixture.SetupHandler(t, rootURL, "POST", discardReq, discardResp, 200)
}
