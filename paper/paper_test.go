package paper

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yellyoshua/whatsapp-chat-parser/whatsapp"
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
		Messages: []whatsapp.Message{
			{Date: whatsapp.DateFormat{}, Author: "Tester", IsSender: true, IsReceiver: false, Message: "", Attachment: whatsapp.Attachment{}},
			{Date: whatsapp.DateFormat{}, Author: "Tester2", IsSender: false, IsReceiver: true, Message: "", Attachment: whatsapp.Attachment{}},
		},
		Background: "red",
	}

	receptor := new(bytes.Buffer)
	err := renderTemplate(testerTemplate, data, templateFuncs, receptor)

	assert.Equal(t, nil, err)
	assert.Equal(t, expected, receptor.String())
}
