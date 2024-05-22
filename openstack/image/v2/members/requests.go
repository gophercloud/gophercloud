package members

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

/*
Create member for specific image

# Preconditions

  - The specified images must exist.
  - You can only add a new member to an image which 'visibility' attribute is
    private.
  - You must be the owner of the specified image.

# Synchronous Postconditions

With correct permissions, you can see the member status of the image as
pending through API calls.

More details here:
http://developer.openstack.org/api-ref-image-v2.html#createImageMember-v2
*/
func Create(ctx context.Context, client *gophercloud.ServiceClient, id string, member string) (r CreateResult) {
	b := map[string]any{"member": member}
	resp, err := client.Post(ctx, createMemberURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// List members returns list of members for specifed image id.
func List(client *gophercloud.ServiceClient, id string) pagination.Pager {
	return pagination.NewPager(client, listMembersURL(client, id), func(r pagination.PageResult) pagination.Page {
		return MemberPage{pagination.SinglePageBase(r)}
	})
}

// Get image member details.
func Get(ctx context.Context, client *gophercloud.ServiceClient, imageID string, memberID string) (r DetailsResult) {
	resp, err := client.Get(ctx, getMemberURL(client, imageID, memberID), &r.Body, &gophercloud.RequestOpts{OkCodes: []int{200}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete membership for given image. Callee should be image owner.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, imageID string, memberID string) (r DeleteResult) {
	resp, err := client.Delete(ctx, deleteMemberURL(client, imageID, memberID), &gophercloud.RequestOpts{OkCodes: []int{204}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional attributes to the
// Update request.
type UpdateOptsBuilder interface {
	ToImageMemberUpdateMap() (map[string]any, error)
}

// UpdateOpts represents options to an Update request.
type UpdateOpts struct {
	Status string
}

// ToMemberUpdateMap formats an UpdateOpts structure into a request body.
func (opts UpdateOpts) ToImageMemberUpdateMap() (map[string]any, error) {
	return map[string]any{
		"status": opts.Status,
	}, nil
}

// Update function updates member.
func Update(ctx context.Context, client *gophercloud.ServiceClient, imageID string, memberID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToImageMemberUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(ctx, updateMemberURL(client, imageID, memberID), b, &r.Body,
		&gophercloud.RequestOpts{OkCodes: []int{200}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
