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

func HandleCreateReplicaSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/instances", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestJSONRequest(t, r, `
{
  "instance": {
    "volume": {
      "size": 1
    },
    "flavorRef": "9",
    "name": "t2s1_ALT_GUEST",
    "replica_of": "6bdca2fc-418e-40bd-a595-62abda61862d"
  }
}
`)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
  "instance": {
    "status": "BUILD",
    "updated": "2014-10-14T18:42:15",
    "name": "t2s1_ALT_GUEST",
    "links": [
      {
        "href": "https://ord.databases.api.rackspacecloud.com/v1.0/5919009/instances/8367c312-7c40-4a66-aab1-5767478914fc",
        "rel": "self"
      },
      {
        "href": "https://ord.databases.api.rackspacecloud.com/instances/8367c312-7c40-4a66-aab1-5767478914fc",
        "rel": "bookmark"
      }
    ],
    "created": "2014-10-14T18:42:15",
    "id": "8367c312-7c40-4a66-aab1-5767478914fc",
    "volume": {"size": 1},
    "flavor": {"id": "9"},
    "datastore": {
      "version": "5.6",
      "type": "mysql"
    },
    "replica_of": {"id": "6bdca2fc-418e-40bd-a595-62abda61862d"}
  }
}
`)
	})
}

func HandleListReplicasSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/instances", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")

		fmt.Fprintf(w, `
{
	"instances": [
		{
			"status": "ACTIVE",
			"name": "t1s1_ALT_GUEST",
			"links": [
				{
					"href": "https://ord.databases.api.rackspacecloud.com/v1.0/1234/instances/3c691f06-bf9a-4618-b7ec-2817ce0cf254",
					"rel": "self"
				},
				{
					"href": "https://ord.databases.api.rackspacecloud.com/instances/3c691f06-bf9a-4618-b7ec-2817ce0cf254",
					"rel": "bookmark"
				}
			],
			"ip": [
				"10.0.0.3"
			],
			"id": "3c691f06-bf9a-4618-b7ec-2817ce0cf254",
			"volume": {
				"size": 1
			},
			"flavor": {
				"id": "9"
			},
			"datastore": {
				"version": "5.6",
				"type": "mysql"
			},
			"replica_of": {
				"id": "8b499b45-52d6-402d-b398-f9d8f279c69a"
			}
		}
	]
}
`)
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

func HandleGetReplicaSuccessfully(t *testing.T, id string) {
	th.Mux.HandleFunc("/instances/"+id, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")

		fmt.Fprintf(w, `
{
  "instance": {
    "status": "ACTIVE",
    "updated": "2014-09-26T19:15:57",
    "name": "t1_ALT_GUEST",
    "created": "2014-09-26T19:15:50",
    "ip": [
      "10.0.0.2"
    ],
    "replicas": [
			{"id": "3c691f06-bf9a-4618-b7ec-2817ce0cf254"}
    ],
    "id": "8b499b45-52d6-402d-b398-f9d8f279c69a",
    "volume": {
      "used": 0.54,
      "size": 1
    },
    "flavor": {"id": "9"},
    "datastore": {
      "version": "5.6",
      "type": "mysql"
    }
  }
}
`)
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

func HandleDetachReplicaSuccessfully(t *testing.T, id string) {
	th.Mux.HandleFunc("/instances/"+id, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		th.TestJSONRequest(t, r, `
{
	"instance": {
		"replica_of": "",
		"slave_of": ""
	}
}
`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
	})
}
