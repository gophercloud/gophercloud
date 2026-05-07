package clouds_test

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/config/clouds"
)

func ExampleWithCloudName() {
	const exampleClouds = `clouds:
  openstack:
    auth:
      auth_url: https://example.com:13000`

	ao, _, _, err := clouds.Parse(
		clouds.WithCloudsYAML(strings.NewReader(exampleClouds)),
		clouds.WithCloudName("openstack"),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(ao.IdentityEndpoint)
	// Output: https://example.com:13000
}

func ExampleWithUserID() {
	const exampleClouds = `clouds:
  openstack:
    auth:
      auth_url: https://example.com:13000`

	ao, _, _, err := clouds.Parse(
		clouds.WithCloudsYAML(strings.NewReader(exampleClouds)),
		clouds.WithCloudName("openstack"),
		clouds.WithUsername("Kris"),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(ao.Username)
	// Output: Kris
}

func ExampleWithRegion() {
	const exampleClouds = `clouds:
  openstack:
    auth:
      auth_url: https://example.com:13000`

	_, eo, _, err := clouds.Parse(
		clouds.WithCloudsYAML(strings.NewReader(exampleClouds)),
		clouds.WithCloudName("openstack"),
		clouds.WithRegion("mars"),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(eo.Region)
	// Output: mars
}

func TestParse(t *testing.T) {
	const tempDirPrefix = "gophercloud-test-"

	rmTmpDirOrPanic := func(tmpDir string) {
		if err := os.RemoveAll(tmpDir); err != nil {
			panic("unable to remove the temporary files: " + err.Error())
		}
	}

	t.Run("parses the local clouds.yaml and secure.yaml if present", func(t *testing.T) {
		const cloudsYAML = `clouds:
  gophercloud-test:
    auth:
      auth_url: https://example.com/gophercloud-test-12345:13000`
		const secureYAML = `clouds:
  gophercloud-test:
    auth:
      password: secret
      username: gophercloud-test-username`

		tmpDir, err := os.MkdirTemp(os.TempDir(), tempDirPrefix)
		if err != nil {
			t.Fatalf("unable to create a temporary directory: %v", err)
		}
		defer rmTmpDirOrPanic(tmpDir)

		cwd, err := os.Getwd()
		if err != nil {
			t.Fatalf("unable to determine the current working directory: %v", err)
		}
		if err := os.Chdir(tmpDir); err != nil {
			t.Fatalf("unable to move to a temporary directory: %v", err)
		}
		defer func() {
			if err := os.Chdir(cwd); err != nil {
				panic("unable to reset the current working directory: " + err.Error())
			}
		}()

		if err := os.WriteFile("clouds.yaml", []byte(cloudsYAML), 0644); err != nil {
			t.Fatalf("unable to create a mock clouds.yaml file: %v", err)
		}

		if err := os.WriteFile("secure.yaml", []byte(secureYAML), 0644); err != nil {
			t.Fatalf("unable to create a mock secure.yaml file: %v", err)
		}

		ao, _, _, err := clouds.Parse(
			clouds.WithCloudName("gophercloud-test"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := ao.IdentityEndpoint; got != "https://example.com/gophercloud-test-12345:13000" {
			t.Errorf("unexpected identity endpoint: %q", got)
		}

		if got := ao.Username; got != "gophercloud-test-username" {
			t.Errorf("unexpected username: %q", got)
		}
	})

	t.Run("parses the locations in order", func(t *testing.T) {
		const cloudsYAML1 = `clouds:
  gophercloud-test:
    auth:
      auth_url: https://example.com/gophercloud-test-1:13000`
		const cloudsYAML2 = `clouds:
  gophercloud-test:
    auth:
      auth_url: https://example.com/gophercloud-test-2:13000`

		tmpDir1, err := os.MkdirTemp(os.TempDir(), tempDirPrefix)
		if err != nil {
			t.Fatalf("unable to create a temporary directory: %v", err)
		}
		defer rmTmpDirOrPanic(tmpDir1)

		tmpDir2, err := os.MkdirTemp(os.TempDir(), tempDirPrefix)
		if err != nil {
			t.Fatalf("unable to create a temporary directory: %v", err)
		}
		defer rmTmpDirOrPanic(tmpDir2)

		cloudsPath1, cloudsPath2 := path.Join(tmpDir1, "clouds.yaml"), path.Join(tmpDir2, "clouds.yaml")

		if err := os.WriteFile(cloudsPath1, []byte(cloudsYAML1), 0644); err != nil {
			t.Fatalf("unable to create a mock clouds.yaml file in path %q: %v", cloudsPath1, err)
		}
		if err := os.WriteFile(cloudsPath2, []byte(cloudsYAML2), 0644); err != nil {
			t.Fatalf("unable to create a mock clouds.yaml file in path %q: %v", cloudsPath2, err)
		}

		ao, _, _, err := clouds.Parse(
			clouds.WithCloudName("gophercloud-test"),
			clouds.WithLocations(cloudsPath1, cloudsPath2),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := ao.IdentityEndpoint; got != "https://example.com/gophercloud-test-1:13000" {
			t.Errorf("unexpected identity endpoint: %q", got)
		}
	})

	t.Run("uses XDG_CONFIG_HOME for clouds.yaml location", func(t *testing.T) {
		const cloudsYAML = `clouds:
  gophercloud-test:
    auth:
      auth_url: https://example.com/xdg-config:13000`

		// Create a temp dir to use as XDG_CONFIG_HOME with clouds.yaml inside.
		xdgDir, err := os.MkdirTemp(os.TempDir(), tempDirPrefix)
		if err != nil {
			t.Fatalf("unable to create a temporary directory: %v", err)
		}
		defer rmTmpDirOrPanic(xdgDir)

		openstackDir := path.Join(xdgDir, "openstack")
		if err := os.MkdirAll(openstackDir, 0755); err != nil {
			t.Fatalf("unable to create openstack config directory: %v", err)
		}
		if err := os.WriteFile(path.Join(openstackDir, "clouds.yaml"), []byte(cloudsYAML), 0644); err != nil {
			t.Fatalf("unable to create a mock clouds.yaml file: %v", err)
		}

		// Change to an empty temp dir so cwd has no clouds.yaml.
		emptyDir, err := os.MkdirTemp(os.TempDir(), tempDirPrefix)
		if err != nil {
			t.Fatalf("unable to create a temporary directory: %v", err)
		}
		defer rmTmpDirOrPanic(emptyDir)

		cwd, err := os.Getwd()
		if err != nil {
			t.Fatalf("unable to determine the current working directory: %v", err)
		}
		if err := os.Chdir(emptyDir); err != nil {
			t.Fatalf("unable to move to a temporary directory: %v", err)
		}
		defer func() {
			if err := os.Chdir(cwd); err != nil {
				panic("unable to reset the current working directory: " + err.Error())
			}
		}()

		t.Setenv("XDG_CONFIG_HOME", xdgDir)
		t.Setenv("OS_CLIENT_CONFIG_FILE", "")

		ao, _, _, err := clouds.Parse(
			clouds.WithCloudName("gophercloud-test"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := ao.IdentityEndpoint; got != "https://example.com/xdg-config:13000" {
			t.Errorf("unexpected identity endpoint: %q", got)
		}
	})

	t.Run("falls back to ~/.config when XDG_CONFIG_HOME is not set", func(t *testing.T) {
		const cloudsYAML = `clouds:
  gophercloud-test:
    auth:
      auth_url: https://example.com/home-config:13000`

		// Create a temp dir to use as HOME with clouds.yaml in .config/openstack/.
		homeDir, err := os.MkdirTemp(os.TempDir(), tempDirPrefix)
		if err != nil {
			t.Fatalf("unable to create a temporary directory: %v", err)
		}
		defer rmTmpDirOrPanic(homeDir)

		openstackDir := path.Join(homeDir, ".config", "openstack")
		if err := os.MkdirAll(openstackDir, 0755); err != nil {
			t.Fatalf("unable to create openstack config directory: %v", err)
		}
		if err := os.WriteFile(path.Join(openstackDir, "clouds.yaml"), []byte(cloudsYAML), 0644); err != nil {
			t.Fatalf("unable to create a mock clouds.yaml file: %v", err)
		}

		// Change to an empty temp dir so cwd has no clouds.yaml.
		emptyDir, err := os.MkdirTemp(os.TempDir(), tempDirPrefix)
		if err != nil {
			t.Fatalf("unable to create a temporary directory: %v", err)
		}
		defer rmTmpDirOrPanic(emptyDir)

		cwd, err := os.Getwd()
		if err != nil {
			t.Fatalf("unable to determine the current working directory: %v", err)
		}
		if err := os.Chdir(emptyDir); err != nil {
			t.Fatalf("unable to move to a temporary directory: %v", err)
		}
		defer func() {
			if err := os.Chdir(cwd); err != nil {
				panic("unable to reset the current working directory: " + err.Error())
			}
		}()

		t.Setenv("XDG_CONFIG_HOME", "")
		t.Setenv("HOME", homeDir)
		t.Setenv("OS_CLIENT_CONFIG_FILE", "")

		ao, _, _, err := clouds.Parse(
			clouds.WithCloudName("gophercloud-test"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := ao.IdentityEndpoint; got != "https://example.com/home-config:13000" {
			t.Errorf("unexpected identity endpoint: %q", got)
		}
	})

	t.Run("falls back to the next location if clouds.yaml is not found", func(t *testing.T) {
		const cloudsYAML1 = `clouds:
  gophercloud-test:
    auth:
      auth_url: https://example.com/gophercloud-test-1:13000`
		const cloudsYAML2 = `clouds:
  gophercloud-test:
    auth:
      auth_url: https://example.com/gophercloud-test-2:13000`

		tmpDir0, err := os.MkdirTemp(os.TempDir(), tempDirPrefix)
		if err != nil {
			t.Fatalf("unable to create a temporary directory: %v", err)
		}
		defer rmTmpDirOrPanic(tmpDir0)

		tmpDir1, err := os.MkdirTemp(os.TempDir(), tempDirPrefix)
		if err != nil {
			t.Fatalf("unable to create a temporary directory: %v", err)
		}
		defer rmTmpDirOrPanic(tmpDir1)

		tmpDir2, err := os.MkdirTemp(os.TempDir(), tempDirPrefix)
		if err != nil {
			t.Fatalf("unable to create a temporary directory: %v", err)
		}
		defer rmTmpDirOrPanic(tmpDir2)

		cloudsPath0, cloudsPath1, cloudsPath2 := path.Join(tmpDir0, "clouds.yaml"), path.Join(tmpDir1, "clouds.yaml"), path.Join(tmpDir2, "clouds.yaml")

		if err := os.WriteFile(cloudsPath1, []byte(cloudsYAML1), 0644); err != nil {
			t.Fatalf("unable to create a mock clouds.yaml file in path %q: %v", cloudsPath1, err)
		}
		if err := os.WriteFile(cloudsPath2, []byte(cloudsYAML2), 0644); err != nil {
			t.Fatalf("unable to create a mock clouds.yaml file in path %q: %v", cloudsPath2, err)
		}

		ao, _, _, err := clouds.Parse(
			clouds.WithCloudName("gophercloud-test"),
			clouds.WithLocations(cloudsPath0, cloudsPath1, cloudsPath2),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := ao.IdentityEndpoint; got != "https://example.com/gophercloud-test-1:13000" {
			t.Errorf("unexpected identity endpoint: %q", got)
		}
	})

	t.Run("supports user in one domain and project in another domain using names", func(t *testing.T) {
		const cloudsYAML = `clouds:
  gophercloud-test:
    auth:
      auth_url: https://example.com:5000/v3
      username: myuser
      password: mypassword
      user_domain_name: Default
      project_name: myproject
      project_domain_name: some_domain`

		ao, _, _, err := clouds.Parse(
			clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
			clouds.WithCloudName("gophercloud-test"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := ao.Username; got != "myuser" {
			t.Errorf("unexpected username: %q", got)
		}
		if got := ao.DomainName; got != "Default" {
			t.Errorf("unexpected user domain name: %q, expected 'Default'", got)
		}
		if got := ao.TenantName; got != "myproject" {
			t.Errorf("unexpected project name: %q", got)
		}
		if ao.Scope == nil {
			t.Fatal("expected Scope to be set")
		}
		if got := ao.Scope.ProjectName; got != "myproject" {
			t.Errorf("unexpected scope project name: %q", got)
		}
		if got := ao.Scope.DomainName; got != "some_domain" {
			t.Errorf("unexpected scope domain name: %q, expected 'some_domain'", got)
		}
	})

	t.Run("supports user in one domain and project in another domain using IDs", func(t *testing.T) {
		const cloudsYAML = `clouds:
  gophercloud-test:
    auth:
      auth_url: https://example.com:5000/v3
      username: myuser
      password: mypassword
      user_domain_id: default-domain-id
      project_id: project-123
      project_domain_id: other-domain-id`

		ao, _, _, err := clouds.Parse(
			clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
			clouds.WithCloudName("gophercloud-test"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := ao.Username; got != "myuser" {
			t.Errorf("unexpected username: %q", got)
		}
		if got := ao.DomainID; got != "default-domain-id" {
			t.Errorf("unexpected user domain ID: %q, expected 'default-domain-id'", got)
		}
		if got := ao.TenantID; got != "project-123" {
			t.Errorf("unexpected project ID: %q", got)
		}
		if ao.Scope == nil {
			t.Fatal("expected Scope to be set")
		}
		if got := ao.Scope.ProjectID; got != "project-123" {
			t.Errorf("unexpected scope project ID: %q", got)
		}
		// When using project_id, the domain is not needed in scope
		if ao.Scope.DomainID != "" {
			t.Errorf("expected scope domain ID to be empty when using project_id, got: %q", ao.Scope.DomainID)
		}
	})

	t.Run("falls back to domain_name for both user and project when specific domains not set", func(t *testing.T) {
		const cloudsYAML = `clouds:
  gophercloud-test:
    auth:
      auth_url: https://example.com:5000/v3
      username: myuser
      password: mypassword
      domain_name: shared-domain
      project_name: myproject`

		ao, _, _, err := clouds.Parse(
			clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
			clouds.WithCloudName("gophercloud-test"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := ao.Username; got != "myuser" {
			t.Errorf("unexpected username: %q", got)
		}
		if got := ao.DomainName; got != "shared-domain" {
			t.Errorf("unexpected user domain name: %q, expected 'shared-domain'", got)
		}
		if ao.Scope == nil {
			t.Fatal("expected Scope to be set")
		}
		if got := ao.Scope.ProjectName; got != "myproject" {
			t.Errorf("unexpected scope project name: %q", got)
		}
		if got := ao.Scope.DomainName; got != "shared-domain" {
			t.Errorf("unexpected scope domain name: %q, expected 'shared-domain'", got)
		}
	})

	t.Run("user_domain_name takes precedence over domain_name for user identity", func(t *testing.T) {
		const cloudsYAML = `clouds:
  gophercloud-test:
    auth:
      auth_url: https://example.com:5000/v3
      username: myuser
      password: mypassword
      user_domain_name: user-specific-domain
      domain_name: fallback-domain
      project_name: myproject`

		ao, _, _, err := clouds.Parse(
			clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
			clouds.WithCloudName("gophercloud-test"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := ao.DomainName; got != "user-specific-domain" {
			t.Errorf("unexpected user domain name: %q, expected 'user-specific-domain'", got)
		}
		if ao.Scope == nil {
			t.Fatal("expected Scope to be set")
		}
		if got := ao.Scope.DomainName; got != "fallback-domain" {
			t.Errorf("unexpected scope domain name: %q, expected 'fallback-domain'", got)
		}
	})

	t.Run("project_domain_name takes precedence over domain_name for project scope", func(t *testing.T) {
		const cloudsYAML = `clouds:
  gophercloud-test:
    auth:
      auth_url: https://example.com:5000/v3
      username: myuser
      password: mypassword
      domain_name: fallback-domain
      project_name: myproject
      project_domain_name: project-specific-domain`

		ao, _, _, err := clouds.Parse(
			clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
			clouds.WithCloudName("gophercloud-test"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := ao.DomainName; got != "fallback-domain" {
			t.Errorf("unexpected user domain name: %q, expected 'fallback-domain'", got)
		}
		if ao.Scope == nil {
			t.Fatal("expected Scope to be set")
		}
		if got := ao.Scope.DomainName; got != "project-specific-domain" {
			t.Errorf("unexpected scope domain name: %q, expected 'project-specific-domain'", got)
		}
	})

	t.Run("project_id scoping does not require domain information", func(t *testing.T) {
		const cloudsYAML = `clouds:
  gophercloud-test:
    auth:
      auth_url: https://example.com:5000/v3
      username: myuser
      password: mypassword
      user_domain_name: Default
      project_id: unique-project-id-123`

		ao, _, _, err := clouds.Parse(
			clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
			clouds.WithCloudName("gophercloud-test"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if ao.Scope == nil {
			t.Fatal("expected Scope to be set")
		}
		if got := ao.Scope.ProjectID; got != "unique-project-id-123" {
			t.Errorf("unexpected scope project ID: %q", got)
		}
		if ao.Scope.DomainID != "" || ao.Scope.DomainName != "" {
			t.Errorf("expected no domain information in scope when using project_id, got DomainID=%q, DomainName=%q",
				ao.Scope.DomainID, ao.Scope.DomainName)
		}
	})

	t.Run("supports system_scope: all for system-level operations", func(t *testing.T) {
		const cloudsYAML = `clouds:
  gophercloud-test:
    auth:
      auth_url: https://example.com:5000/v3
      username: admin
      password: adminpassword
      user_domain_name: Default
      system_scope: all`

		ao, _, _, err := clouds.Parse(
			clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
			clouds.WithCloudName("gophercloud-test"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := ao.Username; got != "admin" {
			t.Errorf("unexpected username: %q", got)
		}
		if got := ao.DomainName; got != "Default" {
			t.Errorf("unexpected user domain name: %q", got)
		}
		if ao.Scope == nil {
			t.Fatal("expected Scope to be set")
		}
		if !ao.Scope.System {
			t.Error("expected Scope.System to be true")
		}
		if ao.Scope.ProjectID != "" || ao.Scope.ProjectName != "" {
			t.Errorf("expected no project information in system scope, got ProjectID=%q, ProjectName=%q",
				ao.Scope.ProjectID, ao.Scope.ProjectName)
		}
	})

	t.Run("system_scope takes precedence over project scope", func(t *testing.T) {
		const cloudsYAML = `clouds:
  gophercloud-test:
    auth:
      auth_url: https://example.com:5000/v3
      username: admin
      password: adminpassword
      user_domain_name: Default
      system_scope: all
      project_name: should-be-ignored
      project_domain_name: also-ignored`

		ao, _, _, err := clouds.Parse(
			clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
			clouds.WithCloudName("gophercloud-test"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if ao.Scope == nil {
			t.Fatal("expected Scope to be set")
		}
		if !ao.Scope.System {
			t.Error("expected Scope.System to be true")
		}
		if ao.Scope.ProjectName != "" {
			t.Errorf("expected project_name to be ignored when system_scope is set, got: %q", ao.Scope.ProjectName)
		}
	})

	t.Run("ignores system_scope when value is not all", func(t *testing.T) {
		const cloudsYAML = `clouds:
  gophercloud-test:
    auth:
      auth_url: https://example.com:5000/v3
      username: myuser
      password: mypassword
      user_domain_name: Default
      system_scope: something-else
      project_name: myproject
      project_domain_name: Default`

		ao, _, _, err := clouds.Parse(
			clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
			clouds.WithCloudName("gophercloud-test"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if ao.Scope == nil {
			t.Fatal("expected Scope to be set")
		}
		if ao.Scope.System {
			t.Error("expected Scope.System to be false when system_scope is not all")
		}
		if got := ao.Scope.ProjectName; got != "myproject" {
			t.Errorf("expected project scoping to be used, got ProjectName=%q", got)
		}
	})

	t.Run("supports domain scoping by domain_id when no project specified", func(t *testing.T) {
		const cloudsYAML = `clouds:
  gophercloud-test:
    auth:
      auth_url: https://example.com:5000/v3
      username: domainadmin
      password: mypassword
      user_domain_name: Default
      domain_id: domain-123`

		ao, _, _, err := clouds.Parse(
			clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
			clouds.WithCloudName("gophercloud-test"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := ao.Username; got != "domainadmin" {
			t.Errorf("unexpected username: %q", got)
		}
		if got := ao.DomainName; got != "Default" {
			t.Errorf("unexpected user domain name: %q", got)
		}
		if ao.Scope == nil {
			t.Fatal("expected Scope to be set")
		}
		if got := ao.Scope.DomainID; got != "domain-123" {
			t.Errorf("unexpected scope domain ID: %q, expected 'domain-123'", got)
		}
		if ao.Scope.ProjectID != "" || ao.Scope.ProjectName != "" {
			t.Errorf("expected no project in domain scope, got ProjectID=%q, ProjectName=%q",
				ao.Scope.ProjectID, ao.Scope.ProjectName)
		}
	})

	t.Run("supports domain scoping by domain_name when no project specified", func(t *testing.T) {
		const cloudsYAML = `clouds:
  gophercloud-test:
    auth:
      auth_url: https://example.com:5000/v3
      username: domainadmin
      password: mypassword
      user_domain_name: UserDomain
      domain_name: ScopeDomain`

		ao, _, _, err := clouds.Parse(
			clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
			clouds.WithCloudName("gophercloud-test"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := ao.Username; got != "domainadmin" {
			t.Errorf("unexpected username: %q", got)
		}
		if got := ao.DomainName; got != "UserDomain" {
			t.Errorf("unexpected user domain name: %q, expected 'UserDomain'", got)
		}
		if ao.Scope == nil {
			t.Fatal("expected Scope to be set")
		}
		if got := ao.Scope.DomainName; got != "ScopeDomain" {
			t.Errorf("unexpected scope domain name: %q, expected 'ScopeDomain'", got)
		}
		if ao.Scope.ProjectID != "" || ao.Scope.ProjectName != "" {
			t.Errorf("expected no project in domain scope, got ProjectID=%q, ProjectName=%q",
				ao.Scope.ProjectID, ao.Scope.ProjectName)
		}
	})

	t.Run("project scope takes precedence over domain scope", func(t *testing.T) {
		const cloudsYAML = `clouds:
  gophercloud-test:
    auth:
      auth_url: https://example.com:5000/v3
      username: myuser
      password: mypassword
      user_domain_name: Default
      project_name: myproject
      project_domain_name: ProjectDomain
      domain_name: should-be-ignored`

		ao, _, _, err := clouds.Parse(
			clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
			clouds.WithCloudName("gophercloud-test"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if ao.Scope == nil {
			t.Fatal("expected Scope to be set")
		}
		if got := ao.Scope.ProjectName; got != "myproject" {
			t.Errorf("expected project scope, got ProjectName=%q", got)
		}
		if got := ao.Scope.DomainName; got != "ProjectDomain" {
			t.Errorf("expected project domain in scope, got DomainName=%q", got)
		}
	})

	t.Run("domain scoping with user_domain_id and domain_id", func(t *testing.T) {
		const cloudsYAML = `clouds:
  gophercloud-test:
    auth:
      auth_url: https://example.com:5000/v3
      username: domainadmin
      password: mypassword
      user_domain_id: user-domain-id-123
      domain_id: scope-domain-id-456`

		ao, _, _, err := clouds.Parse(
			clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
			clouds.WithCloudName("gophercloud-test"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := ao.DomainID; got != "user-domain-id-123" {
			t.Errorf("unexpected user domain ID: %q, expected 'user-domain-id-123'", got)
		}
		if ao.Scope == nil {
			t.Fatal("expected Scope to be set")
		}
		if got := ao.Scope.DomainID; got != "scope-domain-id-456" {
			t.Errorf("unexpected scope domain ID: %q, expected 'scope-domain-id-456'", got)
		}
	})
}
