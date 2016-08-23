// +build windows

package scripts

import (
	"bytes"
	"strings"
	"text/template"
)

func compileTemplate(name string, contents string, data interface{}) []byte {
	// Encoding explicit carriage returns so that the file renders correctly
	// in notepad on Windows
	contents = strings.Replace(contents, "\n", "\r\n", -1)

	t := template.Must(template.New(name).Parse(contents))

	var script bytes.Buffer
	t.Execute(&script, data)

	return script.Bytes()
}
