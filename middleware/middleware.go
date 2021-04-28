package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/yellyoshua/whatsapp-chat-parser/constants"
	"github.com/yellyoshua/whatsapp-chat-parser/storage"
	"github.com/yellyoshua/whatsapp-chat-parser/utils"
)

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
	// ctx.Request.ParseMultipartForm(10 << 20)
	// _, header, err := ctx.Request.FormFile("file")
	// if err != nil {
	// 	ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
	// 	ctx.Done()
	// } else {
	// 	isTXTFile := strings.Contains(header.Filename, ".txt")

	// 	if !isTXTFile {
	// 		err := fmt.Errorf("submit a .txt file")
	// 		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
	// 		ctx.Done()
	// 	} else {
	// 		ctx.Next()
	// 	}
	// }
}

// ParseFullChatZIP this check request submit .zip full chat file is valid
func ParseFullChatZIP(ctx *gin.Context) {
	ctx.Next()
	// ctx.Request.ParseMultipartForm(10 << 20)
	// _, header, err := ctx.Request.FormFile("file")

	// if err != nil {
	// 	ctx.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprintf("error with file %s", err.Error()))
	// 	ctx.Done()
	// } else {
	// 	if isZipFile := strings.Contains(header.Filename, ".zip"); !isZipFile {

	// 		ctx.AbortWithStatusJSON(http.StatusInternalServerError, "required be a zip file")
	// 		ctx.Done()
	// 	} else {
	// 		ctx.Next()
	// 	}
	// }
}
