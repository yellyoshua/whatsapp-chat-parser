package whatsapp

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestSplitChatMessages(t *testing.T) {
	m := splitChatMessages("6/1/20, 23:07 - Yoshua: Que le pasó?")
	assert.Equal(t, 1, len(m))

	m1 := splitChatMessages(inputMessages)
	assert.Equal(t, 19, len(m1))
}

func TestParsePlainMessages(t *testing.T) {
	m := splitChatMessages("6/1/20, 23:07 - Yoshua: Que le pasó?")
	rawMessages := parsePlainMessages(m)
	assert.Equal(t, 1, len(rawMessages))

	m1 := splitChatMessages(inputMessages)
	rawMessages1 := parsePlainMessages(m1)
	assert.Equal(t, 14, len(rawMessages1))
}

func TestRawMessagesToMessages(t *testing.T) {
	m := splitChatMessages("6/1/20, 23:07 - Yoshua: Que le pasó?")
	rawMessages := parsePlainMessages(m)
	messages := rawMessagesToMessages(rawMessages)

	assert.Equal(t, len(rawMessages)+1, len(messages))
}

func TestReplaceAttachment(t *testing.T) {
	attachmentNoSpaces := func() {
		message := "3/19/20, 19:00 - Yoshua Lopez: IMG-20200319-WA0011.jpg (file attached)\nJefecito ... un bug!!!!"
		expected := "3/19/20, 19:00 - Yoshua Lopez: <attached: IMG-20200319-WA0011.jpg>\nJefecito ... un bug!!!!"
		withAttachment := formatAttachment(message)

		assert.Equal(t, expected, withAttachment)
	}

	attachmentWithSpaces := func() {
		message := "3/19/20, 19:00 - Yoshua Lopez: Frank sinatra.jpg (file attached)\nJefecito ... un bug!!!!"
		expected := "3/19/20, 19:00 - Yoshua Lopez: <attached: Frank%20sinatra.jpg>\nJefecito ... un bug!!!!"
		withAttachment := formatAttachment(message)

		assert.Equal(t, expected, withAttachment)
	}

	attachmentNoSpaces()
	attachmentWithSpaces()
}
