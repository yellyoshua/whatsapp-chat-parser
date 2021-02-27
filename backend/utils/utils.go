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

// CopyFile _
func CopyFile(f io.ReadCloser) (io.Reader, error) {
	// TODO: solve how to copy this buffer
	b, err := ioutil.ReadAll(f)
	return bytes.NewBuffer(b), err
	// var buf1 bytes.Buffer
	// w := io.MultiWriter(&buf1)

	// if _, err := io.Copy(w, f); err != nil {
	// 	return io.Reader(&buf1), err
	// }

	// return io.Reader(&buf1), nil
}

// // CopyFile _
// func CopyFile(f io.Reader) (io.Reader, error) {
// 	// b, err := ioutil.ReadAll(f)
// 	// return io.Reader(bytes.NewBuffer(b)), err
// 	var buf1 bytes.Buffer
// 	w := io.MultiWriter(&buf1)

// 	if _, err := io.Copy(w, f); err != nil {
// 		return io.Reader(&buf1), err
// 	}

// 	return io.Reader(&buf1), nil
// }

// ValueOfTextFile __
func ValueOfTextFile(f io.Reader, text *string) error {
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	*text = string(buf)

	return nil
}
