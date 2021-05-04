package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsQRFile(t *testing.T) {
	assert.Equal(t, true, isQRFile("qr/4424ad33-8115-49bd-bc1a-26e930c72ab8.png"))
	assert.Equal(t, true, isQRFile("/qr/4424ad33-8115-49bd-bc1a-26e930c72ab8.png"))
	assert.Equal(t, true, isQRFile("https://aws.com/qr/4424ad33-8115-49bd-bc1a-26e930c72ab8.png"))
	assert.Equal(t, false, isQRFile("https://aws.com/4424ad33-8115-49bd-bc1a-26e930c72ab8.png"))
	assert.Equal(t, false, isQRFile("/4424ad33-8115-49bd-bc1a-26e930c72ab8.png"))
}
