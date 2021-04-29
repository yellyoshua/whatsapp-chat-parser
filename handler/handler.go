package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"sync"
	"time"

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
	var ch chan string = make(chan string)
	// ctxCopy := ctx.Copy()

	go func(ch chan string) {
		time.Sleep(time.Second * time.Duration(4))
		ch <- "Holy shit!"
	}(ch)

	go func(ch chan string) {
		time.Sleep(time.Second * time.Duration(5))
		ch <- "Holy shit 2!"
	}(ch)

	go func(ch chan string) {

		for m := range ch {
			logger.Info(m)
		}

		close(ch)
	}(ch)

	ctx.String(200, "Holy shit!")
}

func closeConnection(ctx *gin.Context) {
	ctx.Request.Context().Done()
	ctx.Done()
}

func PostUploadChatFiles(ctx *gin.Context) {
	// var clientStorage *storage.Uploader
	var attachmentURL string

	// client_storage, _ := ctx.Get(constants.KEY_MIDDLEWARE_CLIENT_STORAGE)
	attachment_url, _ := ctx.Get(constants.KEY_MIDDLEWARE_ATTACHMENT_URL)

	// clientStorage = client_storage.(*storage.Uploader)
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
		closeConnection(ctx)
	} else {
		var filesToScanText map[string]io.Reader = make(map[string]io.Reader)
		var filesToUpload map[string]io.Reader = make(map[string]io.Reader)
		var filesToFilterQR map[string]io.Reader = make(map[string]io.Reader)

		if err := utils.DuplicateReaders(filesInZip, filesToUpload, filesToScanText, filesToFilterQR); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, "error duplicating readers")
			closeConnection(ctx)
		} else {

			var chChat chan string = make(chan string)
			var chFilesReplacedWithQR chan map[string]io.Reader = make(chan map[string]io.Reader)
			var chQRFilesPath chan map[string]string = make(chan map[string]string)

			var wgChat sync.WaitGroup

			wgChat.Add(1)
			go api.ExtractChatFromFiles(filesToScanText, chChat, &wgChat)
			wgChat.Add(1)
			go api.FilterQRFilesExtensions(filesToFilterQR, chQRFilesPath, &wgChat)
			wgChat.Add(1)
			go api.GenerateQR(chQRFilesPath, chFilesReplacedWithQR, &wgChat)

			var chat string
			var qrFilesPath map[string]string

			go func() {
				wgChat.Wait()
				close(chChat)
				close(chQRFilesPath)
				close(chFilesReplacedWithQR)
			}()

			chat = <-chChat
			qrFilesPath = <-chQRFilesPath
			qrFiles := <-chFilesReplacedWithQR

			// TODO: upload files in background
			var chUploads chan error = make(chan error)
			go uploadFilesS3(uuid, filesToUpload, chUploads)
			go uploadFilesS3(uuid, qrFiles, chUploads)

			go func(chUploads chan error) {
				for ch := range chUploads {
					if ch != err {
						logger.Info("Error uploading files -> " + ch.Error())
					}
				}
				logger.Info("Files uploaded")
				close(chUploads)
			}(chUploads)

			book, err := api.ParseWhatsappChatMessages(chat, qrFilesPath, attachmentsURL)

			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, "error processing messages")
				closeConnection(ctx)
			} else {
				resWithFormat(ctx, book, responseFormat, uuid)
			}
		}
	}
}

func PostParseOnlyChat(ctx *gin.Context) {
	file, fileHeader, _ := ctx.Request.FormFile("file")

	files := map[string]io.Reader{
		fileHeader.Filename: file,
	}

	var wg sync.WaitGroup
	var chChat chan string = make(chan string)
	wg.Add(1)
	go api.ExtractChatFromFiles(files, chChat, &wg)
	wg.Wait()

	book, err := api.ParseWhatsappChatMessages(<-chChat, nil, "/")

	defer close(chChat)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, "error processing messages")
		closeConnection(ctx)
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

// TODO: Retorna un [signal: killed] de error
func uploadFilesS3(uuid string, files map[string]io.Reader, chUploads chan error) {
	// st := *clientStorage

	files_copy := make(map[string]io.Reader)

	for file_path, f := range files {
		copyOfFile := new(bytes.Buffer)
		if err := utils.CopyReader(f, copyOfFile); err == nil {
			files_copy[filepath.Join(uuid, file_path)] = copyOfFile
		}
	}

	st := storage.New()

	if err := st.UploadFiles(files_copy); err != nil {
		logger.Info("error uploading files -> %s", err)
		chUploads <- err
	} else {
		chUploads <- nil
	}
	chUploads <- nil
}
