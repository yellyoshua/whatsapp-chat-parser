package utils

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDuplicateReaders(t *testing.T) {
	readers := map[string]io.Reader{
		"id_v1": strings.NewReader("v1"),
		"id_v2": strings.NewReader("v2"),
		"id_v3": strings.NewReader("v3"),
	}
	var clone1 map[string]io.Reader = make(map[string]io.Reader)
	var clone2 map[string]io.Reader = make(map[string]io.Reader)
	if err := DuplicateReaders(readers, clone1, clone2); err != nil {
		t.Error(err)
	}

	assert.NotEqual(t, clone1["id_v1"], nil)

	clone1Result1, errClone1Result1 := ioutil.ReadAll(clone1["id_v1"])
	if errClone1Result1 != nil {
		t.Error(errClone1Result1)
	}

	clone1Result2, errClone1Result2 := ioutil.ReadAll(clone1["id_v2"])
	if errClone1Result2 != nil {
		t.Error(errClone1Result2)
	}

	clone1Result3, errClone1Result3 := ioutil.ReadAll(clone1["id_v3"])
	if errClone1Result3 != nil {
		t.Error(errClone1Result3)
	}

	assert.Equal(t, "v1", string(clone1Result1))
	assert.Equal(t, "v2", string(clone1Result2))
	assert.Equal(t, "v3", string(clone1Result3))
}
func TestCopyFile(t *testing.T) {
	var expected string = "Hola mundo"

	r := ioutil.NopCloser(strings.NewReader("Hola mundo"))
	r.Close()

	var rCopy bytes.Buffer

	if err := CopyReader(r, &rCopy); err != nil {
		t.Errorf("Error copying reader -> %s", err)
	}

	var text string

	if err := ValueOfTextFile(&rCopy, &text); err != nil {
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
