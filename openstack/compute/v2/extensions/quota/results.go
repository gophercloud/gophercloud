package quotas

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// Quota is a set of operational limits that allow for control of compute usage.
const sample = `
{
  "quota_set" : {
	"fixed_ips" : -1,
	"security_groups" : 10,
	"id" : "56b6c3eb639e48c691052919e5a60dc3",
	"injected_files" : 5,
	"injected_file_path_bytes" : 255,
	"cores" : 108,
	"security_group_rules" : 20,
	"keypairs" : 10,
	"instances" : 25,
	"ram" : 204800,
	"metadata_items" : 128,
	"injected_file_content_bytes" : 10240
  }
}
`

type Quota struct {
	//ID is tenant associated with this quota_set
	ID string `mapstructure:"id"`
	//FixedIps is number of fixed ips alloted this quota_set
	FixedIps int `mapstructure:"fixed_ips"`
	// FloatingIps is number of floatinh ips alloted this quota_set
	FloatingIps int `mapstructure:"floating_ips"`
	// InjectedFileContentBytes is content bytes allowed for each injected file
	InjectedFileContentBytes int `mapstructure:"injected_file_content_bytes"`
	// InjectedFilePathBytes is allowed bytes for each injected file path
	InjectedFilePathBytes int `mapstructure:"injected_file_path_bytes"`
	// InjectedFiles is injected files allowed for each project
	InjectedFiles int `mapstructure:"injected_files"`
	// KeyPairs is number of ssh keypairs
	KeyPairs int `mapstructure:"keypairs"`
	// MetadataItems is number of metadata items allowed for each instance
	MetadataItems int `mapstructure:"metadata_items"`
	// Ram is megabytes allowed for each instance
	Ram int `mapstructure:"ram"`
	// SecurityGroupRules is rules allowed for each security group
	SecurityGroupRules int `mapstructure:"security_group_rules"`
	// SecurityGroups security groups allowed for each project
	SecurityGroups int `mapstructure:"security_groups"`
	// Cores is number of instance cores allowed for each project
	Cores int `mapstructure:"cores"`
	// Instances is number of instances allowed for each project
	Instances int `mapstructure:"instances"`
}

// QuotaPage stores a single, only page of Quota results from a List call.
type QuotaPage struct {
	pagination.SinglePageBase
}

// IsEmpty determines whether or not a QuotaPage is empty.
func (page QuotaPage) IsEmpty() (bool, error) {
	ks, err := ExtractQuotas(page)
	return len(ks) == 0, err
}

// ExtractQuotas interprets a page of results as a slice of Quotas.
func ExtractQuotas(page pagination.Page) ([]Quota, error) {
	var resp struct {
		Quotas []Quota `mapstructure:"quotas"`
	}

	err := mapstructure.Decode(page.(QuotaPage).Body, &resp)
	results := make([]Quota, len(resp.Quotas))
	for i, q := range resp.Quotas {
		results[i] = q
	}
	return results, err
}

type quotaResult struct {
	gophercloud.Result
}

// Extract is a method that attempts to interpret any Quota resource response as a Quota struct.
func (r quotaResult) Extract() (*Quota, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var res struct {
		Quota *Quota `json:"quota_set" mapstructure:"quota_set"`
	}

	err := mapstructure.Decode(r.Body, &res)
	return res.Quota, err
}

// GetResult is the response from a Get operation. Call its Extract method to interpret it
// as a Quota.
type GetResult struct {
	quotaResult
}
