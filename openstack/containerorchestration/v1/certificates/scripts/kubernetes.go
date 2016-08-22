package scripts

import (
	"bytes"
	"html/template"

	"github.com/gophercloud/gophercloud/openstack/containerorchestration/v1/bays"
)

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
	t := template.Must(template.New("kubectl.config").Parse(`apiVersion: v1
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
`))

	var script bytes.Buffer
	t.Execute(&script, data)
	return script.Bytes()
}

func (w *kubernetesWriter) buildBashScript(data interface{}) []byte {
	t := template.Must(template.New("kubectl.env").Parse(`__CARINA_ENV_SOURCE="$_"
if [ -n "$BASH_SOURCE" ]; then
  __CARINA_ENV_SOURCE="${BASH_SOURCE[0]}"
fi
DIR="$(cd "$(dirname "${__CARINA_ENV_SOURCE:-$0}")" > /dev/null && \pwd)"
unset __CARINA_ENV_SOURCE 2> /dev/null

export KUBECONFIG=$DIR/kubectl.config
export DOCKER_VERSION={{.DockerVersion}}
`))

	var script bytes.Buffer
	t.Execute(&script, data)
	return script.Bytes()
}

func (w *kubernetesWriter) buildCmdScript(data interface{}) []byte {
	t := template.Must(template.New("kubectl.cmd").Parse(`set KUBECONFIG=%~dp0\kubectl.config
set DOCKER_VERSION={{.DockerVersion}}
  `))

	var script bytes.Buffer
	t.Execute(&script, data)
	return script.Bytes()
}

func (w *kubernetesWriter) buildPs1Script(data interface{}) []byte {
	t := template.Must(template.New("kubectl.cmd").Parse(`$env:KUBECONFIG=$PSScriptRoot\kubectl.config
$env:DOCKER_VERSION="{{.DockerVersion}}"
`))

	var script bytes.Buffer
	t.Execute(&script, data)
	return script.Bytes()
}

func (w *kubernetesWriter) buildFishScript(data interface{}) []byte {
	t := template.Must(template.New("kubectl.fish").Parse(`set DIR (dirname (status -f))

set -x DOCKER_CERT_PATH $DIR/kubectl.config
set -x DOCKER_VERSION {{.DockerVersion}}
`))

	var script bytes.Buffer
	t.Execute(&script, data)
	return script.Bytes()
}
