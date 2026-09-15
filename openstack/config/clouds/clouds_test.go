package clouds_test

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/config/clouds"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func authField(cloud clouds.Cloud, key string) string {
	v, _ := cloud.Auth[key].(string)
	return v
}

func ExampleWithCloudName() {
	const exampleClouds = `clouds:
  openstack:
    auth:
      auth_url: https://example.com:13000`

	cloud, _, _, err := clouds.Parse(
		clouds.WithCloudsYAML(strings.NewReader(exampleClouds)),
		clouds.WithCloudName("openstack"),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(authField(cloud, "auth_url"))
	// Output: https://example.com:13000
}

func ExampleWithUserID() {
	const exampleClouds = `clouds:
  openstack:
    auth:
      auth_url: https://example.com:13000
      username: Kris`

	cloud, _, _, err := clouds.Parse(
		clouds.WithCloudsYAML(strings.NewReader(exampleClouds)),
		clouds.WithCloudName("openstack"),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(authField(cloud, "username"))
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

		cloud, _, _, err := clouds.Parse(
			clouds.WithCloudName("gophercloud-test"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := authField(cloud, "auth_url"); got != "https://example.com/gophercloud-test-12345:13000" {
			t.Errorf("unexpected identity endpoint: %q", got)
		}

		if got := authField(cloud, "username"); got != "gophercloud-test-username" {
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

		cloud, _, _, err := clouds.Parse(
			clouds.WithCloudName("gophercloud-test"),
			clouds.WithLocations(cloudsPath1, cloudsPath2),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := authField(cloud, "auth_url"); got != "https://example.com/gophercloud-test-1:13000" {
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

		cloud, _, _, err := clouds.Parse(
			clouds.WithCloudName("gophercloud-test"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := authField(cloud, "auth_url"); got != "https://example.com/xdg-config:13000" {
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

		cloud, _, _, err := clouds.Parse(
			clouds.WithCloudName("gophercloud-test"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := authField(cloud, "auth_url"); got != "https://example.com/home-config:13000" {
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

		cloud, _, _, err := clouds.Parse(
			clouds.WithCloudName("gophercloud-test"),
			clouds.WithLocations(cloudsPath0, cloudsPath1, cloudsPath2),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := authField(cloud, "auth_url"); got != "https://example.com/gophercloud-test-1:13000" {
			t.Errorf("unexpected identity endpoint: %q", got)
		}
	})

	t.Run("parses clouds-public.yaml if present", func(t *testing.T) {
		const cloudsYAML = `clouds:
  gophercloud-test-0:
    cloud: gophercloud-test-1
    auth:
      user_domain_name: CustomDomain`
		const cloudsPublicYAML = `public-clouds:
  gophercloud-test-1:
    auth:
      auth_url: https://example.com/gophercloud-test-1:13000
      user_domain_name: Default`

		tmpDir, err := os.MkdirTemp(os.TempDir(), tempDirPrefix)
		th.AssertNoErr(t, err)
		defer rmTmpDirOrPanic(tmpDir)

		cwd, err := os.Getwd()
		th.AssertNoErr(t, err)
		err = os.Chdir(tmpDir)
		th.AssertNoErr(t, err)
		defer func() {
			err = os.Chdir(cwd)
			th.AssertNoErr(t, err)
		}()

		err = os.WriteFile("clouds-public.yaml", []byte(cloudsPublicYAML), 0644)
		th.AssertNoErr(t, err)

		cloud, _, _, err := clouds.Parse(
			clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
			clouds.WithCloudName("gophercloud-test-0"),
		)
		th.AssertNoErr(t, err)

		if got := authField(cloud, "auth_url"); got != "https://example.com/gophercloud-test-1:13000" {
			t.Errorf("unexpected identity endpoint: %q", got)
		}

		if got := authField(cloud, "user_domain_name"); got != "CustomDomain" {
			t.Errorf("unexpected UserDomainName: %q", got)
		}
	})
	t.Run("parses clouds-public.yaml locations in order", func(t *testing.T) {
		const cloudsYAML = `clouds:
  gophercloud-test-0:
    cloud: gophercloud-test-1
    auth:
      user_domain_name: CustomDomain`
		const cloudsPublicYAML1 = `public-clouds:
  gophercloud-test-1:
    auth:
      auth_url: https://example.com/gophercloud-test-1:13000
      user_domain_name: Default`
		const cloudsPublicYAML2 = `public-clouds:
  gophercloud-test-1:
    auth:
      auth_url: https://example.com/gophercloud-test-2:13000
      user_domain_name: Default`

		tmpDir1, err := os.MkdirTemp(os.TempDir(), tempDirPrefix)
		th.AssertNoErr(t, err)
		defer rmTmpDirOrPanic(tmpDir1)

		tmpDir2, err := os.MkdirTemp(os.TempDir(), tempDirPrefix)
		th.AssertNoErr(t, err)
		defer rmTmpDirOrPanic(tmpDir2)

		publicPath1 := path.Join(tmpDir1, "clouds-public.yaml")
		publicPath2 := path.Join(tmpDir2, "clouds-public.yaml")

		err = os.WriteFile(publicPath1, []byte(cloudsPublicYAML1), 0644)
		th.AssertNoErr(t, err)
		err = os.WriteFile(publicPath2, []byte(cloudsPublicYAML2), 0644)
		th.AssertNoErr(t, err)

		cloud, _, _, err := clouds.Parse(
			clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
			clouds.WithCloudName("gophercloud-test-0"),
			clouds.WithPublicLocations(publicPath1, publicPath2),
		)
		th.AssertNoErr(t, err)

		if got := authField(cloud, "auth_url"); got != "https://example.com/gophercloud-test-1:13000" {
			t.Errorf("unexpected identity endpoint: %q", got)
		}

		if got := authField(cloud, "user_domain_name"); got != "CustomDomain" {
			t.Errorf("unexpected UserDomainName: %q", got)
		}
	})
	t.Run("fall back to next location if clouds-public.yaml is not found", func(t *testing.T) {
		const cloudsYAML = `clouds:
  gophercloud-test-0:
    cloud: gophercloud-test-1
    auth:
      user_domain_name: CustomDomain`
		const cloudsPublicYAML1 = `public-clouds:
  gophercloud-test-1:
    auth:
      auth_url: https://example.com/gophercloud-test-1:13000
      user_domain_name: Default`
		const cloudsPublicYAML2 = `public-clouds:
  gophercloud-test-1:
    auth:
      auth_url: https://example.com/gophercloud-test-2:13000
      user_domain_name: Default`

		tmpDir0, err := os.MkdirTemp(os.TempDir(), tempDirPrefix)
		th.AssertNoErr(t, err)
		defer rmTmpDirOrPanic(tmpDir0)

		tmpDir1, err := os.MkdirTemp(os.TempDir(), tempDirPrefix)
		th.AssertNoErr(t, err)
		defer rmTmpDirOrPanic(tmpDir1)

		tmpDir2, err := os.MkdirTemp(os.TempDir(), tempDirPrefix)
		th.AssertNoErr(t, err)
		defer rmTmpDirOrPanic(tmpDir2)

		publicPath0 := path.Join(tmpDir0, "clouds-public.yaml")
		publicPath1 := path.Join(tmpDir1, "clouds-public.yaml")
		publicPath2 := path.Join(tmpDir2, "clouds-public.yaml")

		err = os.WriteFile(publicPath1, []byte(cloudsPublicYAML1), 0644)
		th.AssertNoErr(t, err)
		err = os.WriteFile(publicPath2, []byte(cloudsPublicYAML2), 0644)
		th.AssertNoErr(t, err)

		cloud, _, _, err := clouds.Parse(
			clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
			clouds.WithCloudName("gophercloud-test-0"),
			clouds.WithPublicLocations(publicPath0, publicPath1, publicPath2),
		)
		th.AssertNoErr(t, err)

		if got := authField(cloud, "auth_url"); got != "https://example.com/gophercloud-test-1:13000" {
			t.Errorf("unexpected identity endpoint: %q", got)
		}

		if got := authField(cloud, "user_domain_name"); got != "CustomDomain" {
			t.Errorf("unexpected UserDomainName: %q", got)
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

		cloud, _, _, err := clouds.Parse(
			clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
			clouds.WithCloudName("gophercloud-test"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := authField(cloud, "username"); got != "myuser" {
			t.Errorf("unexpected username: %q", got)
		}
		if got := authField(cloud, "user_domain_name"); got != "Default" {
			t.Errorf("unexpected user domain name: %q, expected 'Default'", got)
		}
		if got := authField(cloud, "project_name"); got != "myproject" {
			t.Errorf("unexpected project name: %q", got)
		}

		if got := authField(cloud, "project_domain_name"); got != "some_domain" {
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

		cloud, _, _, err := clouds.Parse(
			clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
			clouds.WithCloudName("gophercloud-test"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := authField(cloud, "username"); got != "myuser" {
			t.Errorf("unexpected username: %q", got)
		}
		if got := authField(cloud, "user_domain_id"); got != "default-domain-id" {
			t.Errorf("unexpected user domain ID: %q, expected 'default-domain-id'", got)
		}
		if got := authField(cloud, "project_id"); got != "project-123" {
			t.Errorf("unexpected project ID: %q", got)
		}
		if got := authField(cloud, "project_domain_id"); got != "other-domain-id" {
			t.Errorf("unexpected project domain ID: %q, expected 'other-domain-id'", got)
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

		cloud, _, _, err := clouds.Parse(
			clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
			clouds.WithCloudName("gophercloud-test"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := authField(cloud, "username"); got != "myuser" {
			t.Errorf("unexpected username: %q", got)
		}
		if got := authField(cloud, "domain_name"); got != "shared-domain" {
			t.Errorf("unexpected user domain name: %q, expected 'shared-domain'", got)
		}
		if got := authField(cloud, "project_name"); got != "myproject" {
			t.Errorf("unexpected scope project name: %q", got)
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

		cloud, _, _, err := clouds.Parse(
			clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
			clouds.WithCloudName("gophercloud-test"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := authField(cloud, "user_domain_name"); got != "user-specific-domain" {
			t.Errorf("unexpected user domain name: %q, expected 'user-specific-domain'", got)
		}
		if got := authField(cloud, "domain_name"); got != "fallback-domain" {
			t.Errorf("unexpected fallback domain name: %q, expected 'fallback-domain'", got)
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

		cloud, _, _, err := clouds.Parse(
			clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
			clouds.WithCloudName("gophercloud-test"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := authField(cloud, "domain_name"); got != "fallback-domain" {
			t.Errorf("unexpected fallback domain name: %q, expected 'fallback-domain'", got)
		}
		if got := authField(cloud, "project_domain_name"); got != "project-specific-domain" {
			t.Errorf("unexpected project domain name: %q, expected 'project-specific-domain'", got)
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

		cloud, _, _, err := clouds.Parse(
			clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
			clouds.WithCloudName("gophercloud-test"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := authField(cloud, "project_id"); got != "unique-project-id-123" {
			t.Errorf("unexpected project ID: %q", got)
		}
		if authField(cloud, "domain_id") != "" || authField(cloud, "domain_name") != "" {
			t.Errorf("expected no domain information when using project_id, got domain_id=%q, domain_name=%q",
				authField(cloud, "domain_id"), authField(cloud, "domain_name"))
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

		cloud, _, _, err := clouds.Parse(
			clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
			clouds.WithCloudName("gophercloud-test"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := authField(cloud, "username"); got != "admin" {
			t.Errorf("unexpected username: %q", got)
		}
		if got := authField(cloud, "user_domain_name"); got != "Default" {
			t.Errorf("unexpected user domain name: %q", got)
		}
		if got := authField(cloud, "system_scope"); got != "all" {
			t.Errorf("unexpected system scope: %q", got)
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

		cloud, _, _, err := clouds.Parse(
			clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
			clouds.WithCloudName("gophercloud-test"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := authField(cloud, "system_scope"); got != "all" {
			t.Errorf("unexpected system scope: %q", got)
		}
		if got := authField(cloud, "project_name"); got != "should-be-ignored" {
			t.Errorf("unexpected project name: %q", got)
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

		cloud, _, _, err := clouds.Parse(
			clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
			clouds.WithCloudName("gophercloud-test"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := authField(cloud, "username"); got != "domainadmin" {
			t.Errorf("unexpected username: %q", got)
		}
		if got := authField(cloud, "user_domain_name"); got != "Default" {
			t.Errorf("unexpected user domain name: %q", got)
		}
		if got := authField(cloud, "domain_id"); got != "domain-123" {
			t.Errorf("unexpected scope domain ID: %q, expected 'domain-123'", got)
		}
		if authField(cloud, "project_id") != "" || authField(cloud, "project_name") != "" {
			t.Errorf("expected no project in domain scope, got project_id=%q, project_name=%q",
				authField(cloud, "project_id"), authField(cloud, "project_name"))
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

		cloud, _, _, err := clouds.Parse(
			clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
			clouds.WithCloudName("gophercloud-test"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := authField(cloud, "username"); got != "domainadmin" {
			t.Errorf("unexpected username: %q", got)
		}
		if got := authField(cloud, "user_domain_name"); got != "UserDomain" {
			t.Errorf("unexpected user domain name: %q, expected 'UserDomain'", got)
		}
		if got := authField(cloud, "domain_name"); got != "ScopeDomain" {
			t.Errorf("unexpected scope domain name: %q, expected 'ScopeDomain'", got)
		}
		if authField(cloud, "project_id") != "" || authField(cloud, "project_name") != "" {
			t.Errorf("expected no project in domain scope, got project_id=%q, project_name=%q",
				authField(cloud, "project_id"), authField(cloud, "project_name"))
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

		cloud, _, _, err := clouds.Parse(
			clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
			clouds.WithCloudName("gophercloud-test"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := authField(cloud, "project_name"); got != "myproject" {
			t.Errorf("expected project scope, got project_name=%q", got)
		}
		if got := authField(cloud, "project_domain_name"); got != "ProjectDomain" {
			t.Errorf("expected project domain, got project_domain_name=%q", got)
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

		cloud, _, _, err := clouds.Parse(
			clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
			clouds.WithCloudName("gophercloud-test"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := authField(cloud, "user_domain_id"); got != "user-domain-id-123" {
			t.Errorf("unexpected user domain ID: %q, expected 'user-domain-id-123'", got)
		}
		if got := authField(cloud, "domain_id"); got != "scope-domain-id-456" {
			t.Errorf("unexpected scope domain ID: %q, expected 'scope-domain-id-456'", got)
		}
	})
}
