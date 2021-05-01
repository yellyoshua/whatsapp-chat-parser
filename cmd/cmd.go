package main

import (
	"context"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/yellyoshua/whatsapp-chat-parser/constants"
	"github.com/yellyoshua/whatsapp-chat-parser/handler"
	"github.com/yellyoshua/whatsapp-chat-parser/logger"
	"github.com/yellyoshua/whatsapp-chat-parser/middleware"
	"github.com/yellyoshua/whatsapp-chat-parser/storage"
	"github.com/yellyoshua/whatsapp-chat-parser/utils"
)

func setupEnvironments() {
	godotenv.Load(".env.local")

	// ENV_NAME : isRequerid?
	envs := map[string]bool{
		"PORT":           false,
		"S3_BUCKET_NAME": true,
		"AWS_REGION":     true,
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

func checkChatParserCLI() {
	output, err := exec.Command(constants.CLI_WP_PARSER, "--is-ok").Output()
	if err != nil {
		logger.Fatal("error executing [%s] CLI command -> %s", constants.CLI_WP_PARSER, err.Error())
	} else {
		message := string(output[:])

		if noOkey := !utils.IsEqualString(message, "ok"); noOkey {
			logger.Fatal("[%s] CLI is not installed, output -> %s", constants.CLI_WP_PARSER, message)
		}
	}
}

func notExistFolder(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

func main() {
	checkChatParserCLI()
	setupFolders()
	setupEnvironments()

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	port := os.Getenv("PORT")
	clientStorage := storage.New()

	exit := make(chan bool)
	go handleInterrupt(exit)

	router.GET("/", handler.HolyShit)

	whatsappHandler := router.Group("/whatsapp").Use(middleware.InjectDependencies(&clientStorage))

	whatsappHandler.POST("/:format/chat", middleware.ParseFullChatZIP, handler.PostUploadChatFiles)
	whatsappHandler.POST("/:format/messages", middleware.ParseOnlyChat, handler.PostParseOnlyChat)

	if noPort := len(port) == 0; noPort {
		port = constants.DefaultPort
	}

	// HTTP SERVER
	httpServer := &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			panic("Error trying start whatsapp-chat-parser server -> " + err.Error())
		}
	}()

	<-exit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		panic("Gracefull HTTP server shutdown failed: " + err.Error())
	}

}

func handleInterrupt(exit chan bool) {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)
	<-ch

	logger.Info("Closing server")
	exit <- true
}
