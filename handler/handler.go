package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/yellyoshua/whatsapp-chat-parser/api"
	"github.com/yellyoshua/whatsapp-chat-parser/constants"
	"github.com/yellyoshua/whatsapp-chat-parser/logger"
	"github.com/yellyoshua/whatsapp-chat-parser/paper"
	"github.com/yellyoshua/whatsapp-chat-parser/storage"
	"github.com/yellyoshua/whatsapp-chat-parser/utils"
)

func HolyShit(ctx *gin.Context) {
	defer ctx.Done()
	ctx.String(200, "Holy shit!")
}

func PostUploadChatFiles(ctx *gin.Context) {
	var clientStorage *storage.Uploader
	var attachmentURL string

	client_storage, _ := ctx.Get(constants.KEY_MIDDLEWARE_CLIENT_STORAGE)
	attachment_url, _ := ctx.Get(constants.KEY_MIDDLEWARE_ATTACHMENT_URL)

	clientStorage = client_storage.(*storage.Uploader)
	attachmentURL = attachment_url.(string)

	responseFormat, _ := ctx.Params.Get("format")

	uuid := utils.NewUniqueID()
	attachmentsURL := filepath.Join(attachmentURL, uuid)
	ctx.Request.ParseMultipartForm(10 << 20)
	file, header, _ := ctx.Request.FormFile("file")

	filesInZip, err := utils.ExtractZipFile(file, header.Size)
	defer file.Close()

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, "error reading zip file")
		ctx.Done()
	} else {
		var filesToScanText map[string]io.Reader = make(map[string]io.Reader)
		var filesToUpload map[string]io.Reader = make(map[string]io.Reader)
		var filesToQR map[string]io.Reader = make(map[string]io.Reader)

		if err := utils.DuplicateReaders(filesInZip, filesToUpload, filesToScanText, filesToQR); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, "error duplicating readers")
			ctx.Done()
		} else {

			var isDone chan string = make(chan string)
			var chChat chan string = make(chan string)
			var chUploads chan error = make(chan error)
			var chUploadsQR chan error = make(chan error)
			var chFilesReplacedWithQR chan map[string]io.Reader = make(chan map[string]io.Reader)
			var chQRFilesURL chan map[string]string = make(chan map[string]string)

			go api.GenerateQR(filesToQR, chFilesReplacedWithQR, chQRFilesURL)
			go api.ExtractChatFromFiles(filesToScanText, chChat)
			// TODO: upload files in background
			go uploadFilesS3(uuid, filesToUpload, clientStorage, chUploads)
			go uploadFilesS3(uuid, <-chFilesReplacedWithQR, clientStorage, chUploadsQR)

			var chat string
			var qrFilesURL map[string]string
			var errUploads error
			var errQRUploads error

			go func(isDone chan string) {
				chat, qrFilesURL, errUploads, errQRUploads = <-chChat, <-chQRFilesURL, <-chUploads, <-chUploadsQR
				isDone <- "ok"
			}(isDone)

			defer func() {
				close(isDone)
				close(chFilesReplacedWithQR)
				close(chQRFilesURL)
				close(chChat)
				close(chUploads)
			}()

			<-isDone

			if errUploads != nil || errQRUploads != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, "Error uploading files -> "+errUploads.Error())
				ctx.Done()
			} else {
				book, err := api.ParseWhatsappChatMessages(chat, qrFilesURL, attachmentsURL)

				if err != nil {
					ctx.AbortWithStatusJSON(http.StatusInternalServerError, "error processing messages")
					ctx.Done()
				} else {
					resWithFormat(ctx, book, responseFormat, uuid)
				}
			}
		}
	}
}

func PostParseOnlyChat(ctx *gin.Context) {
	file, fileHeader, _ := ctx.Request.FormFile("file")

	files := map[string]io.Reader{
		fileHeader.Filename: file,
	}

	var chChat chan string = make(chan string)
	go api.ExtractChatFromFiles(files, chChat)

	book, err := api.ParseWhatsappChatMessages(<-chChat, nil, "/")

	defer close(chChat)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, "error processing messages")
		ctx.Done()
	} else {
		responseFormat, _ := ctx.Params.Get("format")
		resWithFormat(ctx, book, responseFormat, "")
	}
}

func resWithFormat(ctx *gin.Context, book paper.Book, responseFormat string, uuid string) {
	switch format := utils.UniqLowerCase(responseFormat); format {
	case constants.FormatJSON:
		{
			var result interface{}
			messages, err := book.ExportJSON()

			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
				ctx.Done()
			} else {
				if err := json.Unmarshal(messages.Value, &result); err != nil {
					ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
					ctx.Done()
				} else {
					ctx.JSON(http.StatusOK, map[string]interface{}{
						"uuid":     uuid,
						"count":    messages.Count,
						"messages": result,
					})
					ctx.Done()
				}
			}
		}
	case constants.FormatHTML:
		{
			messages, err := book.ExportHTML(paper.Minimal)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
				ctx.Done()
			} else {
				ctx.String(http.StatusOK, messages)
				ctx.Done()
			}
		}
	default:
		{
			ctx.String(http.StatusOK, "OK")
			ctx.Done()
		}
	}
}

func uploadFilesS3(uuid string, files map[string]io.Reader, clientStorage *storage.Uploader, chUploads chan error) {
	st := *clientStorage

	files_copy := files

	for file_path, f := range files_copy {
		files[filepath.Join(uuid, file_path)] = f
	}

	if err := st.UploadFiles(files_copy); err != nil {
		logger.Info("error uploading files -> %s", err)
		chUploads <- err
	} else {
		chUploads <- nil
	}
}
