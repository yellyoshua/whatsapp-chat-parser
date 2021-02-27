package utils

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestCopyFile(t *testing.T) {
	var expected string = "Hola mundo"

	r := ioutil.NopCloser(strings.NewReader("Hola mundo"))
	r.Close()

	rCopy, err := CopyFile(r)
	if err != nil {
		t.Errorf("Error copying reader -> %s", err)
	}

	var text string

	if err := ValueOfTextFile(rCopy, &text); err != nil {
		t.Errorf("Error reading TextFile -> %s", err)
	}

	assert.Equal(t, text, expected)
}

func TestValueOfTextFile(t *testing.T) {
	var expected string = "Hola mundo"

	r := ioutil.NopCloser(strings.NewReader("Hola mundo"))
	r.Close()

	var text string

	if err := ValueOfTextFile(r, &text); err != nil {
		t.Errorf("Error reading TextFile -> %s", err)
	}

	assert.Equal(t, text, expected)
}
