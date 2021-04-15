package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/yellyoshua/whatsapp-chat-parser/constants"
	"github.com/yellyoshua/whatsapp-chat-parser/logger"
	"github.com/yellyoshua/whatsapp-chat-parser/paper"
	"github.com/yellyoshua/whatsapp-chat-parser/storage"
	"github.com/yellyoshua/whatsapp-chat-parser/utils"
	"github.com/yellyoshua/whatsapp-chat-parser/whatsapp"
)

func initParsers() (whatsapp.Parser, paper.Writer) {
	var chat whatsapp.Parser = whatsapp.New()
	var writer = paper.New()
	return chat, writer
}

func HolyShit(ctx *gin.Context) {
	defer ctx.Done()
	ctx.String(200, "Holy shit!")
}

func PostUploadChatFiles(clientStorage storage.Uploader, attachmentURI string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		responseFormat, _ := ctx.Params.Get("format")
		chat, writer := initParsers()

		var text string
		var fullMessages string

		uuid := utils.NewUniqueID()
		ctx.Request.ParseMultipartForm(10 << 20)
		file, header, _ := ctx.Request.FormFile("file")

		filesInZip, err := utils.ExtractZipFile(file, header.Size, uuid)
		defer file.Close()

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, "error reading zip file")
			ctx.Done()
		} else {
			var wg sync.WaitGroup
			var filesToScanText map[string]io.Reader = make(map[string]io.Reader)
			var filesToUpload map[string]io.Reader = make(map[string]io.Reader)

			if err := utils.DuplicateReaders(filesInZip, filesToUpload, filesToScanText); err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, "error duplicating readers")
				ctx.Done()
			} else {

				wg.Add(1)
				go extractChatFromFiles(filesToScanText, &text, &wg)
				wg.Add(1)
				go uploadFilesS3(filesToUpload, clientStorage, &wg)

				wg.Wait()
				if err := chat.ParserMessages([]byte(text), &fullMessages); err != nil {
					ctx.AbortWithStatusJSON(http.StatusInternalServerError, "error processing messages")
					ctx.Done()
				} else {
					book := writer.UnmarshalMessagesAndSort(fullMessages, filepath.Join(attachmentURI, uuid))
					resWithFormat(ctx, book, responseFormat)
				}
			}
		}
	}
}

func PostParseOnlyChat(ctx *gin.Context) {
	chat, writer := initParsers()
	file, fileHeader, _ := ctx.Request.FormFile("file")

	var text string
	var wg sync.WaitGroup
	var fullMessages string

	files := map[string]io.Reader{
		fileHeader.Filename: file,
	}

	wg.Add(1)
	go extractChatFromFiles(files, &text, &wg)

	wg.Wait()

	if err := chat.ParserMessages([]byte(text), &fullMessages); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, "error processing messages")
		ctx.Done()
	} else {
		book := writer.UnmarshalMessagesAndSort(fullMessages, "/")
		responseFormat, _ := ctx.Params.Get("format")
		resWithFormat(ctx, book, responseFormat)
	}
}

func resWithFormat(ctx *gin.Context, book paper.Book, responseFormat string) {
	format := utils.UniqLowerCase(responseFormat)
	if format == constants.FormatJSON {
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
					"uuid":     nil,
					"count":    messages.Count,
					"messages": result,
				})
				ctx.Done()
			}
		}
	} else if format == constants.FormatHTML {
		messages, err := book.ExportHTML(paper.Minimal)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			ctx.Done()
		} else {
			ctx.String(http.StatusOK, messages)
			ctx.Done()
		}

	} else {
		ctx.String(http.StatusOK, "OK")
		ctx.Done()
	}
}

func uploadFilesS3(files map[string]io.Reader, clientStorage storage.Uploader, wg *sync.WaitGroup) {
	defer wg.Done()

	if err := clientStorage.UploadFiles(files); err != nil {
		logger.Info("Error uploading files -> %s", err)
	}
}

func extractChatFromFiles(files map[string]io.Reader, messages *string, wg *sync.WaitGroup) {
	defer wg.Done()

	for fullPath, f := range files {
		if isTXT := strings.Contains(fullPath, ".txt"); isTXT {
			if err := utils.ValueOfTextFile(f, messages); err != nil {
				logger.Info("Error ValueOfTextFile -> %s", err)
			}
		}
	}
}
