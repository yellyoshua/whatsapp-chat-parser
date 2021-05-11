package utils

import (
	"bytes"
	"fmt"
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

func TestStringToInt(t *testing.T) {
	val := "10"
	num := StringToInt(val)
	fmt.Println(num)
	assert.EqualValues(t, 10, num)

	val1 := "01"
	num1 := StringToInt(val1)
	fmt.Println(num1)
	assert.EqualValues(t, 1, num1)

	val2 := " 01"
	num2 := StringToInt(val2)
	fmt.Println(num2)
	assert.EqualValues(t, 1, num2)

	val3 := "a01"
	num3 := StringToInt(val3)
	fmt.Println(num3)
	assert.EqualValues(t, 0, num3)

	val4 := "2001"
	num4 := StringToInt(val4)
	fmt.Println(num4)
	assert.EqualValues(t, 2001, num4)

	val5 := ""
	num5 := StringToInt(val5)
	fmt.Println(num5)
	assert.EqualValues(t, 0, num5)
}

func TestPadStart(t *testing.T) {
	expected := "20002010"
	result := PadStart("10", "2000", 8)
	assert.Equal(t, expected, result)

	expected1 := "2010"
	result1 := PadStart("10", "2000", 4)
	assert.Equal(t, expected1, result1)

	expected2 := "02"
	result2 := PadStart("2", "0", 2)
	assert.Equal(t, expected2, result2)
}
