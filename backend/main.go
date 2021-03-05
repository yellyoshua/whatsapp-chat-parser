package main

import (
	"archive/zip"
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/yellyoshua/whatsapp-chat-parser/api"
	"github.com/yellyoshua/whatsapp-chat-parser/chatparser"
	"github.com/yellyoshua/whatsapp-chat-parser/logger"
	"github.com/yellyoshua/whatsapp-chat-parser/paper"
	"github.com/yellyoshua/whatsapp-chat-parser/storage"
	"github.com/yellyoshua/whatsapp-chat-parser/utils"
)

// "github.com/kyokomi/emoji/v2"
// github.com/boombuler/barcode

/** TODO: Parse images with Regex.
 * Convert all to go routines
 *
 */

// Input message format: `3/12/20, 21:37 - Pepe: Example message`

// RegexContact format `$date - Carlos perez: $message`
var RegexContact string = `(\d{1,2}/\d{1,2}/\d{2,4})+(, )[0-9:]+(.+?)(: )`

func setupEnvironments() {
	logger.CheckError("Error loading .env file", godotenv.Load(".env"))
	// ENV_NAME : isRequerid?
	envs := map[string]bool{
		// App expose port default is "4000"
		"PORT": false,
		// Storage bucket name `myproject-production`
		"GCS_BUCKET": true,
		// Google Cloud Project Id `demoproject`
		"GCS_PROJECT_ID": true,
		// Especific IAM service account `example@gserviceaccount.com`
		"GCS_IAM_SVC": true,
	}

	for name, isRequired := range envs {
		if value := os.Getenv(name); len(value) == 0 && isRequired {
			logger.Fatal("%v environment is requerid", name)
		}
	}
}

func setupFolders() {
	folders := map[string]os.FileMode{
		api.TempPath: os.ModePerm,
	}

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

func main() {
	setupFolders()
	setupEnvironments()

	var port string = "4000"
	var router api.API = api.New()
	var clientStorage storage.Uploader = storage.New()
	var attachmentURI string = filepath.Join("https://storage.googleapis.com", os.Getenv("GCS_BUCKET"))

	router.POST("/upload", handlerUploadChats(clientStorage, attachmentURI))

	// var filePath string = "samples/Jhonny/chat.txt"
	// data, err := ioutil.ReadFile(filePath)
	// if err != nil {
	// 	logger.Fatal("Error reading file -> %s", err)
	// }

	// chat := chatparser.New()

	// var plainMessages string

	// logger.CheckError("Error parsing messages", chat.ParserMessages(data, &plainMessages))

	// writer := paper.New()
	// book := writer.UnmarshalMessagesAndSort(plainMessages)

	// pwd := utils.GetCurrentPath()

	// logger.CheckError("Error exporting html", book.ExportHTMLFile(paper.Loves, pwd+"/demo2.html"))

	// router.GET("/sample1", func(w http.ResponseWriter, r *http.Request) {
	// 	data, err := book.ExportJSON()
	// 	if err != nil {
	// 		api.ResponseBadRequest(w, "Error parsing json")
	// 		return
	// 	}
	// 	w.Write(data)
	// })

	logger.CheckError("Error listen server", router.Listen(port))
}

func handlerUploadChats(clientStorage storage.Uploader, attachmentURI string) api.Handler {
	var chat chatparser.Parser = chatparser.New()
	return func(w http.ResponseWriter, r *http.Request) {
		writer := paper.New()

		uuid := utils.NewUniqueID()
		r.ParseMultipartForm(10 << 20)
		file, header, err := r.FormFile("file")

		if err != nil {
			api.ResponseBadRequest(w, "Error file uploading")
			return
		}

		defer file.Close()

		if isZipFile := strings.Contains(header.Filename, ".zip"); isZipFile == false {
			api.ResponseBadRequest(w, "Required be a zip file")
			return
		}

		filesInZip, err := extractZipFile(file, header.Size, uuid)
		if err != nil {
			api.ResponseBadRequest(w, "Error reading zip file")
			return
		}

		var text string
		var fullMessages string

		if err := processFiles(filesInZip, clientStorage, &text); err != nil {
			api.ResponseBadRequest(w, "Error processing zip files")
			return
		}

		if err := chat.ParserMessages([]byte(text), &fullMessages); err != nil {
			api.ResponseBadRequest(w, "Error processing messages")
			return
		}

		book := writer.UnmarshalMessagesAndSort(fullMessages, filepath.Join(attachmentURI, uuid))
		messages, err := book.ExportJSON()
		if err != nil {
			api.ResponseBadRequest(w, "Error marshal messages")
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(messages)
		// fmt.Fprintln(w, "Name of the File: ", header.Filename)
		// fmt.Fprintln(w, "Size of the File: ", header.Size)
		// fmt.Fprintln(w, "Uploaded in folder: ", uuid)
		return
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

func extractZipFile(f multipart.File, size int64, folder string) (map[string]io.Reader, error) {
	var dst map[string]io.Reader = make(map[string]io.Reader)
	reader, err := zip.NewReader(f, size)
	if err != nil {
		return dst, err
	}

	for _, zipFile := range reader.File {
		var customReader bytes.Buffer
		f, _ := zipFile.Open()
		fullPath := filepath.Join(folder, zipFile.Name)

		if err := utils.CopyReader(f, &customReader); err != nil {
			return dst, err
		}

		dst[fullPath] = &customReader

		if err := f.Close(); err != nil {
			logger.Info("Error closing -> %s", err)
			return dst, err
		}
	}

	return dst, nil
}
