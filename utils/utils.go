package utils

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/google/uuid"
	"github.com/yellyoshua/whatsapp-chat-parser/logger"
)

// UniqLowerCase convert in lowecase
func UniqLowerCase(input string) string {
	return strings.ToLower(input)
}

// LowerCase convert many inputs in lowecase
func LowerCase(inputs ...*string) {
	for _, txt := range inputs {
		*txt = UniqLowerCase(*txt)
	}
}

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

func ReflectPointer(dest interface{}, val interface{}) {
	isPointer := func(val interface{}) bool {
		return reflect.TypeOf(val).Kind() == reflect.Ptr
	}

	if isPointer(dest) {
		rGopher := reflect.ValueOf(dest)
		rG2Val := reflect.ValueOf(val)
		rGopher.Elem().Set(rG2Val)
	}
}

func DuplicateReaders(readers map[string]io.Reader, destReaders ...map[string]io.Reader) error {
	// var buffers []*bytes.Buffer = make([]*bytes.Buffer, 0)
	for i, r := range readers {
		var buffer bytes.Buffer
		if err := CopyReader(r, &buffer); err != nil {
			return err
		}

		for _, d := range destReaders {
			d[i] = &buffer
		}
	}
	return nil
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

func ExtractZipFile(f multipart.File, size int64, folder string) (map[string]io.Reader, error) {
	var dst map[string]io.Reader = make(map[string]io.Reader)
	reader, err := zip.NewReader(f, size)
	if err != nil {
		return dst, err
	}

	for _, zipFile := range reader.File {
		var customReader bytes.Buffer
		f, _ := zipFile.Open()
		fullPath := filepath.Join(folder, zipFile.Name)

		if err := CopyReader(f, &customReader); err != nil {
			return dst, err
		}

		dst[fullPath] = &customReader

		if err := f.Close(); err != nil {
			logger.Info("error closing zip file -> %s", err)
			return dst, err
		}
	}

	return dst, nil
}
