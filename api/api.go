package api

import (
	"bytes"
	"io"
	"path"
	"regexp"
	"strings"

	"github.com/skip2/go-qrcode"
	"github.com/yellyoshua/whatsapp-chat-parser/logger"
	"github.com/yellyoshua/whatsapp-chat-parser/utils"
)

var ExtensionsForQR []string = []string{
	".opus",
	".mp4",
}

func findExtension(filename string) bool {
	for _, filter := range ExtensionsForQR {
		r, _ := regexp.Compile(filter)
		if exist := r.MatchString(filename); exist {
			return true
		}
	}
	return false
}

func FilterQRFilesExtensions(files map[string]io.Reader) map[string]string {
	var qrFilesPaths map[string]string = make(map[string]string)

	for path_file := range files {
		canReplaceWithQR := findExtension(path_file)

		if canReplaceWithQR {
			id := utils.NewUniqueID()
			qrPathFile := path.Join("qr", id+".png")
			qrFilesPaths[path_file] = qrPathFile
		}
	}

	return qrFilesPaths
}

func GenerateQR(attachmentURL string, qr_files_path map[string]string) map[string]io.Reader {
	var files_replaced_with_qr map[string]io.Reader = make(map[string]io.Reader)

	for path_file, qrHashPathFile := range qr_files_path {

		if len(qrHashPathFile) != 0 {
			q, _ := qrcode.New(path.Join(attachmentURL, path_file), qrcode.High)
			qrImage, _ := q.PNG(256)
			files_replaced_with_qr[qrHashPathFile] = bytes.NewReader(qrImage)
		}
	}

	return files_replaced_with_qr
}

func ExtractChatFromFiles(files map[string]io.Reader) string {
	var chat string

	for path_file, file := range files {
		if isTXT := strings.Contains(path_file, ".txt"); isTXT {
			if err := utils.ValueOfTextFile(file, &chat); err != nil {
				logger.Info("Error ValueOfTextFile -> %s -> file: %s", err.Error(), path_file)
			}
		}
	}

	return chat
}
