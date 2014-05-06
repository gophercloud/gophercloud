package snapshots

type Snapshot struct {
	Status              string
	Display_name        string
	Created_at          string
	Display_description string
	Volume_id           string
	Metadata            map[string]string
	Id                  string
	Size                int
}

type CreateOpts map[string]interface{}
type GetOpts map[string]string
type DeleteOpts map[string]string
