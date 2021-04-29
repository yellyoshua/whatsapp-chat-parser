package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yellyoshua/whatsapp-chat-parser/constants"
	"github.com/yellyoshua/whatsapp-chat-parser/storage"
	"github.com/yellyoshua/whatsapp-chat-parser/utils"
)

func abortConnection(ctx *gin.Context) {
	conn, _, err := ctx.Writer.Hijack()

	if err != nil {
		ctx.Abort()
	} else {
		if err := conn.Close(); err != nil {
			ctx.Abort()
		}
	}
}

func InjectDependencies(clientStorage *storage.Uploader) gin.HandlerFunc {
	attachmentURL := utils.GetAttachmentURL()

	return func(ctx *gin.Context) {
		ctx.Set(constants.KEY_MIDDLEWARE_CLIENT_STORAGE, clientStorage)
		ctx.Set(constants.KEY_MIDDLEWARE_ATTACHMENT_URL, attachmentURL)
		ctx.Next()
	}
}

// ParseOnlyChat this check request submit .txt chat file is valid
func ParseOnlyChat(ctx *gin.Context) {
	ctx.Next()
	ctx.Request.ParseMultipartForm(10 << 20)
	_, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		abortConnection(ctx)
	} else {
		isTXTFile := strings.Contains(header.Filename, ".txt")

		if !isTXTFile {
			err := fmt.Errorf("submit a .txt file")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			abortConnection(ctx)
		} else {
			ctx.Next()
		}
	}
}

// ParseFullChatZIP this check request submit .zip full chat file is valid
func ParseFullChatZIP(ctx *gin.Context) {
	ctx.Next()
	ctx.Request.ParseMultipartForm(10 << 20)
	_, header, err := ctx.Request.FormFile("file")

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprintf("error with file %s", err.Error()))
		abortConnection(ctx)
	} else {
		if isZipFile := strings.Contains(header.Filename, ".zip"); !isZipFile {

			ctx.AbortWithStatusJSON(http.StatusInternalServerError, "required be a zip file")
			abortConnection(ctx)
		} else {
			ctx.Next()
		}
	}
}
