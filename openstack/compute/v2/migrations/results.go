package migrations

import (
	"encoding/json"
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
	"time"
)

// Migration represents the details of a migration.
type Migration struct {
	// The date and time when the resource was created
	CreatedAt time.Time `json:"-"`
	// The target compute for a migration
	DestCompute string `json:"dest_compute"`
	// The target host for a migration
	DestHost string `json:"dest_host"`
	// The target node for a migration.
	DestNode string `json:"dest_node"`
	// The ID of the server migration
	Id int64 `json:"id"`
	// The UUID of the server
	InstanceUuid string `json:"instance_uuid"`
	// In resize case, the flavor ID for resizing the server
	NewInstanceTypeId int64 `json:"new_instance_type_id"`
	// The flavor ID of the server when the migration was started
	OldInstanceTypeId int64 `json:"old_instance_type_id"`
	// The source compute for a migration
	SourceCompute string `json:"source_compute"`
	// The source node for a migration
	SourceNode string `json:"source_node"`
	// The current status of the migration
	Status string `json:"status"`
	// The date and time when the resource was updated
	UpdatedAt time.Time `json:"-"`
	// The type of the server migration. This is one of live-migration, migration, resize and evacuation
	MigrationType string `json:"migration_type"`
	// The UUID of the migration
	Uuid string `json:"uuid"`
	// The ID of the user which initiated the server migration
	UserId string `json:"user_id"`
	// The ID of the user which initiated the server migration
	ProjectId string `json:"project_id"`
}

// UnmarshalJSON converts our JSON API response into our migration struct
func (i *Migration) UnmarshalJSON(b []byte) error {
	type tmp Migration
	var s struct {
		tmp
		CreatedAt gophercloud.JSONRFC3339MilliNoZ `json:"created_at"`
		UpdatedAt gophercloud.JSONRFC3339MilliNoZ `json:"updated_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*i = Migration(s.tmp)

	i.UpdatedAt = time.Time(s.UpdatedAt)
	i.CreatedAt = time.Time(s.CreatedAt)
	return err
}

type MigrationPage struct {
	pagination.SinglePageBase
}

func (r MigrationPage) IsEmpty() (bool, error) {
	migrations, err := ExtractMigrations(r)
	return len(migrations) == 0, err
}

func ExtractMigrations(r pagination.Page) ([]Migration, error) {
	var resp []Migration
	err := ExtractMigrationsInto(r, &resp)
	return resp, err
}

func ExtractMigrationsInto(r pagination.Page, v any) error {
	return r.(MigrationPage).Result.ExtractIntoSlicePtr(v, "migrations")
}
