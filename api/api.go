package api

import (
	"bytes"
	"io"
	"path"
	"regexp"
	"strings"

	"github.com/skip2/go-qrcode"
	"github.com/yellyoshua/whatsapp-chat-parser/logger"
	"github.com/yellyoshua/whatsapp-chat-parser/paper"
	"github.com/yellyoshua/whatsapp-chat-parser/utils"
	"github.com/yellyoshua/whatsapp-chat-parser/whatsapp"
)

var FilterExtensions []string = []string{
	".img",
	".png",
	".jpg",
	".opus",
	".webp",
}

func findExtension(filename string) bool {
	for _, filter := range FilterExtensions {
		r, _ := regexp.Compile(filter)
		if exist := r.MatchString(filename); exist {
			return true
		}
	}
	return false
}

func GenerateQR(files map[string]io.Reader, files_replaced_with_qr chan map[string]io.Reader, qr_files_path chan map[string]string) {
	var qrFilesPaths map[string]string = make(map[string]string)
	var qrFiles map[string]io.Reader = make(map[string]io.Reader)

	for path_file := range files {
		canReplaceWithQR := findExtension(path_file)

		if canReplaceWithQR {
			id := utils.NewUniqueID()
			qrPathFile := path.Join("qr", id+".png")

			q, _ := qrcode.New(path_file, qrcode.High)
			qrImage, _ := q.PNG(256)
			qrFiles[qrPathFile] = bytes.NewReader(qrImage)
			qrFilesPaths[path_file] = qrPathFile
		}
	}

	files_replaced_with_qr <- qrFiles
	qr_files_path <- qrFilesPaths
}

func ExtractChatFromFiles(files map[string]io.Reader, chChat chan string) {
	var chat string

	for path_file, file := range files {
		if isTXT := strings.Contains(path_file, ".txt"); isTXT {
			if err := utils.ValueOfTextFile(file, &chat); err != nil {
				logger.Info("Error ValueOfTextFile -> %s -> file: %s", err.Error(), path_file)
			}
		}
	}

	chChat <- chat
}

func ParseWhatsappChatMessages(chat string, qr_files map[string]string, attachmentURL string) (paper.Book, error) {
	var messages string
	wp := whatsapp.New()
	writer := paper.New()

	if err := wp.ParserMessages([]byte(chat), &messages); err != nil {
		return nil, err
	}

	book := writer.UnmarshalMessagesAndSort(messages, qr_files, attachmentURL)
	return book, nil
}
