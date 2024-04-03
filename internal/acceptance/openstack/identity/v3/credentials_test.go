//go:build acceptance || identity || credentials

package v3

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/credentials"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/ec2tokens"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestCredentialsCRUD(t *testing.T) {
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
	token, err := tokens.Create(context.TODO(), client, &authOptions).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, token)

	user, err := tokens.Get(context.TODO(), client, token.ID).ExtractUser()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, user)

	project, err := tokens.Get(context.TODO(), client, token.ID).ExtractProject()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, project)

	createOpts := credentials.CreateOpts{
		ProjectID: project.ID,
		Type:      "ec2",
		UserID:    user.ID,
		Blob:      "{\"access\":\"181920\",\"secret\":\"secretKey\"}",
	}

	// Create a credential
	credential, err := credentials.Create(context.TODO(), client, createOpts).Extract()
	th.AssertNoErr(t, err)

	// Delete a credential
	defer credentials.Delete(context.TODO(), client, credential.ID)
	tools.PrintResource(t, credential)

	th.AssertEquals(t, credential.Blob, createOpts.Blob)
	th.AssertEquals(t, credential.Type, createOpts.Type)
	th.AssertEquals(t, credential.UserID, createOpts.UserID)
	th.AssertEquals(t, credential.ProjectID, createOpts.ProjectID)

	// Get a credential
	getCredential, err := credentials.Get(context.TODO(), client, credential.ID).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, getCredential)

	th.AssertEquals(t, getCredential.Blob, createOpts.Blob)
	th.AssertEquals(t, getCredential.Type, createOpts.Type)
	th.AssertEquals(t, getCredential.UserID, createOpts.UserID)
	th.AssertEquals(t, getCredential.ProjectID, createOpts.ProjectID)

	updateOpts := credentials.UpdateOpts{
		ProjectID: project.ID,
		Type:      "ec2",
		UserID:    user.ID,
		Blob:      "{\"access\":\"181920\",\"secret\":\"mySecret\"}",
	}

	// Update a credential
	updateCredential, err := credentials.Update(context.TODO(), client, credential.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, updateCredential)

	th.AssertEquals(t, updateCredential.Blob, updateOpts.Blob)
}

func TestCredentialsValidateS3(t *testing.T) {
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
	token, err := tokens.Create(context.TODO(), client, &authOptions).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, token)

	user, err := tokens.Get(context.TODO(), client, token.ID).ExtractUser()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, user)

	project, err := tokens.Get(context.TODO(), client, token.ID).ExtractProject()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, project)

	createOpts := credentials.CreateOpts{
		ProjectID: project.ID,
		Type:      "ec2",
		UserID:    user.ID,
		Blob:      "{\"access\":\"181920\",\"secret\":\"secretKey\"}",
	}

	// Create a credential
	credential, err := credentials.Create(context.TODO(), client, createOpts).Extract()
	th.AssertNoErr(t, err)

	// Delete a credential
	defer credentials.Delete(context.TODO(), client, credential.ID)
	tools.PrintResource(t, credential)

	th.AssertEquals(t, credential.Blob, createOpts.Blob)
	th.AssertEquals(t, credential.Type, createOpts.Type)
	th.AssertEquals(t, credential.UserID, createOpts.UserID)
	th.AssertEquals(t, credential.ProjectID, createOpts.ProjectID)

	opts := ec2tokens.AuthOptions{
		Access: "181920",
		Secret: "secretKey",
		// auth will fail if this is not s3
		Service: "s3",
	}

	// Validate a credential
	token, err = ec2tokens.ValidateS3Token(context.TODO(), client, &opts).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, token)
}
