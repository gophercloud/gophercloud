package auth_test

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/v2/auth"
)

func ExampleAuthOptionsV3_Authenticate_password() {
	opts := auth.AuthOptionsV3{
		AuthURL: "https://identity.example.com/v3",
		Auth: auth.V3PasswordOpts{
			Username:       "username",
			Password:       "password",
			UserDomainName: "Default",
			Scope:          &auth.Scope{ProjectID: "project-id"},
		},
	}
	result, err := opts.Authenticate(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.TokenID)

	_ = result
}

func ExampleAuthOptionsV2_Authenticate() {
	opts := auth.AuthOptionsV2{
		AuthURL: "https://identity.example.com/v2.0",
		Auth: auth.V2PasswordOpts{
			Username: "username",
			Password: "password",
			TenantID: "tenant-id",
		},
	}
	result, err := opts.Authenticate(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	_ = result
}

func ExampleAuthOptionsV3_Authenticate_token() {
	opts := auth.AuthOptionsV3{
		AuthURL: "https://identity.example.com/v3",
		Auth: auth.V3TokenOpts{
			Token: "existing-token",
			Scope: &auth.Scope{ProjectID: "project-id"},
		},
	}
	result, err := opts.Authenticate(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	_ = result
}

func ExampleAuthOptionsV3_Authenticate_applicationCredential() {
	credentials := auth.V3ApplicationCredentialOpts{
		ApplicationCredentialID:     "credential-id",
		ApplicationCredentialSecret: "credential-secret",
	}
	opts := auth.AuthOptionsV3{
		AuthURL: "https://identity.example.com/v3",
		Auth:    credentials,
	}
	result, err := opts.Authenticate(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	_ = result
}

func ExampleAuthOptionsV3_Authenticate_multifactor() {
	opts := auth.AuthOptionsV3{
		AuthURL: "https://identity.example.com/v3",
		Auth: auth.V3MultifactorOpts{
			AuthMethods: []auth.AuthOptionsBuilderV3{
				auth.V3PasswordOpts{UserID: "user-id", Password: "password"},
				auth.V3TOTPOpts{UserID: "user-id", Passcode: "123456"},
			},
			Scope: &auth.Scope{ProjectID: "project-id"},
		},
	}
	result, err := opts.Authenticate(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	_ = result
}
