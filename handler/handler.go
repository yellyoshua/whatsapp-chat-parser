package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"strings"

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

			if err := processFiles(filesInZip, clientStorage, &text); err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, "error processing zip files")
				ctx.Done()
			} else {
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

	file, _, _ := ctx.Request.FormFile("file")

	var plainMessages string
	var fullMessages string

	if err := utils.ValueOfTextFile(file, &plainMessages); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		ctx.Done()
	} else {

		if err := chat.ParserMessages([]byte(plainMessages), &fullMessages); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, "error processing messages")
			ctx.Done()
		} else {
			book := writer.UnmarshalMessagesAndSort(fullMessages, "/")
			responseFormat, _ := ctx.Params.Get("format")
			resWithFormat(ctx, book, responseFormat)
		}
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

func processFiles(files map[string]io.Reader, clientStorage storage.Uploader, messages *string) error {
	var filesToUpload map[string]io.Reader = make(map[string]io.Reader)
	var filesToScanText map[string]io.Reader = make(map[string]io.Reader)

	for i, f := range files {
		var toUpload bytes.Buffer
		var toScanText bytes.Buffer
		err := utils.CopyReader(f, &toUpload, &toScanText)
		if err != nil {
			return err
		}
		filesToUpload[i] = &toUpload
		filesToScanText[i] = &toScanText
	}

	if err := clientStorage.UploadFiles(filesToUpload); err != nil {
		logger.Info("Error uploading files -> %s", err)
		return err
	}

	for fullPath, f := range filesToScanText {
		if isTXT := strings.Contains(fullPath, ".txt"); isTXT {
			if err := utils.ValueOfTextFile(f, messages); err != nil {
				logger.Info("Error ValueOfTextFile -> %s", err)
				return err
			}
		}
	}

	return nil
}
