package paper

import (
	"bytes"
	"encoding/json"
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

var demo_messages []Message = []Message{
	{Date: "06_01_2020=23:07", Author: "S1", IsSender: false, IsReceiver: true, IsInfo: false, Message: "S1-m"},
	{Date: "06_01_2020=23:10", Author: "S2", IsSender: true, IsReceiver: false, IsInfo: false, Message: "S1-m"},
	{Date: "05_12_2020=20:27", Author: "S2", IsSender: true, IsReceiver: true, IsInfo: false, Message: "S1-m"},
	{Date: "05_13_2020=20:27", Author: "S2", IsSender: true, IsReceiver: true, IsInfo: false, Message: "S1-m"},
	{Date: "05_14_2020=20:27", Author: "S2", IsSender: true, IsReceiver: true, IsInfo: false, Message: "S1-m"},
}

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
	err := renderTemplate(testerTemplate, data, templateFuncs, receptor)

	assert.Equal(t, nil, err)
	assert.Equal(t, expected, receptor.String())
}

func TestUnmarshalMessagesAndSort(t *testing.T) {
	expected := 4
	json_messages, _ := json.Marshal(demo_messages)

	writer := New()
	book := writer.UnmarshalJSONMessages(string(json_messages), map[string]string{}, "home")
	messages := book.Export()

	isInfoCount := 0

	for _, m := range messages {
		if m.IsInfo {
			isInfoCount++
		}
	}

	assert.Equal(t, expected, isInfoCount)
}

func TestGetTranslateDate(t *testing.T) {
	expected := "Marzo 01, 2001"
	dateParsed := getTranslateDate("es", "03", "01", "2001")

	assert.Equal(t, expected, dateParsed)
}
