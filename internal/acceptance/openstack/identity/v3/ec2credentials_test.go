//go:build acceptance || identity || ec2credentials

package v3

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/ec2credentials"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestEC2CredentialsCRD(t *testing.T) {
	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	ao, err := openstack.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	authOptions := tokens.AuthOptions{
		Username:   ao.Username,
		Password:   ao.Password,
		DomainName: ao.DomainName,
		DomainID:   ao.DomainID,
		// We need a scope to get the token roles list
		Scope: tokens.Scope{
			ProjectID:   ao.TenantID,
			ProjectName: ao.TenantName,
			DomainID:    ao.DomainID,
			DomainName:  ao.DomainName,
		},
	}

	res := tokens.Create(context.TODO(), client, &authOptions)
	th.AssertNoErr(t, res.Err)
	token, err := res.Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, token)

	user, err := res.ExtractUser()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, user)

	project, err := res.ExtractProject()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, project)

	createOpts := ec2credentials.CreateOpts{
		TenantID: project.ID,
	}

	ec2credential, err := ec2credentials.Create(context.TODO(), client, user.ID, createOpts).Extract()
	th.AssertNoErr(t, err)
	defer ec2credentials.Delete(context.TODO(), client, user.ID, ec2credential.Access)
	tools.PrintResource(t, ec2credential)

	access := ec2credential.Access
	secret := ec2credential.Secret
	if access == "" {
		t.Fatalf("EC2 credential access was not generated")
	}

	if secret == "" {
		t.Fatalf("EC2 credential secret was not generated")
	}

	th.AssertEquals(t, ec2credential.UserID, user.ID)
	th.AssertEquals(t, ec2credential.TenantID, project.ID)

	// Get an ec2 credential
	getEC2Credential, err := ec2credentials.Get(context.TODO(), client, user.ID, ec2credential.Access).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, getEC2Credential)

	th.AssertEquals(t, getEC2Credential.UserID, user.ID)
	th.AssertEquals(t, getEC2Credential.TenantID, project.ID)
	th.AssertEquals(t, getEC2Credential.Access, access)
	th.AssertEquals(t, getEC2Credential.Secret, secret)

	allPages, err := ec2credentials.List(client, user.ID).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	credentials, err := ec2credentials.ExtractCredentials(allPages)
	th.AssertNoErr(t, err)

	if v := len(credentials); v != 1 {
		t.Fatalf("expected to list one credential, got %d", v)
	}

	th.AssertEquals(t, credentials[0].UserID, user.ID)
	th.AssertEquals(t, credentials[0].TenantID, project.ID)
	th.AssertEquals(t, credentials[0].Access, access)
	th.AssertEquals(t, credentials[0].Secret, secret)
}
