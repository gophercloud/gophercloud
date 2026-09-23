package testing

import (
	"errors"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
)

func TestNoAuthV2Opts(t *testing.T) {
	tests := []struct {
		name string
		opts auth.NoAuthV2Opts
		want string
	}{
		{
			name: "explicit username and tenant name",
			opts: auth.NoAuthV2Opts{Username: "user", TenantName: "tenant"},
			want: "user:tenant",
		},
		{
			name: "defaults",
			want: "admin:admin",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.opts.ToNoAuthToken()
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Fatalf("token = %q, want %q", got, test.want)
			}
		})
	}
}

func TestNoAuthV2OptsRejectsDelimiter(t *testing.T) {
	for _, opts := range []auth.NoAuthV2Opts{
		{Username: "user:name", TenantName: "tenant"},
		{Username: "user", TenantName: "tenant:name"},
	} {
		_, err := opts.ToNoAuthToken()
		if err == nil {
			t.Fatal("expected an error")
		}
		var invalid gophercloud.ErrInvalidInput
		if !errors.As(err, &invalid) {
			t.Fatalf("error type = %T, want gophercloud.ErrInvalidInput", err)
		}
	}
}

func TestNoAuthV3Opts(t *testing.T) {
	tests := []struct {
		name string
		opts auth.NoAuthV3Opts
		want string
	}{
		{
			name: "explicit user and project IDs",
			opts: auth.NoAuthV3Opts{UserID: "user-id", ProjectID: "project-id"},
			want: "user-id:project-id",
		},
		{
			name: "project ID defaults to user ID",
			opts: auth.NoAuthV3Opts{UserID: "user-id"},
			want: "user-id:user-id",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.opts.ToNoAuthToken()
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Fatalf("token = %q, want %q", got, test.want)
			}
		})
	}
}

func TestNoAuthV3OptsRequiresUserID(t *testing.T) {
	_, err := (auth.NoAuthV3Opts{ProjectID: "project-id"}).ToNoAuthToken()
	if err == nil {
		t.Fatal("expected an error")
	}
	var missing gophercloud.ErrMissingInput
	if !errors.As(err, &missing) {
		t.Fatalf("error type = %T, want gophercloud.ErrMissingInput", err)
	}
	if missing.Argument != "UserID" {
		t.Fatalf("missing argument = %q, want UserID", missing.Argument)
	}
}

func TestNoAuthV3OptsRejectsDelimiter(t *testing.T) {
	for _, opts := range []auth.NoAuthV3Opts{
		{UserID: "user:id", ProjectID: "project-id"},
		{UserID: "user-id", ProjectID: "project:id"},
	} {
		_, err := opts.ToNoAuthToken()
		if err == nil {
			t.Fatal("expected an error")
		}
		var invalid gophercloud.ErrInvalidInput
		if !errors.As(err, &invalid) {
			t.Fatalf("error type = %T, want gophercloud.ErrInvalidInput", err)
		}
	}
}
