package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yellyoshua/whatsapp-chat-parser/api"
	"github.com/yellyoshua/whatsapp-chat-parser/constants"
	"github.com/yellyoshua/whatsapp-chat-parser/logger"
	"github.com/yellyoshua/whatsapp-chat-parser/paper"
	"github.com/yellyoshua/whatsapp-chat-parser/storage"
	"github.com/yellyoshua/whatsapp-chat-parser/utils"
	"github.com/yellyoshua/whatsapp-chat-parser/whatsapp"
)

type Handl struct {
	clientStorage storage.Uploader
}

func New() *Handl {
	clientStorage := storage.New()
	return &Handl{clientStorage: clientStorage}
}

func isQRFile(file_path string) bool {
	return strings.Contains(file_path, "qr/")
}

// merged_files_path includes a map with names of files and qr_files
func (h *Handl) backgroundUploadFiles(ctx *gin.Context, uuid string, attachmentsURL string, merged_files_path map[string]string) {
	file, header, _ := ctx.Request.FormFile("file")

	filesInZip, _ := utils.ExtractZipFile(file, header.Size)
	defer file.Close()

	var uploadsQueue chan string = make(chan string)
	var isDone chan error = make(chan error)
	var qr_files_paths map[string]string = make(map[string]string)

	for path_origin, http_uri := range merged_files_path {
		if isQRFile(http_uri) {
			qr_files_paths[path_origin] = http_uri
		}
	}

	go func(files map[string]io.Reader, ch <-chan string, isDone chan<- error) {
		uuid := <-ch
		qr_files := api.GenerateQR(attachmentsURL, qr_files_paths)

		for file_path, qr_file := range qr_files {
			filesInZip[file_path] = qr_file
		}

		err := uploadFilesS3(h.clientStorage, uuid, files)
		logger.Info("is uploaded %v files", len(files))

		isDone <- err
		close(isDone)
	}(filesInZip, uploadsQueue, isDone)

	go func(uuid string, ch chan<- string, isDone <-chan error) {
		ch <- uuid
		close(uploadsQueue)
		<-isDone
	}(uuid, uploadsQueue, isDone)
}

func (h *Handl) HolyShit(ctx *gin.Context) {
	defer ctx.Done()
	ctx.String(200, "Holy shit!")
}

func closeConnection(ctx *gin.Context) {
	ctx.Request.Context().Done()
	ctx.Done()
}

func (h *Handl) PostUploadChatFiles(ctx *gin.Context) {
	var attachmentURL string

	attachment_url, _ := ctx.Get(constants.KEY_MIDDLEWARE_ATTACHMENT_URL)
	attachmentURL = attachment_url.(string)

	responseFormat, _ := ctx.Params.Get("format")

	ctxCopy := ctx.Copy()
	uuid := utils.NewUniqueID()
	attachmentsURL := path.Join(attachmentURL, uuid)
	ctx.Request.ParseMultipartForm(10 << 20)
	file, header, _ := ctx.Request.FormFile("file")

	filesInZip, err := utils.ExtractZipFile(file, header.Size)
	defer file.Close()

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, "error reading zip file")
		closeConnection(ctx)
	} else {
		var filesToScanText map[string]io.Reader = make(map[string]io.Reader)
		var filesToFilterQR map[string]io.Reader = make(map[string]io.Reader)

		if err := utils.DuplicateReaders(filesInZip, filesToScanText, filesToFilterQR); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, "error duplicating readers")
			closeConnection(ctx)
		} else {

			var chat string
			var qrFilesPath map[string]string

			var isDone chan bool = make(chan bool)

			go func(isDone chan<- bool) {
				chat = api.ExtractChatFromFiles(filesToScanText)
				isDone <- true
			}(isDone)

			go func(isDone chan<- bool) {
				qrFilesPath = api.FilterQRFilesExtensions(filesToFilterQR)
				isDone <- true
			}(isDone)

			<-isDone
			<-isDone
			close(isDone)

			var mergedFiles map[string]string = make(map[string]string)
			mergedFiles = qrFilesPath

			for file_path := range filesToScanText {
				if len(mergedFiles[file_path]) == 0 {
					mergedFiles[file_path] = file_path
				}
			}

			h.backgroundUploadFiles(ctxCopy, uuid, attachmentsURL, mergedFiles)

			chatBuilder := whatsapp.New(uuid, chat)
			messages := chatBuilder.Messages()

			writer := paper.New(messages)
			book := writer.AttachFiles(mergedFiles, attachmentsURL)

			resWithFormat(ctx, book, responseFormat, uuid)
		}
	}
}

func PostParseOnlyChat(ctx *gin.Context) {
	uuid := utils.NewUniqueID()
	file, fileHeader, _ := ctx.Request.FormFile("file")

	files := map[string]io.Reader{
		fileHeader.Filename: file,
	}

	chat := api.ExtractChatFromFiles(files)

	chatBuilder := whatsapp.New(uuid, chat)
	messages := chatBuilder.Messages()

	writer := paper.New(messages)
	book := writer.AttachFiles(nil, "/")

	responseFormat, _ := ctx.Params.Get("format")
	resWithFormat(ctx, book, responseFormat, uuid)
}

func resWithFormat(ctx *gin.Context, book paper.Book, responseFormat string, uuid string) {
	switch format := utils.UniqLowerCase(responseFormat); format {
	case constants.FormatJSON:
		{
			var result interface{}

			messages, err := book.ExportJSON()

			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
				closeConnection(ctx)
			} else {
				if err := json.Unmarshal(messages.Value, &result); err != nil {
					ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
					closeConnection(ctx)
				} else {
					ctx.JSON(http.StatusOK, map[string]interface{}{
						"uuid":     uuid,
						"count":    messages.Count,
						"messages": result,
					})
					closeConnection(ctx)
				}
			}
		}
	case constants.FormatHTML:
		{
			messages, err := book.ExportHTML(paper.Minimal)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
				closeConnection(ctx)
			} else {
				ctx.String(http.StatusOK, messages)
				closeConnection(ctx)
			}
		}
	default:
		{
			ctx.String(http.StatusOK, "OK")
			closeConnection(ctx)
		}
	}
}

func uploadFilesS3(bucket storage.Uploader, uuid string, files map[string]io.Reader) error {
	files_copy := make(map[string]io.Reader)

	for file_path, f := range files {
		copyOfFile := new(bytes.Buffer)
		if err := utils.CopyReader(f, copyOfFile); err == nil {
			files_copy[path.Join(uuid, file_path)] = copyOfFile
		}
	}
	err := bucket.UploadFiles(files_copy)
	return err
}
