package whatsapp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yellyoshua/whatsapp-chat-parser/utils"
)

var inputMessages string = `
20/06/19, 15:59 - Messages to this group are now secured with end-to-end encryption.
20/06/19, 15:59 - Loris created group “WhatsApp Chat Parser Example 2”
20/06/19, 15:59 - Loris added Emily
20/06/19, 15:59 - Loris added John
20/06/19, 15:59 - John: Hey 👋
20/06/19, 15:59 - Loris: Welcome to the chat example!
20/06/19, 15:59 - John: Thanks
20/06/19, 15:59 - Loris: Is everybody here?
20/06/19, 15:59 - Emily: Yes
20/06/19, 15:59 - Loris: Good
10/08/19, 22:17 - Jony: La teoría A es q andas en la huelga mijo.. Jaja


10/08/19, 22:17 - Jony: La teoría A es q andas en la huelga mijo.. Jaja
10/09/19, 12:58 - Yoshua Lopez: Jaja calla oye, no hay buses para ir
03/09/20, 19:00 - Yoshua Lopez: IMG-20200319-WA0011.jpg (file attached)
Jefecito ... un bug!!!!
	`

func TestByteToStringMessages(t *testing.T) {

	messages := byteToStringMessages([]byte(inputMessages))

	assert.Equal(t, true, len(messages) < len(inputMessages))
}
func TestReplaceAttachment(t *testing.T) {
	attachmentNoSpaces := func() {
		message := "3/19/20, 19:00 - Yoshua Lopez: IMG-20200319-WA0011.jpg (file attached)\nJefecito ... un bug!!!!"
		expected := "3/19/20, 19:00 - Yoshua Lopez: <attached: IMG-20200319-WA0011.jpg>\nJefecito ... un bug!!!!"
		withAttachment := replaceAttachment(message)

		assert.Equal(t, expected, withAttachment)
	}

	attachmentWithSpaces := func() {
		message := "3/19/20, 19:00 - Yoshua Lopez: Frank sinatra.jpg (file attached)\nJefecito ... un bug!!!!"
		expected := "3/19/20, 19:00 - Yoshua Lopez: <attached: Frank%20sinatra.jpg>\nJefecito ... un bug!!!!"
		withAttachment := replaceAttachment(message)

		assert.Equal(t, expected, withAttachment)
	}

	attachmentNoSpaces()
	attachmentWithSpaces()
}
func TestParserMessages(t *testing.T) {
	var outputMessages string

	wp := New()
	uuid := "demo-" + utils.NewUniqueID()

	rawChat := wp.ChatParser(uuid, []byte(inputMessages))

	if err := rawChat.ParserMessages(&outputMessages); err != nil {
		t.Errorf("error parsing messages -> " + err.Error())
	}

	// assert.Equal(t, true, len(outputMessages) > 0)
}
