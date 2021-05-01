package paper

import (
	"bytes"
	"html/template"
	"testing"
)

type attachmentWrapper struct {
	Attachment Attachment
}

func TestCardImage(t *testing.T) {
	attachment := Attachment{
		FileName:  "holamundo.opus",
		Exist:     true,
		Extension: ".opus",
	}
	var buffer = new(bytes.Buffer)
	tmpl, errTemplate := template.New("book-rendered").Funcs(templateFuncs).Parse(cardImage())
	if errTemplate != nil {
		t.Error(errTemplate)
	}

	if err := tmpl.Execute(buffer, attachmentWrapper{Attachment: attachment}); err != nil {
		t.Error(err)
	}

}
