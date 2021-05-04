package whatsapp

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/yellyoshua/whatsapp-chat-parser/constants"
)

// RegexContact format input `$date - Carlos perez: $message`
var RegexContact string = `(\d{1,2}/\d{1,2}/\d{2,4})+(, )[0-9:]+(.+?)(: )`

// RegexAttachment format input `$date - $contact: IMG-20200319-WA0011.jpg (file attached)`
var RegexAttachment string = `(: )+[\S\s]+(\.\w{2,4}\s)+\(+(file attached)+\)`

// RegexTextAttachment format input `IMG-20200319-WA0011.jpg (file attached)`
var RegexTextAttachment string = `\(file attached\)`

// RawBuilder _
type RawBuilder interface {
	ChatParser(user_id string, chat []byte) (Parser, error)
}

type rawbuilder struct{}

// Parser _
type Parser interface {
	ParserMessages(outputMessages *string) error
}

type parserstruct struct {
	chat *string
}

func New() RawBuilder {
	return &rawbuilder{}
}

func getTempChat(user_id string) string {
	pwd, _ := os.Getwd()

	if len(user_id) > 0 {
		return filepath.Join(pwd, ".tmp", user_id)
	} else {
		return filepath.Join(pwd, ".tmp")
	}
}

func (r *rawbuilder) ChatParser(user_id string, chat []byte) (Parser, error) {
	var chatValue string
	chatTemp := getTempChat(user_id)
	tempPath := getTempChat("")

	err := os.MkdirAll(tempPath, os.ModeDir)
	if err != nil {
		return &parserstruct{chat: &chatValue}, err
	}

	f, errInitTempFile := os.Create(chatTemp)
	if errInitTempFile != nil {
		return &parserstruct{chat: &chatValue}, errInitTempFile
	}

	if errCHMod := f.Chmod(0777); errCHMod != nil {
		return &parserstruct{chat: &chatValue}, errCHMod
	}

	_, errorWrite := f.Write([]byte(byteToStringMessages(chat)))
	if errorWrite != nil {
		return &parserstruct{chat: &chatValue}, errorWrite
	}

	errCloseTempChat := f.Close()
	if errCloseTempChat != nil {
		return &parserstruct{chat: &chatValue}, errCloseTempChat
	}

	parsedChat, errChatParser := exec.Command(constants.CLI_WP_PARSER, chatTemp).Output()
	if errChatParser != nil {
		return &parserstruct{chat: &chatValue}, errChatParser
	}

	chatValue = string(parsedChat[:])

	return &parserstruct{chat: &chatValue}, nil
}

// ParserMessages _
func (p *parserstruct) ParserMessages(outputMessages *string) error {
	*outputMessages = *p.chat
	return nil
}

func byteToStringMessages(data []byte) string {
	var messages string

	regexContact, _ := regexp.Compile(RegexContact)

	var whatsappMessages []string

	plainMessages := strings.TrimSpace(string(data))
	bytesOfMessages := []byte(plainMessages)
	messagesIndexes := regexContact.FindAllStringIndex(plainMessages, -1)

	for i := 0; i < len(messagesIndexes); i++ {
		axis := messagesIndexes[i]
		nextIndex := i + 1
		existMessage := len(axis) == 2

		if existMessage {
			start, _ := axis[0], axis[1]

			if nextIndex < len(messagesIndexes) {
				nextAxis := messagesIndexes[nextIndex]
				existNextAxis := len(nextAxis) == 2

				if existNextAxis {
					message := string(bytesOfMessages[start:nextAxis[0]])
					message = strings.TrimSpace(replaceAttachment(message))
					whatsappMessages = append(whatsappMessages, message)
				}

			} else {
				message := string(bytesOfMessages[start:])
				message = strings.TrimSpace(replaceAttachment(message))
				whatsappMessages = append(whatsappMessages, message)
			}

		}
	}

	messages = strings.Join(whatsappMessages, "\n")

	return messages
}

func replaceAttachment(message string) string {
	regexAttachment, _ := regexp.Compile(RegexAttachment)
	regexTextAttachment, _ := regexp.Compile(RegexTextAttachment)

	attachment := regexTextAttachment.ReplaceAllString(regexAttachment.FindString(message), "${1}$2")

	attachmentBytes := []byte(attachment)

	if len(attachmentBytes) == 0 {
		return message
	}

	fileName := strings.TrimSpace(string(attachmentBytes[1:]))

	fileName = strings.ReplaceAll(fileName, " ", "%20")

	repl := fmt.Sprintf(": <attached: %s>", fileName)
	result := regexAttachment.ReplaceAllString(message, repl)
	return result
}
