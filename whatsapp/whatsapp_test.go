package whatsapp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParserMessages(t *testing.T) {
	chat := New("demo-123", "6/1/20, 23:07 - Yoshua: Que le pasó?")
	messages := chat.Messages()

	assert.Equal(t, 2, len(messages))
}
