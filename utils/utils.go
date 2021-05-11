package utils

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math/bits"
	"mime/multipart"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/yellyoshua/whatsapp-chat-parser/constants"
	"github.com/yellyoshua/whatsapp-chat-parser/logger"
)

// UniqLowerCase convert in lowecase
func UniqLowerCase(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

func IsEqualString(a string, b string) bool {
	return UniqLowerCase(a) == UniqLowerCase(b)
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

func ExtractZipFile(f multipart.File, size int64) (map[string]io.Reader, error) {
	var dst map[string]io.Reader = make(map[string]io.Reader)
	reader, err := zip.NewReader(f, size)
	if err != nil {
		return dst, err
	}

	for _, zipFile := range reader.File {
		var customReader bytes.Buffer
		f, _ := zipFile.Open()
		fileURL := path.Join(zipFile.Name)

		if err := CopyReader(f, &customReader); err != nil {
			return dst, err
		}

		dst[fileURL] = &customReader

		if err := f.Close(); err != nil {
			logger.Info("error closing zip file -> %s", err)
			return dst, err
		}
	}

	return dst, nil
}

func GetAttachmentURL() string {
	var s3AttachmentURL string
	var defaultS3BucketName = constants.S3BucketName
	var s3BucketName = os.Getenv("S3_BUCKET_NAME")
	var s3BucketRegion = os.Getenv("AWS_REGION")

	if len(s3BucketName) == 0 {
		s3AttachmentURL = strings.Replace(constants.S3BucketEndpoint, "BUCKET_NAME", defaultS3BucketName, -1)
	} else {
		s3AttachmentURL = strings.Replace(constants.S3BucketEndpoint, "BUCKET_NAME", s3BucketName, -1)
	}

	return strings.Replace(s3AttachmentURL, "BUCKET_REGION", s3BucketRegion, -1)
}

func StringToInt(value string) int {
	cleanValue := strings.TrimSpace(value)
	num, _ := strconv.ParseInt(cleanValue, 0, 64)
	return int(num)
}

func IntToString(value int) string {
	return fmt.Sprintf("%v", value)
}

// PadStart (value = 10, pad_start = 2000, length = 4) -> 2010
func PadStart(value string, pad_start string, length int) string {
	var lengthValue uint = uint(len(value))
	var lengthExpected uint = uint(length)

	uintDiff, _ := bits.Sub(lengthExpected, lengthValue, 0)
	lengthDifference := int(uintDiff)

	var indexPad int = 1
	var valuePadStart string

	if lengthDifference > 0 {
		for i := 0; i < lengthDifference; i++ {
			if indexPad > len(pad_start) {
				indexPad = 1
			}

			if indexPad == 1 {
				valuePadStart = valuePadStart + pad_start[:indexPad]
			} else {
				valuePadStart = valuePadStart + pad_start[indexPad-1:indexPad]
			}

			indexPad++
		}
	}

	return valuePadStart + value
}
