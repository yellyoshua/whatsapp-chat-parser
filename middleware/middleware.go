package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// MiddlewareParseOnlyChat this check request submit .txt chat file is valid
func MiddlewareParseOnlyChat(ctx *gin.Context) {
	ctx.Request.ParseMultipartForm(10 << 20)
	_, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		ctx.Done()
	} else {
		isTXTFile := strings.Contains(header.Filename, ".txt")

		if !isTXTFile {
			err := fmt.Errorf("submit a .txt file")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			ctx.Done()
		} else {
			ctx.Next()
		}
	}
}

// MiddlewareParseFullChatZIP this check request submit .zip full chat file is valid
func MiddlewareParseFullChatZIP(ctx *gin.Context) {
	ctx.Request.ParseMultipartForm(10 << 20)
	_, header, err := ctx.Request.FormFile("file")

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprintf("error with file %s", err.Error()))
		ctx.Done()
	} else {
		if isZipFile := strings.Contains(header.Filename, ".zip"); !isZipFile {

			ctx.AbortWithStatusJSON(http.StatusInternalServerError, "required be a zip file")
			ctx.Done()
		} else {
			ctx.Next()
		}
	}
}
