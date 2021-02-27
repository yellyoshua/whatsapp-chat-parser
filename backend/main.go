package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/yellyoshua/whatsapp-chat-parser/api"
	"github.com/yellyoshua/whatsapp-chat-parser/logger"
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

	router.POST("/upload", handlerUploadChats(clientStorage))

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

func handlerUploadChats(clientStorage storage.Uploader) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) {

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

		reader, err := zip.NewReader(file, header.Size)
		if err != nil {
			api.ResponseBadRequest(w, "Error reading zip file")
			return
		}

		var filesToUpload map[string]io.Reader = make(map[string]io.Reader)
		var filesToScanText map[string]io.Reader = make(map[string]io.Reader)

		for _, zipFile := range reader.File {

			// TODO: Solve copy file buffer

			// var toUpload io.Reader
			// var toScan io.Reader
			// var err error

			f, _ := zipFile.Open()
			fullPath := filepath.Join(uuid, zipFile.Name)

			// b, _ := ioutil.ReadAll(f)
			var a, c bytes.Buffer
			w := io.MultiWriter(&a, &c)

			if _, err := io.Copy(w, f); err != nil {
				continue
			}

			// toScan = bytes.NewBuffer(b)
			// toUpload = bytes.NewBuffer(b)

			// toUpload, err = utils.CopyFile(f)
			// if err != nil {
			// 	logger.Info("Error toUpload -> %s", err)
			// 	continue
			// }
			// toScan, err = utils.CopyFile(f)
			// if err != nil {
			// 	logger.Info("Error toScan -> %s", err)
			// 	continue
			// }

			// filesToUpload[fullPath] = toUpload
			// filesToScanText[fullPath] = toScan
			filesToUpload[fullPath] = io.Reader(&a)
			filesToScanText[fullPath] = io.Reader(&c)

			if err := f.Close(); err != nil {
				logger.Info("Error closing -> %s", err)
				continue
			}
		}

		// if err := clientStorage.UploadFiles(filesToUpload); err != nil {
		// 	logger.Info("Error uploading files -> %s", err)
		// 	api.ResponseBadRequest(w, "Error uploading files")
		// 	return
		// }

		var fullMessages string

		for fullPath, f := range filesToUpload {

			if isTXT := strings.Contains(fullPath, ".txt"); isTXT {
				if err := utils.ValueOfTextFile(f, &fullMessages); err != nil {
					logger.Info("Error ValueOfTextFile -> %s", err)
					continue
				}
			}
		}

		logger.Info("FM (%v)", len(fullMessages))

		var fullMessages2 string

		for fullPath, f := range filesToScanText {

			if isTXT := strings.Contains(fullPath, ".txt"); isTXT {
				if err := utils.ValueOfTextFile(f, &fullMessages); err != nil {
					logger.Info("Error ValueOfTextFile -> %s", err)
					continue
				}
			}
		}

		logger.Info("FM2 (%v)", len(fullMessages2))

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Name of the File: ", header.Filename)
		fmt.Fprintln(w, "Size of the File: ", header.Size)
		fmt.Fprintln(w, "Uploaded in folder: ", uuid)
		return
	}
}
