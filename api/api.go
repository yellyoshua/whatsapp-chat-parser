package api

import (
	"bytes"
	"io"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/skip2/go-qrcode"
	"github.com/yellyoshua/whatsapp-chat-parser/logger"
	"github.com/yellyoshua/whatsapp-chat-parser/paper"
	"github.com/yellyoshua/whatsapp-chat-parser/utils"
	"github.com/yellyoshua/whatsapp-chat-parser/whatsapp"
)

var FilterQRExtensions []string = []string{
	".opus",
	".mp4",
}

// TODO: Testear esta funcionalidad
func findExtension(filename string) bool {
	for _, filter := range FilterQRExtensions {
		r, _ := regexp.Compile(filter)
		if exist := r.MatchString(filename); exist {
			return true
		}
	}
	return false
}

func FilterQRFilesExtensions(files map[string]io.Reader, qr_files_path chan map[string]string, wg *sync.WaitGroup) {
	defer wg.Done()

	var qrFilesPaths map[string]string = make(map[string]string)

	for path_file := range files {
		canReplaceWithQR := findExtension(path_file)

		if canReplaceWithQR {
			id := utils.NewUniqueID()
			qrPathFile := path.Join("qr", id+".png")
			qrFilesPaths[path_file] = qrPathFile
		}
	}

	// Este canal se va a leer dentro de otra go routine y afuera.
	qr_files_path <- qrFilesPaths
	qr_files_path <- qrFilesPaths
}

func GenerateQR(qr_files_path <-chan map[string]string, files_replaced_with_qr chan map[string]io.Reader, wg *sync.WaitGroup) {
	qrPaths := <-qr_files_path

	defer wg.Done()

	var qrFiles map[string]io.Reader = make(map[string]io.Reader)

	for _, qrHashPathFile := range qrPaths {

		if len(qrHashPathFile) != 0 {
			q, _ := qrcode.New(qrHashPathFile, qrcode.High)
			qrImage, _ := q.PNG(256)
			qrFiles[qrHashPathFile] = bytes.NewReader(qrImage)
		}
	}

	files_replaced_with_qr <- qrFiles
}

func ExtractChatFromFiles(files map[string]io.Reader, chChat chan string, wg *sync.WaitGroup) {
	defer wg.Done()

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

func ParseWhatsappChatMessages(user_id string, chat string, qr_files map[string]string, attachmentURL string) (paper.Book, error) {
	var messages string
	wp := whatsapp.New()
	writer := paper.New()

	rawChat := wp.ChatParser(user_id, []byte(chat))

	if err := rawChat.ParserMessages(&messages); err != nil {
		return nil, err
	}

	book := writer.UnmarshalMessagesAndSort(messages, qr_files, attachmentURL)
	return book, nil
}
