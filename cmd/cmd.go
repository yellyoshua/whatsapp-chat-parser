package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/yellyoshua/whatsapp-chat-parser/constants"
	"github.com/yellyoshua/whatsapp-chat-parser/handler"
	"github.com/yellyoshua/whatsapp-chat-parser/logger"
	"github.com/yellyoshua/whatsapp-chat-parser/middleware"
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

func notExistFolder(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

func main() {
	setupFolders()
	setupEnvironments()

	h := handler.New()
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	port := os.Getenv("PORT")

	exit := make(chan bool)
	go handleInterrupt(exit)

	router.GET("/", h.HolyShit)

	whatsappHandler := router.Group("/whatsapp").Use(middleware.InjectDependencies())

	whatsappHandler.POST("/:format/chat", middleware.ParseFullChatZIP, h.PostUploadChatFiles)
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
	// the channel used with signal.Notify should be buffered (SA1017)
	signal.Notify(ch, os.Interrupt)
	<-ch

	logger.Info("Closing server")
	exit <- true
}
