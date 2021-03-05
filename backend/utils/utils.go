package utils

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"github.com/google/uuid"
)

// NewUniqueID return a unique uuid string
func NewUniqueID() string {
	token := uuid.New()
	return token.String()
}

// GetCurrentPath __
func GetCurrentPath() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}

// CopyReader _
func CopyReader(src io.Reader, dst ...*bytes.Buffer) error {
	b, err := ioutil.ReadAll(src)
	if err != nil {
		return err
	}

	for i := 0; i < len(dst); i++ {
		s := bytes.NewBuffer(b)
		*dst[i] = *s
	}

	return nil
}

// ValueOfTextFile __
func ValueOfTextFile(f io.Reader, text *string) error {
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	*text = string(buf)

	return nil
}
