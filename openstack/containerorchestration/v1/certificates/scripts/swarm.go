package scripts

import (
	"bytes"
	"html/template"

	"github.com/gophercloud/gophercloud/openstack/containerorchestration/v1/bays"
)

type swarmWriter struct{}

func (w *swarmWriter) Generate(bay *bays.Bay) map[string][]byte {
	scripts := make(map[string][]byte)

	data := w.getScriptData(bay)
	scripts["docker.env"] = w.buildBashScript(data)
	scripts["docker.cmd"] = w.buildCmdScript(data)
	scripts["docker.ps1"] = w.buildPs1Script(data)
	scripts["docker.fish"] = w.buildFishScript(data)
	return scripts
}

func (w *swarmWriter) getScriptData(bay *bays.Bay) interface{} {
	var data struct {
		DockerHost    string
		DockerVersion string
	}
	data.DockerHost = bay.COEEndpoint
	data.DockerVersion = bay.ContainerVersion

	return data
}

func (w *swarmWriter) buildBashScript(data interface{}) []byte {
	t := template.Must(template.New("docker.env").Parse(`__CARINA_ENV_SOURCE="$_"
if [ -n "$BASH_SOURCE" ]; then
  __CARINA_ENV_SOURCE="${BASH_SOURCE[0]}"
fi
DIR="$(cd "$(dirname "${__CARINA_ENV_SOURCE:-$0}")" > /dev/null && \pwd)"
unset __CARINA_ENV_SOURCE 2> /dev/null

export DOCKER_HOST={{.DockerHost}}
export DOCKER_TLS_VERIFY=1
export DOCKER_CERT_PATH=$DIR
export DOCKER_VERSION={{.DockerVersion}}
`))

	var script bytes.Buffer
	t.Execute(&script, data)
	return script.Bytes()
}

func (w *swarmWriter) buildCmdScript(data interface{}) []byte {
	t := template.Must(template.New("docker.cmd").Parse(`set DOCKER_HOST={{.DockerHost}}
set DOCKER_TLS_VERIFY=1
set DOCKER_CERT_PATH=%~dp0
set DOCKER_VERSION={{.DockerVersion}}
  `))

	var script bytes.Buffer
	t.Execute(&script, data)
	return script.Bytes()
}

func (w *swarmWriter) buildPs1Script(data interface{}) []byte {
	t := template.Must(template.New("docker.cmd").Parse(`$env:DOCKER_HOST="{{.DockerHost}}"
$env:DOCKER_TLS_VERIFY=1
$env:DOCKER_CERT_PATH=$PSScriptRoot
$env:DOCKER_VERSION="{{.DockerVersion}}"
`))

	var script bytes.Buffer
	t.Execute(&script, data)
	return script.Bytes()
}

func (w *swarmWriter) buildFishScript(data interface{}) []byte {
	t := template.Must(template.New("docker.fish").Parse(`set DIR (dirname (status -f))

set -x DOCKER_HOST {{.DockerHost}}
set -x DOCKER_TLS_VERIFY 1
set -x DOCKER_CERT_PATH $DIR
set -x DOCKER_VERSION {{.DockerVersion}}
`))

	var script bytes.Buffer
	t.Execute(&script, data)
	return script.Bytes()
}
