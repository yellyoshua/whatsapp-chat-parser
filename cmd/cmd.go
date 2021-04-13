package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/yellyoshua/whatsapp-chat-parser/api"
	"github.com/yellyoshua/whatsapp-chat-parser/constants"
	"github.com/yellyoshua/whatsapp-chat-parser/handler"
	"github.com/yellyoshua/whatsapp-chat-parser/logger"
	"github.com/yellyoshua/whatsapp-chat-parser/middleware"
	"github.com/yellyoshua/whatsapp-chat-parser/storage"
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

func handlerHolyShit(ctx *gin.Context) {
	defer ctx.Done()
	ctx.String(200, "Holy shit!")
}

func startApp() {
	var port string = os.Getenv("PORT")
	var router api.API = api.New()
	var clientStorage storage.Uploader = storage.New()
	var attachmentURI string = constants.S3BucketEndpoint

	router.GET("/", handlerHolyShit)

	router.Use(middleware.MiddlewareParseFullChatZIP).
		POST("/whatsapp/:format/chat", handler.PostUploadChatFiles(clientStorage, attachmentURI))

	router.Use(middleware.MiddlewareParseOnlyChat).
		POST("/whatsapp/:format/messages", handler.PostParseOnlyChat)

	logger.CheckError("Error listen server", router.Listen(port))
}

func main() {
	setupFolders()
	setupEnvironments()
	startApp()
}
