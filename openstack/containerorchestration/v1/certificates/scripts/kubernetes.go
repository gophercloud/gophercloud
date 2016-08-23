package scripts

import "github.com/gophercloud/gophercloud/openstack/containerorchestration/v1/bays"

type kubernetesWriter struct{}

func (w *kubernetesWriter) Generate(bay *bays.Bay) map[string][]byte {
	scripts := make(map[string][]byte)

	data := w.getScriptData(bay)
	scripts["kubectl.config"] = w.buildKubernetesConfig(data)
	scripts["kubectl.env"] = w.buildBashScript(data)
	scripts["kubectl.cmd"] = w.buildCmdScript(data)
	scripts["kubectl.ps1"] = w.buildPs1Script(data)
	scripts["kubectl.fish"] = w.buildFishScript(data)
	return scripts
}

func (w *kubernetesWriter) getScriptData(bay *bays.Bay) interface{} {
	var data struct {
		Name          string
		Server        string
		DockerVersion string
	}
	data.Name = bay.Name
	data.Server = bay.COEEndpoint
	data.DockerVersion = bay.ContainerVersion

	return data
}

func (w *kubernetesWriter) buildKubernetesConfig(data interface{}) []byte {
	template := `apiVersion: v1
clusters:
- cluster:
    certificate-authority: ca.pem
    server: {{.Server}}
  name: {{.Name}}
contexts:
- context:
    cluster: {{.Name}}
    user: {{.Name}}
  name: {{.Name}}
current-context: {{.Name}}
kind: Config
users:
- name: {{.Name}}
  user:
    client-certificate: cert.pem
    client-key: key.pem
`

	return compileTemplate("kubectl.config", template, data)
}

func (w *kubernetesWriter) buildBashScript(data interface{}) []byte {
	template := `__CARINA_ENV_SOURCE="$_"
if [ -n "$BASH_SOURCE" ]; then
  __CARINA_ENV_SOURCE="${BASH_SOURCE[0]}"
fi
DIR="$(cd "$(dirname "${__CARINA_ENV_SOURCE:-$0}")" > /dev/null && \pwd)"
unset __CARINA_ENV_SOURCE 2> /dev/null

export KUBECONFIG=$DIR/kubectl.config
export DOCKER_VERSION={{.DockerVersion}}
`

	return compileTemplate("kubectl.env", template, data)
}

func (w *kubernetesWriter) buildCmdScript(data interface{}) []byte {
	template := `set KUBECONFIG=%~dp0kubectl.config
set DOCKER_VERSION={{.DockerVersion}}
`

	return compileTemplate("kubectl.cmd", template, data)
}

func (w *kubernetesWriter) buildPs1Script(data interface{}) []byte {
	template := `$env:KUBECONFIG="$PSScriptRoot\kubectl.config"
$env:DOCKER_VERSION="{{.DockerVersion}}"
`

	return compileTemplate("kubectl.ps1", template, data)
}

func (w *kubernetesWriter) buildFishScript(data interface{}) []byte {
	template := `set DIR (dirname (status -f))

set -x DOCKER_CERT_PATH $DIR/kubectl.config
set -x DOCKER_VERSION {{.DockerVersion}}
`

	return compileTemplate("kubectl.fish", template, data)
}
