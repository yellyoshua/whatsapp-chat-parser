package paper

import (
	"bytes"
	"encoding/json"
	"html"
	"html/template"
	"io"
	"io/ioutil"
	"path/filepath"
	"sort"
	"time"

	"github.com/urakozz/go-emoji"
	"github.com/yellyoshua/whatsapp-chat-parser/logger"
)

// Type _
type Type func(messages []Message) *BookData

// Loves is a template with hearts background
func Loves(messages []Message) *BookData {
	return &BookData{
		Messages:   messages,
		Background: "green",
	}
}

// Friends is a template dedicated for friends
func Friends(messages []Message) *BookData {
	return &BookData{
		Messages:   messages,
		Background: "red",
	}
}

// Minimal is a template dedicated for everyone
func Minimal(messages []Message) *BookData {
	return &BookData{
		Messages:   messages,
		Background: "orange",
	}
}

// Attachment _
type Attachment struct {
	FileName string `json:"fileName,omitempty"`
}

// Message _
type Message struct {
	Date       string     `json:"date"`
	Author     string     `json:"author"`
	IsSender   bool       `json:"isSender"`
	IsInfo     bool       `json:"isInfo"`
	IsReceiver bool       `json:"isReceiver"`
	Message    string     `json:"message"`
	Attachment Attachment `json:"attachment"`
}

// Book __
type Book interface {
	Export() []Message
	ExportJSON() (MessagesJSON, error)
	ExportHTML(paper Type) (string, error)
	ExportHTMLFile(paper Type, filePathName string) error
}

type export struct {
	messages []Message
}

type MessagesJSON struct {
	Value []byte
	Count int
}

// Writer _
type Writer interface {
	UnmarshalMessagesAndSort(plainMessages string, attachmentFiles map[string]string, attachmentURL string) Book
}

type writertruct struct{}

// BookData _
type BookData struct {
	Messages   []Message `json:"messages"`
	Background string    `json:"background"`
}

// New __
func New() Writer {
	return &writertruct{}
}

func attachURLFile(attachmentURL string, attachmentFiles map[string]string, attachment Attachment) Attachment {
	if len(attachmentFiles[attachment.FileName]) != 0 {
		return Attachment{
			FileName: filepath.Join(attachmentURL, attachmentFiles[attachment.FileName]),
		}
	}

	return Attachment{
		FileName: filepath.Join(attachmentURL, attachment.FileName),
	}
}

func (p *writertruct) UnmarshalMessagesAndSort(plainMessages string, attachmentFiles map[string]string, attachmentURL string) Book {

	var temporalMessages []Message
	var messages []Message

	json.Unmarshal([]byte(plainMessages), &temporalMessages)

	var sender string = ""
	var receiver string = ""

	emojiConvert := emoji.NewEmojiParser()

	for _, m := range temporalMessages {
		messageValue := emojiConvert.ToHtmlEntities(m.Message)

		if notBeDefined := len(sender) == 0; notBeDefined && m.Author != receiver {
			sender = m.Author
		}

		if notBeDefined := len(receiver) == 0; notBeDefined && m.Author != sender {
			receiver = m.Author
		}

		if sender == m.Author {
			var attachment Attachment
			if len(m.Attachment.FileName) > 0 {
				attachment = attachURLFile(attachmentURL, attachmentFiles, m.Attachment)
			}

			currentMessage := Message{
				Date:       m.Date,
				Author:     m.Author,
				Message:    messageValue,
				Attachment: attachment,
				IsSender:   true,
				IsReceiver: false,
				IsInfo:     false,
			}
			messages = append(messages, currentMessage)
			continue
		}

		if receiver == m.Author {
			var attachment Attachment
			if len(m.Attachment.FileName) > 0 {
				attachment = attachURLFile(attachmentURL, attachmentFiles, m.Attachment)
			}

			currentMessage := Message{
				Date:       m.Date,
				Author:     m.Author,
				Message:    messageValue,
				Attachment: attachment,
				IsSender:   false,
				IsReceiver: true,
				IsInfo:     false,
			}
			messages = append(messages, currentMessage)
			continue
		}

		currentMessage := Message{
			Date:       m.Date,
			Author:     "Info",
			Message:    messageValue,
			Attachment: Attachment{},
			IsSender:   false,
			IsReceiver: false,
			IsInfo:     true,
		}
		messages = append(messages, currentMessage)
	}

	sort.SliceStable(messages, func(i, j int) bool {
		return messages[i].Date < messages[j].Date
	})

	return &export{
		messages: messages,
	}
}

func parseHTMLEntities(scapedHTML string) string {
	return html.UnescapeString(scapedHTML)
}

func renderTemplate(bookTemplate string, data BookData, buffer io.Writer) error {
	tmpl, err := template.New("book-rendered").Parse(bookTemplate)
	if err != nil {
		return err
	}

	return tmpl.Execute(buffer, data)
}

// ExportHTMLFile _
func (e *export) ExportHTMLFile(paper Type, filePathName string) error {
	paperProps := *paper(e.messages)

	book, err := paintPaper(paperProps)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePathName, []byte(book), 0644)
}

// ExportJSON
func (e *export) ExportJSON() (MessagesJSON, error) {
	messages, err := json.Marshal(e.messages)
	return MessagesJSON{
		Value: messages,
		Count: len(e.messages),
	}, err
}

// Export
func (e *export) Export() []Message {
	return e.messages
}

// ExportHTML __
func (e *export) ExportHTML(paper Type) (string, error) {
	paperProps := *paper(e.messages)
	return paintPaper(paperProps)
}

func paintPaper(bookData BookData) (string, error) {
	begin := time.Now()
	var book string

	var buffer = new(bytes.Buffer)
	bookTemplate := createBookTemplate()
	if err := renderTemplate(bookTemplate, bookData, buffer); err != nil {
		return book, nil
	}

	defer func() {
		end := time.Now()
		logger.Info("time: %v", end.Sub(begin))
	}()

	book = parseHTMLEntities(buffer.String())

	// should exporter pdf file here filePathName
	return book, nil
}

func createBookTemplate() string {
	var bookTemplate string

	bookStyles := `
  <style>
    :root {
      --sender-bg-color: white;
      --receiver-bg-color: rgba(39, 255, 118, 0.611);
    }

    * {
      box-sizing: border-box;
      padding: 0;
      margin: 0;
      font-family: 'SymbolaRegular';
      font-weight: normal;
      font-style: normal;
    }

    @page {
      size: auto;
      /* auto is the initial value */

      /* this affects the margin in the printer settings */
      margin: 5mm 5mm 5mm 5mm;
    }

    body {
      /* display: flex;
      flex-direction: column;
      justify-content: space-between;
      flex-flow: wrap;
      flex-wrap: wrap; */
      /* height: 100vh; */
      height: fit-content;
      width: 100%;
      background: pink;
      -webkit-columns: 2;
      -moz-columns: 2;
      columns: 2;
      column-rule: 3px solid lightblue;
      column-rule-style: dotted;
      column-count: 2;
      column-gap: 10px;
      column-span: all;
      column-fill: auto;
    }

    .receiver-message-container {
      margin-right: 5%;
      justify-content: flex-end;
    }

    .sender-message-container {
      margin-left: 5%;
      justify-content: flex-start;
    }

    .message-container {

      display: flex;
      width: 95%;
      height: auto;
      padding: 5px 0px;
    }

    .sender-mini-box {
      left: -3px;
      background: linear-gradient(45deg, var(--sender-bg-color) 50%, transparent 50%);
      transform: rotate(45deg);
    }

    .receiver-mini-box {
      right: -3px;
      background: linear-gradient(45deg, var(--receiver-bg-color) 50%, transparent 50%);
      transform: rotate(225deg);
    }

    .mini-box {
      position: absolute;
      top: 5px;
      border-radius: 1px;
      width: 20px;
      height: 20px;
    }

    .receiver-message-bubble {
      background: var(--receiver-bg-color);
    }

    .sender-message-bubble {
      background: var(--sender-bg-color);
    }

    .message-bubble {
      position: relative;
      width: 75%;
      padding: 10px;
      border-radius: 5px;
    }

    .message-author {
      font-size: 14px;
    }

    .message-message {
      font-size: 16px;
      text-align: justify;
    }

    img.emoji {
      width: 19px;
      height: 19px;
    }

    .message-date {
      font-size: 12px;
    }

    @media print {

      .pg-break {
        clear: both;
        /* page-break-after: always; */
        -webkit-column-break-inside: avoid;
        /* Chrome, Safari, Opera */
        page-break-inside: avoid;
        /* Firefox */
        break-inside: avoid;
        /* IE 10+ */
        page-break-before: avoid;
      }
    }
  </style>
  `

	cardIfBeSender := `
  {{ if .IsSender -}}
			<div class="pg-break message-container sender-message-container">
				<div class="message-bubble sender-message-bubble">
					<div class="sender-mini-box mini-box"></div>
					<div class="message-author"><strong>{{.Author}} (sender)</strong></div>
					<div class="message-message">{{.Message}}</div>
					<div class="message-date"><strong>{{.Date}}</strong></div>
				</div>
			</div>
	{{- end}}
  `

	cardIfBeReceiver := `
  {{ if .IsReceiver -}}
  <div class="pg-break message-container receiver-message-container">
    <div class="message-bubble receiver-message-bubble">
      <div class="receiver-mini-box mini-box"></div>
      <div class="message-author"><strong>{{.Author}} (receiver)</strong></div>
      <div class="message-message">{{.Message}}</div>
      <div class="message-date"><strong>{{.Date}}</strong></div>
    </div>
  </div>
  {{- end}}
  `

	mappingMessages := `
  {{- range .Messages }}
  ` + cardIfBeSender + `
  ` + cardIfBeReceiver + `
	{{- end}}
  `

	baseHTMLTemplate := `
  <!DOCTYPE html>
  <html lang="en">
  <head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Document</title>
  </head>
  <body>
  ` + mappingMessages + `
  ` + bookStyles + `
  </body>
  </html>
  `

	bookTemplate = baseHTMLTemplate

	return bookTemplate
}
