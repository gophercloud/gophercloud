//go:build acceptance || identity || credentials

package v3

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/credentials"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/ec2tokens"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// assertBlobJSONEquals compares a credential's serialized Blob against an
// expected JSON string.
func assertBlobJSONEquals(t *testing.T, expected, actual string) {
	t.Helper()

	var v any
	err := json.Unmarshal([]byte(actual), &v)
	th.AssertNoErr(t, err)

	th.AssertJSONEquals(t, expected, v)
}

func TestCredentialsCRUD(t *testing.T) {
	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	ao, err := auth.AuthOptionsFromEnvV3()
	th.AssertNoErr(t, err)

	result, err := ao.Authenticate(context.TODO(), client.ProviderClient)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, result)

	user := &result.User
	tools.PrintResource(t, user)

	project := result.Project
	if project == nil {
		t.Fatal("authentication did not return a project")
	}
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

	assertBlobJSONEquals(t, createOpts.Blob, credential.Blob)
	th.AssertEquals(t, credential.Type, createOpts.Type)
	th.AssertEquals(t, credential.UserID, createOpts.UserID)
	th.AssertEquals(t, credential.ProjectID, createOpts.ProjectID)

	// Get a credential
	getCredential, err := credentials.Get(context.TODO(), client, credential.ID).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, getCredential)

	assertBlobJSONEquals(t, createOpts.Blob, getCredential.Blob)
	th.AssertEquals(t, getCredential.Type, createOpts.Type)
	th.AssertEquals(t, getCredential.UserID, createOpts.UserID)
	th.AssertEquals(t, getCredential.ProjectID, createOpts.ProjectID)

	updateOpts := credentials.UpdateOpts{
		Blob: "{\"access\":\"181920\",\"secret\":\"mySecret\"}",
	}

	// Update a credential
	updateCredential, err := credentials.Update(context.TODO(), client, credential.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, updateCredential)

	assertBlobJSONEquals(t, updateOpts.Blob, updateCredential.Blob)
	th.AssertEquals(t, createOpts.Type, updateCredential.Type)
	th.AssertEquals(t, createOpts.UserID, updateCredential.UserID)
	th.AssertEquals(t, createOpts.ProjectID, updateCredential.ProjectID)
}

func TestCredentialsValidateS3(t *testing.T) {
	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	ao, err := auth.AuthOptionsFromEnvV3()
	th.AssertNoErr(t, err)

	result, err := ao.Authenticate(context.TODO(), client.ProviderClient)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, result)

	user := &result.User
	tools.PrintResource(t, user)

	project := result.Project
	if project == nil {
		t.Fatal("authentication did not return a project")
	}
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

	assertBlobJSONEquals(t, createOpts.Blob, credential.Blob)
	th.AssertEquals(t, credential.Type, createOpts.Type)
	th.AssertEquals(t, credential.UserID, createOpts.UserID)
	th.AssertEquals(t, credential.ProjectID, createOpts.ProjectID)

	opts := auth.EC2TokenOpts{
		Access: "181920",
		Secret: "secretKey",
		// auth will fail if this is not s3
		Service: "s3",
	}

	// Validate a credential
	credentialsBody, err := opts.ToAuthBody()
	th.AssertNoErr(t, err)
	credentials := credentialsBody["credentials"]
	for _, field := range []string{"body_hash", "headers", "host", "params", "path", "verb"} {
		delete(credentials, field)
	}
	token, err := ec2tokens.ValidateS3Token(context.TODO(), client, gophercloud.RawRequest{
		Body: map[string]any{"credentials": credentials},
	}).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, token)
}
