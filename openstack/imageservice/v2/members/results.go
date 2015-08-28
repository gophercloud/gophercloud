package members

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
)

// ImageMember model
type ImageMember struct {
	CreatedAt string `mapstructure:"created_at"`
	ImageID   string `mapstructure:"image_id"`
	MemberID  string `mapstructure:"member_id"`
	Schema    string
	// Status could be one of pending, accepted, reject
	Status    string
	UpdatedAt string `mapstructure:"updated_at"`
}

// CreateMemberResult result model
type CreateMemberResult struct {
	gophercloud.ErrResult
}

// Extract ImageMember model from request if possible
func (cm CreateMemberResult) Extract() (*ImageMember, error) {
	if cm.Err != nil {
		return nil, cm.Err
	}
	casted := cm.Body.(map[string]interface{})
	var results ImageMember

	err := mapstructure.Decode(casted, &results)
	return &results, err
}

// ListMembersResult model
type ListMembersResult struct {
	gophercloud.ErrResult
}

// Extract returns list of image members
func (lm ListMembersResult) Extract() (*[]ImageMember, error) {
	if lm.Err != nil {
		return nil, lm.Err
	}
	casted := lm.Body.(map[string]interface{})

	var results struct {
		ImageMembers []ImageMember `mapstructure:"members"`
	}

	err := mapstructure.Decode(casted, &results)
	return &results.ImageMembers, err
}

// MemberDetailsResult model
type MemberDetailsResult struct {
	gophercloud.ErrResult
}

// Extract returns image member details
func (md MemberDetailsResult) Extract() (*ImageMember, error) {
	if md.Err != nil {
		return nil, md.Err
	}
	casted := md.Body.(map[string]interface{})

	var results ImageMember

	err := mapstructure.Decode(casted, &results)
	return &results, err

}

// MemberDeleteResult model
type MemberDeleteResult struct {
	gophercloud.ErrResult
}

// MemberUpdateResult model
type MemberUpdateResult struct {
	CreateMemberResult
}
