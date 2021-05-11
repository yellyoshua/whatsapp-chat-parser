package whatsapp

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/urakozz/go-emoji"
	"github.com/yellyoshua/whatsapp-chat-parser/constants"
	"github.com/yellyoshua/whatsapp-chat-parser/utils"
)

func splitChatMessages(chat string) []string {
	regexSplitMessage, _ := regexp.Compile("(?:\r\n|\r|\n)")
	plainMessages := regexSplitMessage.Split(chat, -1)
	return plainMessages
}

func formatAttachment(message string) string {
	regexAttachment, _ := regexp.Compile(RegexAttachment)
	regexTextAttachment, _ := regexp.Compile(RegexAttachmentText)

	attachment := regexTextAttachment.ReplaceAllString(regexAttachment.FindString(message), "${1}$2")

	attachmentBytes := []byte(attachment)

	if len(attachmentBytes) == 0 {
		return message
	}

	fileName := strings.TrimSpace(string(attachmentBytes[1:]))

	fileName = strings.ReplaceAll(fileName, " ", "%20")

	repl := fmt.Sprintf(": <attached: %s>", fileName)
	messageWithAttachment := regexAttachment.ReplaceAllString(message, repl)
	return messageWithAttachment
}

func parseMessageAttachment(message string) Attachment {
	var attachment Attachment
	regexMessageAttachment, _ := regexp.Compile(RegexMessageAttachment)

	if regexMessageAttachment.MatchString(message) {
		fileName := regexMessageAttachment.FindStringSubmatch(message)[1]
		attachment = Attachment{
			FileName: strings.TrimSpace(fileName),
		}
	}

	return attachment
}

func parsePlainMessages(plainMessages []string) []RawMessage {
	rawMessages := make([]RawMessage, 0)

	parserMessage, _ := regexp.Compile(RegexpMessage)
	parserMessageSystem, _ := regexp.Compile(RegexpMessageSystem)

	for _, plainMessage := range plainMessages {
		if existMessage := len(plainMessage) > 0; existMessage {
			message := formatAttachment(plainMessage)

			if !parserMessage.MatchString(plainMessage) {
				if parserMessageSystem.MatchString(plainMessage) {
					messageRaw := RawMessage{
						IsSystem: true,
						Message:  message,
					}
					rawMessages = append(rawMessages, messageRaw)
					continue
				}

				if len(rawMessages)-1 >= 0 {
					prevMessage := rawMessages[len(rawMessages)-1]

					messageRaw := RawMessage{
						IsSystem: prevMessage.IsSystem,
						Message:  prevMessage.Message + "\n" + message,
					}

					rawMessages[len(rawMessages)-1] = messageRaw
				}

			} else {
				messageRaw := RawMessage{
					IsSystem: false,
					Message:  message,
				}
				rawMessages = append(rawMessages, messageRaw)
			}
		}
	}

	return rawMessages
}

func rawMessagesToMessages(rawMessages []RawMessage) []Message {
	parserMessage, _ := regexp.Compile(RegexpMessage)
	parserMessageSystem, _ := regexp.Compile(RegexpMessageSystem)

	var sender string
	var receiver string
	var lastDate string
	var systemAuthor = "system"

	messages := make([]Message, 0)

	emojiConvert := emoji.NewEmojiParser()

	for _, m := range rawMessages {
		if m.IsSystem {
			// [_ , date, time, ampm, message, ...]
			messageGroup := utils.SafeStringArray(parserMessageSystem.FindStringSubmatch(m.Message), 5)

			mDate := messageGroup[constants.MESSAGE_INDEX_DATE]            // index 1
			mTime := messageGroup[constants.MESSAGE_INDEX_TIME]            // index 2
			mTimeAMPM := messageGroup[constants.MESSAGE_INDEX_TIME_FORMAT] // index 3
			mMessage := messageGroup[constants.MESSAGE_INDEX_MESSAGE-1]    // index 4

			messageValue := emojiConvert.ToHtmlEntities(mMessage)

			messageTime := formatDate(mDate, mTime, mTimeAMPM)

			messages = append(messages, Message{
				Date:    messageTime,
				Author:  systemAuthor,
				Message: messageValue,
				IsInfo:  true,
			})
		} else {
			// [_ ,date, time, ampm, author, message, ...]
			messageGroup := utils.SafeStringArray(parserMessage.FindStringSubmatch(m.Message), 6)

			mDate := messageGroup[constants.MESSAGE_INDEX_DATE]            // index 1
			mTime := messageGroup[constants.MESSAGE_INDEX_TIME]            // index 2
			mTimeAMPM := messageGroup[constants.MESSAGE_INDEX_TIME_FORMAT] // index 3
			mAuthor := messageGroup[constants.MESSAGE_INDEX_AUTHOR]        // index 4
			mMessage := messageGroup[constants.MESSAGE_INDEX_MESSAGE]      // index 5

			messageValue := emojiConvert.ToHtmlEntities(mMessage)

			messageTime := formatDate(mDate, mTime, mTimeAMPM)
			currentDate := fmt.Sprintf("%v_%v_%v", messageTime.Month, messageTime.Day, messageTime.Year)

			if len(lastDate) == 0 || !utils.IsEqualString(currentDate, lastDate) {
				badgeMessagesDate := Message{
					Date:       messageTime,
					Author:     systemAuthor,
					Message:    getTranslateDate("es", messageTime.Month, messageTime.Day, messageTime.Year),
					IsSender:   false,
					IsReceiver: false,
					IsInfo:     true,
				}

				messages = append(messages, badgeMessagesDate)
			}

			lastDate = currentDate

			if notBeDefined := len(sender) == 0; notBeDefined && mAuthor != receiver {
				sender = mAuthor
			}

			if notBeDefined := len(receiver) == 0; notBeDefined && mAuthor != sender {
				receiver = mAuthor
			}

			if sender == mAuthor {
				message := Message{
					Date:       messageTime,
					Author:     mAuthor,
					Message:    messageValue,
					Attachment: parseMessageAttachment(messageValue),
					IsInfo:     false,
					IsSender:   true,
					IsReceiver: false,
				}
				messages = append(messages, message)
				continue
			}

			if receiver == mAuthor {
				message := Message{
					Date:       messageTime,
					Author:     mAuthor,
					Message:    messageValue,
					Attachment: parseMessageAttachment(messageValue),
					IsInfo:     false,
					IsSender:   false,
					IsReceiver: true,
				}
				messages = append(messages, message)
				continue
			}
		}
	}

	return messages
}
