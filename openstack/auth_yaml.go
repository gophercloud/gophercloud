package openstack

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gophercloud/gophercloud"
	"gopkg.in/yaml.v2"
)

var (
	ErrNoCloudYaml = fmt.Errorf("Clouds.yaml file not found.")
)

type Clouds struct {
	Clouds map[string]struct {
		Profile       string   `yaml:"profile"`
		RegionName    string   `yaml:"region_name"`
		Regions       []string `yaml:"regions"`
		DNSAPIVersion string   `yaml:"dns_api_version"`
		Auth          struct {
			AuthURL     string `yaml:"auth_url"`
			Username    string `yaml:"username"`
			Password    string `yaml:"password"`
			ProjectName string `yaml:"project_name"`
			ProjectID   string `yaml:"project_id"`
		} `yaml:"auth,omitempty"`
	}
}

// AuthOptionsFromYaml fills out an identity.AuthOptions structure with the settings found in a specific
// cloud entry of a clouds.yaml file (as per http://docs.openstack.org/developer/os-client-config)
// If `env` is set to true, it will also override these values with the various
// OS_* environment variables.  The following variables provide sources of truth: OS_AUTH_URL, OS_USERNAME,
// OS_PASSWORD, OS_TENANT_ID, and OS_TENANT_NAME.  Of these, OS_USERNAME, OS_PASSWORD, and OS_AUTH_URL must
// have settings, or an error will result.  OS_TENANT_ID and OS_TENANT_NAME are optional.
func AuthOptionsFromYaml(cloud string, env bool) (gophercloud.AuthOptions, error) {
	clouds := &Clouds{}
	cloudsContent, err := cloudsYaml()
	if err != nil {
		return nilOptions, err
	}
	err = yaml.Unmarshal(cloudsContent, clouds)
	if err != nil {
		return nilOptions, fmt.Errorf("Failed to unmarshal YAML: %v", err)
	}

	auth := clouds.Clouds[cloud].Auth

	authURL := auth.AuthURL
	username := auth.Username
	password := auth.Password
	tenantID := auth.ProjectID
	tenantName := auth.ProjectName

	if env {
		if v := os.Getenv("OS_AUTH_URL"); v != "" {
			authURL = v
		}

		if v := os.Getenv("OS_USERNAME"); v != "" {
			username = v
		}

		if v := os.Getenv("OS_PASSWORD"); v != "" {
			password = v
		}

		if v := os.Getenv("OS_TENANT_ID"); v != "" {
			tenantID = v
		}

		if v := os.Getenv("OS_TENANT_NAME"); v != "" {
			tenantName = v
		}
	}

	if authURL == "" {
		err := gophercloud.ErrMissingInput{Argument: "authURL"}
		return nilOptions, err
	}

	if username == "" {
		err := gophercloud.ErrMissingInput{Argument: "username"}
		return nilOptions, err
	}

	if password == "" {
		err := gophercloud.ErrMissingInput{Argument: "password"}
		return nilOptions, err
	}

	ao := gophercloud.AuthOptions{
		IdentityEndpoint: authURL,
		Username:         username,
		Password:         password,
		TenantID:         tenantID,
		TenantName:       tenantName,
	}

	return ao, nil
}

func cloudsYaml() ([]byte, error) {
	if f, ok := cloudsYamlFile(os.Getenv("USER_CONFIG_DIR")); ok {
		return ioutil.ReadFile(f)
	}

	if f, ok := cloudsYamlFile(fmt.Sprintf("%s/.config/openstack", os.Getenv("HOME"))); ok {
		return ioutil.ReadFile(f)
	}

	if f, ok := cloudsYamlFile(os.Getenv("SITE_CONFIG_DIR")); ok {
		return ioutil.ReadFile(f)
	}

	if f, ok := cloudsYamlFile("/etc/openstack"); ok {
		return ioutil.ReadFile(f)
	}

	return []byte{}, ErrNoCloudYaml
}

func cloudsYamlFile(dirname string) (string, bool) {
	filename := fmt.Sprintf("%s/clouds.yaml", dirname)
	if f, err := os.Stat(filename); err == nil && f.Mode().IsRegular() {
		return filename, true
	}
	return "", false
}
