package instances

import (
	"testing"

	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/rackspace/db/v1/backups"
	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

func TestGetConfig(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleGetConfigSuccessfully(t, instanceID)

	config, err := GetDefaultConfig(fake.ServiceClient(), instanceID).Extract()

	expected := map[string]string{
		"basedir":                      "/usr",
		"connect_timeout":              "15",
		"datadir":                      "/var/lib/mysql",
		"default_storage_engine":       "innodb",
		"innodb_buffer_pool_instances": "1",
		"innodb_buffer_pool_size":      "175M",
		"innodb_checksum_algorithm":    "crc32",
		"innodb_data_file_path":        "ibdata1:10M:autoextend",
		"innodb_file_per_table":        "1",
		"innodb_io_capacity":           "200",
		"innodb_log_file_size":         "256M",
		"innodb_log_files_in_group":    "2",
		"innodb_open_files":            "8192",
		"innodb_thread_concurrency":    "0",
		"join_buffer_size":             "1M",
		"key_buffer_size":              "50M",
		"local-infile":                 "0",
		"log-error":                    "/var/log/mysql/mysqld.log",
		"max_allowed_packet":           "16M",
		"max_connect_errors":           "10000",
		"max_connections":              "40",
		"max_heap_table_size":          "16M",
		"myisam-recover":               "BACKUP",
		"open_files_limit":             "8192",
		"performance_schema":           "off",
		"pid_file":                     "/var/run/mysqld/mysqld.pid",
		"port":                         "3306",
		"query_cache_limit":            "1M",
		"query_cache_size":             "8M",
		"query_cache_type":             "1",
		"read_buffer_size":             "256K",
		"read_rnd_buffer_size":         "1M",
		"server_id":                    "1",
		"skip-external-locking":        "1",
		"skip_name_resolve":            "1",
		"sort_buffer_size":             "256K",
		"table_open_cache":             "4096",
		"thread_stack":                 "192K",
		"tmp_table_size":               "16M",
		"tmpdir":                       "/var/tmp",
		"user":                         "mysql",
		"wait_timeout":                 "3600",
	}

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected, config)
}

func TestAssociateWithConfigGroup(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleAssociateGroupSuccessfully(t, instanceID)

	configGroupID := "{configGroupID}"
	res := AssociateWithConfigGroup(fake.ServiceClient(), instanceID, configGroupID)
	th.AssertNoErr(t, res.Err)
}

func TestListBackups(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleListBackupsSuccessfully(t, instanceID)
	count := 0

	ListBackups(fake.ServiceClient(), instanceID).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := backups.ExtractBackups(page)
		th.AssertNoErr(t, err)

		expected := []backups.Backup{
			backups.Backup{
				Created:     "2014-06-18T21:23:35",
				Description: "Backup from Restored Instance",
				ID:          "87972694-4be2-40f5-83f8-501656e0032a",
				InstanceID:  "29af2cd9-0674-48ab-b87a-b160f00208e6",
				LocationRef: "http://localhost/path/to/backup",
				Name:        "restored_backup",
				ParentID:    "",
				Size:        0.141026,
				Status:      "COMPLETED",
				Updated:     "2014-06-18T21:24:39",
				Datastore:   backups.Datastore{Version: "5.1", Type: "MySQL", VersionID: "20000000-0000-0000-0000-000000000002"},
			},
		}

		th.AssertDeepEquals(t, expected, actual)

		return true, nil
	})

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}
