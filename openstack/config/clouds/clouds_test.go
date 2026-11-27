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
		th.AssertNoErr(t, err)
		defer rmTmpDirOrPanic(tmpDir)

		cwd, err := os.Getwd()
		th.AssertNoErr(t, err)
		err = os.Chdir(tmpDir)
		th.AssertNoErr(t, err)
		defer func() {
			if err := os.Chdir(cwd); err != nil {
				panic("unable to reset the current working directory: " + err.Error())
			}
		}()

		err = os.WriteFile("clouds.yaml", []byte(cloudsYAML), 0644)
		th.AssertNoErr(t, err)

		err = os.WriteFile("secure.yaml", []byte(secureYAML), 0644)
		th.AssertNoErr(t, err)

		ao, _, _, err := clouds.Parse(
			clouds.WithCloudName("gophercloud-test"),
		)
		th.AssertNoErr(t, err)

		th.AssertEquals(t, "https://example.com/gophercloud-test-12345:13000", ao.IdentityEndpoint)
		th.AssertEquals(t, "gophercloud-test-username", ao.Username)
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
		th.AssertNoErr(t, err)
		defer rmTmpDirOrPanic(tmpDir1)

		tmpDir2, err := os.MkdirTemp(os.TempDir(), tempDirPrefix)
		th.AssertNoErr(t, err)
		defer rmTmpDirOrPanic(tmpDir2)

		cloudsPath1, cloudsPath2 := path.Join(tmpDir1, "clouds.yaml"), path.Join(tmpDir2, "clouds.yaml")

		err = os.WriteFile(cloudsPath1, []byte(cloudsYAML1), 0644)
		th.AssertNoErr(t, err)
		err = os.WriteFile(cloudsPath2, []byte(cloudsYAML2), 0644)
		th.AssertNoErr(t, err)

		ao, _, _, err := clouds.Parse(
			clouds.WithCloudName("gophercloud-test"),
			clouds.WithLocations(cloudsPath1, cloudsPath2),
		)
		th.AssertNoErr(t, err)

		th.AssertEquals(t, "https://example.com/gophercloud-test-1:13000", ao.IdentityEndpoint)
	})

	t.Run("uses XDG_CONFIG_HOME for clouds.yaml location", func(t *testing.T) {
		const cloudsYAML = `clouds:
  gophercloud-test:
    auth:
      auth_url: https://example.com/xdg-config:13000`

		// Create a temp dir to use as XDG_CONFIG_HOME with clouds.yaml inside.
		xdgDir, err := os.MkdirTemp(os.TempDir(), tempDirPrefix)
		th.AssertNoErr(t, err)
		defer rmTmpDirOrPanic(xdgDir)

		openstackDir := path.Join(xdgDir, "openstack")
		err = os.MkdirAll(openstackDir, 0755)
		th.AssertNoErr(t, err)
		err = os.WriteFile(path.Join(openstackDir, "clouds.yaml"), []byte(cloudsYAML), 0644)
		th.AssertNoErr(t, err)

		// Change to an empty temp dir so cwd has no clouds.yaml.
		emptyDir, err := os.MkdirTemp(os.TempDir(), tempDirPrefix)
		th.AssertNoErr(t, err)
		defer rmTmpDirOrPanic(emptyDir)

		cwd, err := os.Getwd()
		th.AssertNoErr(t, err)
		err = os.Chdir(emptyDir)
		th.AssertNoErr(t, err)
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
		th.AssertNoErr(t, err)

		th.AssertEquals(t, "https://example.com/xdg-config:13000", ao.IdentityEndpoint)
	})

	t.Run("falls back to ~/.config when XDG_CONFIG_HOME is not set", func(t *testing.T) {
		const cloudsYAML = `clouds:
  gophercloud-test:
    auth:
      auth_url: https://example.com/home-config:13000`

		// Create a temp dir to use as HOME with clouds.yaml in .config/openstack/.
		homeDir, err := os.MkdirTemp(os.TempDir(), tempDirPrefix)
		th.AssertNoErr(t, err)
		defer rmTmpDirOrPanic(homeDir)

		openstackDir := path.Join(homeDir, ".config", "openstack")
		err = os.MkdirAll(openstackDir, 0755)
		th.AssertNoErr(t, err)
		err = os.WriteFile(path.Join(openstackDir, "clouds.yaml"), []byte(cloudsYAML), 0644)
		th.AssertNoErr(t, err)

		// Change to an empty temp dir so cwd has no clouds.yaml.
		emptyDir, err := os.MkdirTemp(os.TempDir(), tempDirPrefix)
		th.AssertNoErr(t, err)
		defer rmTmpDirOrPanic(emptyDir)

		cwd, err := os.Getwd()
		th.AssertNoErr(t, err)
		err = os.Chdir(emptyDir)
		th.AssertNoErr(t, err)
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
		th.AssertNoErr(t, err)

		th.AssertEquals(t, "https://example.com/home-config:13000", ao.IdentityEndpoint)
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
		th.AssertNoErr(t, err)
		defer rmTmpDirOrPanic(tmpDir0)

		tmpDir1, err := os.MkdirTemp(os.TempDir(), tempDirPrefix)
		th.AssertNoErr(t, err)
		defer rmTmpDirOrPanic(tmpDir1)

		tmpDir2, err := os.MkdirTemp(os.TempDir(), tempDirPrefix)
		th.AssertNoErr(t, err)
		defer rmTmpDirOrPanic(tmpDir2)

		cloudsPath0, cloudsPath1, cloudsPath2 := path.Join(tmpDir0, "clouds.yaml"), path.Join(tmpDir1, "clouds.yaml"), path.Join(tmpDir2, "clouds.yaml")

		err = os.WriteFile(cloudsPath1, []byte(cloudsYAML1), 0644)
		th.AssertNoErr(t, err)
		err = os.WriteFile(cloudsPath2, []byte(cloudsYAML2), 0644)
		th.AssertNoErr(t, err)

		ao, _, _, err := clouds.Parse(
			clouds.WithCloudName("gophercloud-test"),
			clouds.WithLocations(cloudsPath0, cloudsPath1, cloudsPath2),
		)
		th.AssertNoErr(t, err)

		th.AssertEquals(t, "https://example.com/gophercloud-test-1:13000", ao.IdentityEndpoint)
	})

	t.Run("parses clouds-public.yaml if present", func(t *testing.T) {
		const cloudsYAML = `clouds:
  gophercloud-test-0:
    profile: gophercloud-test-1
    auth:
      domain_name: CustomDomain`
		const cloudsPublicYAML = `public-clouds:
  gophercloud-test-1:
    auth:
      auth_url: https://example.com/gophercloud-test-1:13000
      domain_name: Default`

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

		ao, _, _, err := clouds.Parse(
			clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
			clouds.WithCloudName("gophercloud-test-0"),
		)
		th.AssertNoErr(t, err)

		th.AssertEquals(t, "https://example.com/gophercloud-test-1:13000", ao.IdentityEndpoint)
		th.AssertEquals(t, "CustomDomain", ao.DomainName)
	})
	t.Run("parses clouds-public.yaml locations in order", func(t *testing.T) {
		const cloudsYAML = `clouds:
  gophercloud-test-0:
    cloud: gophercloud-test-1
    auth:
      domain_name: CustomDomain`
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

		ao, _, _, err := clouds.Parse(
			clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
			clouds.WithCloudName("gophercloud-test-0"),
			clouds.WithPublicLocations(publicPath1, publicPath2),
		)
		th.AssertNoErr(t, err)

		th.AssertEquals(t, "https://example.com/gophercloud-test-1:13000", ao.IdentityEndpoint)
		th.AssertEquals(t, "CustomDomain", ao.DomainName)
	})
	t.Run("fall back to next location if clouds-public.yaml is not found", func(t *testing.T) {
		const cloudsYAML = `clouds:
  gophercloud-test-0:
    cloud: gophercloud-test-1
    auth:
      domain_name: CustomDomain`
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

		ao, _, _, err := clouds.Parse(
			clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
			clouds.WithCloudName("gophercloud-test-0"),
			clouds.WithPublicLocations(publicPath0, publicPath1, publicPath2),
		)
		th.AssertNoErr(t, err)

		th.AssertEquals(t, "https://example.com/gophercloud-test-1:13000", ao.IdentityEndpoint)
		th.AssertEquals(t, "CustomDomain", ao.DomainName)
	})

	t.Run("don't read fs secure.yaml if WithCloudsYAML is set", func(t *testing.T) {
		const cloudsYAML = `clouds:
  gophercloud-test:
    auth:
      auth_url: https://example.com/clouds:13000`
		const secureYAMLFS = `clouds:
  gophercloud-test:
    auth:
      auth_url: https://example.com/secure:13000`

		// Create a working dir to test if fs is read
		wDir, err := os.MkdirTemp(os.TempDir(), tempDirPrefix)
		th.AssertNoErr(t, err)
		defer rmTmpDirOrPanic(wDir)

		th.AssertNoErr(t, os.WriteFile(path.Join(wDir, "secure.yaml"), []byte(secureYAMLFS), 0644))

		cwd, err := os.Getwd()
		th.AssertNoErr(t, err)
		th.AssertNoErr(t, os.Chdir(wDir))
		defer func() {
			if err := os.Chdir(cwd); err != nil {
				panic("unable to reset the current working directory: " + err.Error())
			}
		}()

		ao, _, _, err := clouds.Parse(
			clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
			clouds.WithCloudName("gophercloud-test"),
		)
		th.AssertNoErr(t, err)

		if got := ao.IdentityEndpoint; got != "https://example.com/clouds:13000" {
			t.Errorf("unexpected identity endpoint: %q", got)
		}
	})

	t.Run("read both clouds.yaml and secure.yaml when OS_CLIENT_CONFIG_FILE is set", func(t *testing.T) {
		const cloudsYAML = `clouds:
  gophercloud-test:
    auth:
      auth_url: https://example.com/clouds:13000
      password: non-secure`
		const secureYAML = `clouds:
  gophercloud-test:
    auth:
      auth_url: https://example.com/secure:13000
      password: secure`

		wDir, err := os.MkdirTemp(os.TempDir(), tempDirPrefix)
		th.AssertNoErr(t, err)
		defer rmTmpDirOrPanic(wDir)

		th.AssertNoErr(t, os.WriteFile(path.Join(wDir, "clouds.yaml"), []byte(cloudsYAML), 0644))
		th.AssertNoErr(t, os.WriteFile(path.Join(wDir, "secure.yaml"), []byte(secureYAML), 0644))

		cwd, err := os.Getwd()
		th.AssertNoErr(t, err)
		th.AssertNoErr(t, os.Chdir(wDir))
		defer func() {
			if err := os.Chdir(cwd); err != nil {
				panic("unable to reset the current working directory: " + err.Error())
			}
		}()

		t.Setenv("OS_CLIENT_CONFIG_FILE", path.Join(wDir, "clouds.yaml"))

		ao, _, _, err := clouds.Parse(
			clouds.WithCloudName("gophercloud-test"),
		)
		th.AssertNoErr(t, err)

		if got := ao.Password; got != "secure" {
			t.Errorf("unexpected password: %q", got)
		}
	})

	t.Run("lookup an aside secure.yaml if WithLocations is set without WithSecureLocations", func(t *testing.T) {
		const cloudsYAML = `clouds:
  gophercloud-test:
    auth:
      auth_url: https://example.com/clouds:13000
      password: non-secure`
		const secureYAML = `clouds:
  gophercloud-test:
    auth:
      auth_url: https://example.com/secure:13000
      password: secure`

		// Create a custom dir to lookup in
		dir, err := os.MkdirTemp(os.TempDir(), tempDirPrefix)
		th.AssertNoErr(t, err)
		defer rmTmpDirOrPanic(dir)

		th.AssertNoErr(t, os.WriteFile(path.Join(dir, "clouds.yaml"), []byte(cloudsYAML), 0644))
		th.AssertNoErr(t, os.WriteFile(path.Join(dir, "secure.yaml"), []byte(secureYAML), 0640))

		ao, _, _, err := clouds.Parse(
			clouds.WithLocations(path.Join(dir, "clouds.yaml")),
			clouds.WithCloudName("gophercloud-test"),
		)
		th.AssertNoErr(t, err)

		if got := ao.Password; got != "secure" {
			t.Errorf("unexpected password: %q", got)
		}
	})

	t.Run("don't look up aside secure if both WithLocations and WithSecureLocations are given", func(t *testing.T) {
		const cloudsYAML = `clouds:
  gophercloud-test:
    auth:
      auth_url: https://example.com/clouds:13000
      password: non-secure`
		const secureYAML1 = `clouds:
  gophercloud-test:
    auth:
      auth_url: https://example.com/secure:13000
      password: secure-1`
		const secureYAML2 = `clouds:
  gophercloud-test:
    auth:
      auth_url: https://example.com/secure:13000
      password: secure-2`

		dir1, err := os.MkdirTemp(os.TempDir(), tempDirPrefix)
		th.AssertNoErr(t, err)
		defer rmTmpDirOrPanic(dir1)

		dir2, err := os.MkdirTemp(os.TempDir(), tempDirPrefix)
		th.AssertNoErr(t, err)
		defer rmTmpDirOrPanic(dir2)

		th.AssertNoErr(t, os.WriteFile(path.Join(dir1, "clouds.yaml"), []byte(cloudsYAML), 0644))
		th.AssertNoErr(t, os.WriteFile(path.Join(dir1, "secure.yaml"), []byte(secureYAML1), 0640))
		th.AssertNoErr(t, os.WriteFile(path.Join(dir2, "secure.yaml"), []byte(secureYAML2), 0640))

		ao, _, _, err := clouds.Parse(
			clouds.WithLocations(path.Join(dir1, "clouds.yaml")),
			clouds.WithSecureLocations(path.Join(dir2, "secure.yaml")),
			clouds.WithCloudName("gophercloud-test"),
		)
		th.AssertNoErr(t, err)

		if got := ao.Password; got != "secure-2" {
			t.Errorf("unexpected password: %q", got)
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
		th.AssertNoErr(t, err)

		th.AssertEquals(t, "myuser", ao.Username)
		th.AssertEquals(t, "", ao.DomainName)
		th.AssertEquals(t, "myproject", ao.TenantName)
		th.AssertTrue(t, ao.Scope != nil)
		th.AssertEquals(t, "myproject", ao.Scope.ProjectName)
		th.AssertEquals(t, "some_domain", ao.Scope.DomainName)
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
		th.AssertNoErr(t, err)

		th.AssertEquals(t, "myuser", ao.Username)
		th.AssertEquals(t, "", ao.DomainID)
		th.AssertEquals(t, "project-123", ao.TenantID)
		th.AssertTrue(t, ao.Scope != nil)
		th.AssertEquals(t, "project-123", ao.Scope.ProjectID)
		// When using project_id, the domain is not needed in scope
		th.AssertEquals(t, "", ao.Scope.DomainID)
	})

	t.Run("domain_name sets user identity but does not fall back to project scope", func(t *testing.T) {
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
		th.AssertNoErr(t, err)

		th.AssertEquals(t, "myuser", ao.Username)
		th.AssertEquals(t, "shared-domain", ao.DomainName)
		th.AssertTrue(t, ao.Scope != nil)
		th.AssertEquals(t, "myproject", ao.Scope.ProjectName)
		th.AssertEquals(t, "", ao.Scope.DomainName)
	})

	t.Run("domain_name sets user identity and user_domain_name is not used", func(t *testing.T) {
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
		th.AssertNoErr(t, err)

		th.AssertEquals(t, "fallback-domain", ao.DomainName)
		th.AssertTrue(t, ao.Scope != nil)
		th.AssertEquals(t, "", ao.Scope.DomainName)
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
		th.AssertNoErr(t, err)

		th.AssertEquals(t, "fallback-domain", ao.DomainName)
		th.AssertTrue(t, ao.Scope != nil)
		th.AssertEquals(t, "project-specific-domain", ao.Scope.DomainName)
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
		th.AssertNoErr(t, err)

		th.AssertTrue(t, ao.Scope != nil)
		th.AssertEquals(t, "unique-project-id-123", ao.Scope.ProjectID)
		th.AssertEquals(t, "", ao.Scope.DomainID)
		th.AssertEquals(t, "", ao.Scope.DomainName)
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
		th.AssertNoErr(t, err)

		th.AssertEquals(t, "admin", ao.Username)
		th.AssertEquals(t, "", ao.DomainName)
		th.AssertTrue(t, ao.Scope != nil)
		th.AssertTrue(t, ao.Scope.System)
		th.AssertEquals(t, "", ao.Scope.ProjectID)
		th.AssertEquals(t, "", ao.Scope.ProjectName)
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
		th.AssertNoErr(t, err)

		th.AssertTrue(t, ao.Scope != nil)
		th.AssertTrue(t, ao.Scope.System)
		th.AssertEquals(t, "", ao.Scope.ProjectName)
	})

	t.Run("fails system_scope when value is not all", func(t *testing.T) {
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

		_, _, _, err := clouds.Parse(
			clouds.WithCloudsYAML(strings.NewReader(cloudsYAML)),
			clouds.WithCloudName("gophercloud-test"),
		)
		th.AssertErr(t, err)
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
		th.AssertNoErr(t, err)

		th.AssertEquals(t, "domainadmin", ao.Username)
		th.AssertEquals(t, "", ao.DomainName)
		th.AssertTrue(t, ao.Scope != nil)
		th.AssertEquals(t, "domain-123", ao.Scope.DomainID)
		th.AssertEquals(t, "", ao.Scope.ProjectID)
		th.AssertEquals(t, "", ao.Scope.ProjectName)
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
		th.AssertNoErr(t, err)

		th.AssertEquals(t, "domainadmin", ao.Username)
		th.AssertEquals(t, "ScopeDomain", ao.DomainName)
		th.AssertTrue(t, ao.Scope != nil)
		th.AssertEquals(t, "ScopeDomain", ao.Scope.DomainName)
		th.AssertEquals(t, "", ao.Scope.ProjectID)
		th.AssertEquals(t, "", ao.Scope.ProjectName)
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
		th.AssertNoErr(t, err)

		th.AssertTrue(t, ao.Scope != nil)
		th.AssertEquals(t, "myproject", ao.Scope.ProjectName)
		th.AssertEquals(t, "ProjectDomain", ao.Scope.DomainName)
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
		th.AssertNoErr(t, err)

		th.AssertEquals(t, "scope-domain-id-456", ao.DomainID)
		th.AssertTrue(t, ao.Scope != nil)
		th.AssertEquals(t, "scope-domain-id-456", ao.Scope.DomainID)
	})
}
