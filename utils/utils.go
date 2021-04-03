package utils

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/yellyoshua/whatsapp-chat-parser/logger"
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
