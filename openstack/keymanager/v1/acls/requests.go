package acls

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
)

// GetContainerACL retrieves the ACL of a container.
func GetContainerACL(ctx context.Context, client *gophercloud.ServiceClient, containerID string) (r ACLResult) {
	resp, err := client.Get(ctx, containerURL(client, containerID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetSecretACL retrieves the ACL of a secret.
func GetSecretACL(ctx context.Context, client *gophercloud.ServiceClient, secretID string) (r ACLResult) {
	resp, err := client.Get(ctx, secretURL(client, secretID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// SetOptsBuilder allows extensions to add additional parameters to the
// Set request.
type SetOptsBuilder interface {
	ToACLSetMap() (map[string]any, error)
}

// SetOpt represents options to set a particular ACL type on a resource.
type SetOpt struct {
	// Type is the type of ACL to set. ie: read.
	Type string `json:"-" required:"true"`

	// Users are the list of Keystone user UUIDs.
	Users *[]string `json:"users,omitempty"`

	// ProjectAccess toggles if all users in a project can access the resource.
	ProjectAccess *bool `json:"project-access,omitempty"`
}

// SetOpts represents options to set an ACL on a resource.
type SetOpts []SetOpt

// ToACLSetMap formats a SetOpts into a set request.
func (opts SetOpts) ToACLSetMap() (map[string]any, error) {
	b := make(map[string]any)
	for _, v := range opts {
		m, err := gophercloud.BuildRequestBody(v, v.Type)
		if err != nil {
			return nil, err
		}
		b[v.Type] = m[v.Type]
	}
	return b, nil
}

// SetContainerACL will set an ACL on a container.
func SetContainerACL(ctx context.Context, client *gophercloud.ServiceClient, containerID string, opts SetOptsBuilder) (r ACLRefResult) {
	b, err := opts.ToACLSetMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Put(ctx, containerURL(client, containerID), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// SetSecretACL will set an ACL on a secret.
func SetSecretACL(ctx context.Context, client *gophercloud.ServiceClient, secretID string, opts SetOptsBuilder) (r ACLRefResult) {
	b, err := opts.ToACLSetMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Put(ctx, secretURL(client, secretID), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateContainerACL will update an ACL on a container.
func UpdateContainerACL(ctx context.Context, client *gophercloud.ServiceClient, containerID string, opts SetOptsBuilder) (r ACLRefResult) {
	b, err := opts.ToACLSetMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Patch(ctx, containerURL(client, containerID), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateSecretACL will update an ACL on a secret.
func UpdateSecretACL(ctx context.Context, client *gophercloud.ServiceClient, secretID string, opts SetOptsBuilder) (r ACLRefResult) {
	b, err := opts.ToACLSetMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Patch(ctx, secretURL(client, secretID), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteContainerACL will delete an ACL from a conatiner.
func DeleteContainerACL(ctx context.Context, client *gophercloud.ServiceClient, containerID string) (r DeleteResult) {
	resp, err := client.Delete(ctx, containerURL(client, containerID), &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteSecretACL will delete an ACL from a secret.
func DeleteSecretACL(ctx context.Context, client *gophercloud.ServiceClient, secretID string) (r DeleteResult) {
	resp, err := client.Delete(ctx, secretURL(client, secretID), &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
