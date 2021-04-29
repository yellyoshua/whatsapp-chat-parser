package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var files = []string{
	"PTT-20200602-WA0011.opus",
	"PTT-20200602.opus",
}

func TestFindExtension(t *testing.T) {
	for _, file_name := range files {
		canReplaceWithQR := findExtension(file_name)
		assert.Equal(t, true, canReplaceWithQR)
	}
}
