package whatsapp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSafeLangSupported(t *testing.T) {
	value := "es"
	supported := safeLangSupported(value)
	assert.Equal(t, "es", supported)

	value1 := "gp"
	supported1 := safeLangSupported(value1)
	assert.Equal(t, "es", supported1)

	value2 := "en"
	supported2 := safeLangSupported(value2)
	assert.Equal(t, "en", supported2)
}
