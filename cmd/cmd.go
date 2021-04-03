package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/yellyoshua/whatsapp-chat-parser/api"
	"github.com/yellyoshua/whatsapp-chat-parser/constants"
	"github.com/yellyoshua/whatsapp-chat-parser/logger"
	"github.com/yellyoshua/whatsapp-chat-parser/paper"
	"github.com/yellyoshua/whatsapp-chat-parser/storage"
	"github.com/yellyoshua/whatsapp-chat-parser/utils"
	"github.com/yellyoshua/whatsapp-chat-parser/whatsapp"
)

func setupEnvironments() {
	godotenv.Load(".env.local")

	// ENV_NAME : isRequerid?
	envs := map[string]bool{
		"PORT":           false,
		"AWS_ACCESS_KEY": true,
		"AWS_SECRET_KEY": true,
	}

	for name, isRequired := range envs {
		if value := os.Getenv(name); len(value) == 0 && isRequired {
			logger.Fatal("%v environment is requerid", name)
		}
	}
}

func setupFolders() {
	folders := map[string]os.FileMode{}

	for folder, permission := range folders {
		if notExistFolder(folder) {
			if err := os.Mkdir(folder, permission); err != nil {
				logger.Fatal("Error creating folder %s -> %v", folder, err)
			}
		}
	}
}

func notExistFolder(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
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

func handlerUploadChats(clientStorage storage.Uploader, attachmentURI string) gin.HandlerFunc {
	// TODO: Post form or uri with the response type {JSON or HTML}
	var chat whatsapp.Parser = whatsapp.New()
	return func(ctx *gin.Context) {
		writer := paper.New()

		uuid := utils.NewUniqueID()
		ctx.Request.ParseMultipartForm(10 << 20)
		file, header, err := ctx.Request.FormFile("file")

		defer ctx.Done()

		if err != nil {
			ctx.AbortWithError(500, fmt.Errorf("error file uploading"))
		} else {

			if isZipFile := strings.Contains(header.Filename, ".zip"); !isZipFile {
				ctx.AbortWithError(500, fmt.Errorf("required be a zip file"))
			} else {

				filesInZip, err := utils.ExtractZipFile(file, header.Size, uuid)
				defer file.Close()
				if err != nil {
					ctx.AbortWithError(500, fmt.Errorf("error reading zip file"))
					return
				}

				var text string
				var fullMessages string

				if err := processFiles(filesInZip, clientStorage, &text); err != nil {
					ctx.AbortWithError(500, fmt.Errorf("error processing zip files"))
					return
				}

				if err := chat.ParserMessages([]byte(text), &fullMessages); err != nil {
					ctx.AbortWithError(500, fmt.Errorf("error processing messages"))
					return
				}

				book := writer.UnmarshalMessagesAndSort(fullMessages, filepath.Join(attachmentURI, uuid))
				messages, err := book.ExportJSON()
				if err != nil {
					ctx.AbortWithError(500, fmt.Errorf("error marshal messages"))
					return
				}

				var result interface{}

				if err := json.Unmarshal(messages.Value, &result); err != nil {
					ctx.AbortWithError(500, fmt.Errorf("error marshal messages"))
					return
				}

				var resJSON = map[string]interface{}{
					"uuid":     uuid,
					"count":    messages.Count,
					"messages": result,
				}

				ctx.JSON(http.StatusOK, resJSON)
			}
		}
	}
}

func handlerHolyShit(ctx *gin.Context) {
	defer ctx.Done()
	ctx.String(200, "Holy shit!")
}

func startApp() {
	var port string = os.Getenv("PORT")
	var router api.API = api.New()
	var clientStorage storage.Uploader = storage.New()
	var attachmentURI string = constants.S3BucketEndpoint

	router.POST("/upload", handlerUploadChats(clientStorage, attachmentURI))
	router.GET("/", handlerHolyShit)
	logger.CheckError("Error listen server", router.Listen(port))
}

func main() {
	setupFolders()
	setupEnvironments()
	startApp()
}
