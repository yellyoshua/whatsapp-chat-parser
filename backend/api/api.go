package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// TempPath __
const TempPath = ".temp/"

// API __
type API interface {
	GET(path string, handler Handler)
	POST(path string, handler Handler)
	Listen(port string) error
}

// Handler __
type Handler func(http.ResponseWriter, *http.Request)

type apistruct struct {
	router *gin.Engine
}

// New __
func New() API {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	return &apistruct{router: router}
}

func (api *apistruct) Listen(port string) error {
	server := createServer(api.router, port)
	return server.ListenAndServe()
}

func (api *apistruct) GET(path string, handler Handler) {
	api.router.GET(path, gingonictohttp(handler))
}

func (api *apistruct) POST(path string, handler Handler) {
	api.router.POST(path, gingonictohttp(handler))
}

// ResponseBadRequest __
func ResponseBadRequest(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

func createServer(router *gin.Engine, port string) *http.Server {
	if noPort := len(port) == 0; noPort {
		port = "3000"
	}

	server := new(http.Server)
	server.Addr = ":" + port
	server.Handler = router
	server.ReadTimeout = 10 * time.Second
	server.WriteTimeout = 10 * time.Second
	server.MaxHeaderBytes = 1 << 20
	return server
}

func gingonictohttp(handler func(http.ResponseWriter, *http.Request)) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		handler(ctx.Writer, ctx.Request)
	}
}
