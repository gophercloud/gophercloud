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
}
