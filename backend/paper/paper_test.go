package paper

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testerTemplate string = `
{{- range .Messages}}
{{ if .IsSender -}}
<h1>hola {{.Author}} (sender)</h1>
{{- end}}{{ if .IsReceiver -}}
<h1>hola {{.Author}} (receiver)</h1>
{{- end}}
{{- end}}
`

func TestRenderTemplate(t *testing.T) {
	expected := `
<h1>hola Tester (sender)</h1>
<h1>hola Tester2 (receiver)</h1>
`
	data := BookData{
		Messages: []Message{
			{Date: "", Author: "Tester", IsSender: true, IsReceiver: false, Message: "", Attachment: Attachment{}},
			{Date: "", Author: "Tester2", IsSender: false, IsReceiver: true, Message: "", Attachment: Attachment{}},
		},
		Background: "red",
	}

	receptor := new(bytes.Buffer)
	err := renderTemplate(testerTemplate, data, receptor)

	assert.Equal(t, nil, err)
	assert.Equal(t, expected, receptor.String())
}
