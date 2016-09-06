// +build !windows

package scripts

import (
	"bytes"
	"text/template"
)

func compileTemplate(name string, contents string, data interface{}) []byte {
	t := template.Must(template.New(name).Parse(contents))

	var script bytes.Buffer
	t.Execute(&script, data)

	return script.Bytes()
}
