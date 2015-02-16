package instances

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

const singleInstanceJson = `
{
  "instance": {
    "created": "2014-02-13T21:47:13",
    "datastore": {
      "type": "mysql",
      "version": "5.6"
    },
    "flavor": {
      "id": "1",
      "links": [
        {
          "href": "https://ord.databases.api.rackspacecloud.com/v1.0/1234/flavors/1",
          "rel": "self"
        },
        {
          "href": "https://ord.databases.api.rackspacecloud.com/v1.0/1234/flavors/1",
          "rel": "bookmark"
        }
      ]
    },
    "links": [
      {
        "href": "https://ord.databases.api.rackspacecloud.com/v1.0/1234/flavors/1",
        "rel": "self"
      }
    ],
    "hostname": "e09ad9a3f73309469cf1f43d11e79549caf9acf2.rackspaceclouddb.com",
    "id": "d4603f69-ec7e-4e9b-803f-600b9205576f",
    "name": "json_rack_instance",
    "status": "BUILD",
    "updated": "2014-02-13T21:47:13",
    "volume": {
      "size": 2
    }
  }
}
`

func HandleCreateInstanceSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/instances", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		th.TestJSONRequest(t, r, `
{
  "instance": {
    "databases": [
      {
        "character_set": "utf8",
        "collate": "utf8_general_ci",
        "name": "sampledb"
      },
      {
        "name": "nextround"
      }
    ],
    "flavorRef": "1",
    "name": "json_rack_instance",
    "users": [
      {
        "databases": [
          {
            "name": "sampledb"
          }
        ],
        "name": "demouser",
        "password": "demopassword"
      }
    ],
    "volume": {
      "size": 2
    },
		"restorePoint": {
			"backupRef": "1234567890"
		}
  }
}
`)

		fmt.Fprintf(w, singleInstanceJson)
	})
}

func HandleGetInstanceSuccessfully(t *testing.T, id string) {
	th.Mux.HandleFunc("/instances/"+id, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")

		fmt.Fprintf(w, singleInstanceJson)
	})
}

func HandleGetConfigSuccessfully(t *testing.T, id string) {
	th.Mux.HandleFunc("/instances/"+id+"/configuration", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")

		fmt.Fprintf(w, `
{
  "instance": {
    "configuration": {
      "basedir": "/usr",
      "connect_timeout": "15",
      "datadir": "/var/lib/mysql",
      "default_storage_engine": "innodb",
      "innodb_buffer_pool_instances": "1",
      "innodb_buffer_pool_size": "175M",
      "innodb_checksum_algorithm": "crc32",
      "innodb_data_file_path": "ibdata1:10M:autoextend",
      "innodb_file_per_table": "1",
      "innodb_io_capacity": "200",
      "innodb_log_file_size": "256M",
      "innodb_log_files_in_group": "2",
      "innodb_open_files": "8192",
      "innodb_thread_concurrency": "0",
      "join_buffer_size": "1M",
      "key_buffer_size": "50M",
      "local-infile": "0",
      "log-error": "/var/log/mysql/mysqld.log",
      "max_allowed_packet": "16M",
      "max_connect_errors": "10000",
      "max_connections": "40",
      "max_heap_table_size": "16M",
      "myisam-recover": "BACKUP",
      "open_files_limit": "8192",
      "performance_schema": "off",
      "pid_file": "/var/run/mysqld/mysqld.pid",
      "port": "3306",
      "query_cache_limit": "1M",
      "query_cache_size": "8M",
      "query_cache_type": "1",
      "read_buffer_size": "256K",
      "read_rnd_buffer_size": "1M",
      "server_id": "1",
      "skip-external-locking": "1",
      "skip_name_resolve": "1",
      "sort_buffer_size": "256K",
      "table_open_cache": "4096",
      "thread_stack": "192K",
      "tmp_table_size": "16M",
      "tmpdir": "/var/tmp",
      "user": "mysql",
      "wait_timeout": "3600"
    }
  }
}
`)
	})
}

func HandleAssociateGroupSuccessfully(t *testing.T, id string) {
	th.Mux.HandleFunc("/instances/"+id, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestJSONRequest(t, r, `{"instance": {"configuration": "{configGroupID}"}}`)

		w.WriteHeader(http.StatusAccepted)
		w.Header().Add("Content-Type", "application/json")

		fmt.Fprintf(w, singleInstanceJson)
	})
}

func HandleListBackupsSuccessfully(t *testing.T, id string) {
	th.Mux.HandleFunc("/instances/"+id+"/backups", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")

		fmt.Fprintf(w, `
{
  "backups": [
    {
      "status": "COMPLETED",
      "updated": "2014-06-18T21:24:39",
      "description": "Backup from Restored Instance",
      "datastore": {
        "version": "5.1",
        "type": "MySQL",
        "version_id": "20000000-0000-0000-0000-000000000002"
      },
      "id": "87972694-4be2-40f5-83f8-501656e0032a",
      "size": 0.141026,
      "name": "restored_backup",
      "created": "2014-06-18T21:23:35",
      "instance_id": "29af2cd9-0674-48ab-b87a-b160f00208e6",
      "parent_id": null,
      "locationRef": "http://localhost/path/to/backup"
    }
  ]
}
`)
	})
}
