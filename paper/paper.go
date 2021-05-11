package paper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"io"
	"io/ioutil"
	"path"
	"path/filepath"
	"time"

	"github.com/yellyoshua/whatsapp-chat-parser/logger"
	"github.com/yellyoshua/whatsapp-chat-parser/whatsapp"
)

var templateFuncs = template.FuncMap{
	"isAttachmentAudio": func(file_extension string) bool {
		fe := file_extension
		if fe == ".opus" || fe == ".mp3" {
			return true
		}
		return false
	},
	"isAttachmentImage": func(file_extension string) bool {
		fe := file_extension
		if fe == ".webp" || fe == ".gift" || fe == ".jpg" || fe == ".png" {
			return true
		}
		return false
	},
	"isAttachmentVideo": func(file_extension string) bool {
		fe := file_extension
		return fe == ".mp4"
	},
	"formatDate": func(date whatsapp.DateFormat) string {
		if len(date.Format) > 0 {
			return fmt.Sprintf("%s:%s %s", date.Hours, date.Mins, date.Format)
		}
		return fmt.Sprintf("%s:%s", date.Hours, date.Mins)
	},
}

// Type _
type Type func(messages []whatsapp.Message) *BookData

// Attachment _
type Attachment struct {
	Exist     bool   `json:"exist"`
	FileName  string `json:"fileName,omitempty"`
	Extension string `json:"extension,omitempty"`
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
	Export() []whatsapp.Message
	ExportJSON() (MessagesJSON, error)
	ExportHTML(paper Type) (string, error)
	ExportHTMLFile(paper Type, filePathName string) error
}

type export struct {
	messages []whatsapp.Message
}

type MessagesJSON struct {
	Value []byte
	Count int
}

// Writer _
type Writer interface {
	AttachFiles(attachmentFiles map[string]string, attachmentURL string) Book
}

type writertruct struct {
	messages []whatsapp.Message
}

// BookData _
type BookData struct {
	Messages   []whatsapp.Message `json:"messages"`
	Background string             `json:"background"`
}

// Loves is a template with hearts background
func Loves(messages []whatsapp.Message) *BookData {
	return &BookData{
		Messages:   messages,
		Background: "green",
	}
}

// Friends is a template dedicated for friends
func Friends(messages []whatsapp.Message) *BookData {
	return &BookData{
		Messages:   messages,
		Background: "red",
	}
}

// Minimal is a template dedicated for everyone
func Minimal(messages []whatsapp.Message) *BookData {
	return &BookData{
		Messages:   messages,
		Background: "orange",
	}
}

// New __
func New(messages []whatsapp.Message) Writer {
	return &writertruct{messages: messages}
}

func attachURLFile(attachmentURL string, attachmentFiles map[string]string, attachment whatsapp.Attachment) whatsapp.Attachment {
	if len(attachmentFiles[attachment.FileName]) != 0 {
		return whatsapp.Attachment{
			Exist:     true,
			Extension: filepath.Ext(attachment.FileName),
			FileName:  path.Join(attachmentURL, attachmentFiles[attachment.FileName]),
		}
	}

	return whatsapp.Attachment{
		Exist:     false,
		Extension: filepath.Ext(attachment.FileName),
		FileName:  attachment.FileName,
	}
}

func (p *writertruct) AttachFiles(attachmentFiles map[string]string, attachmentURL string) Book {
	var newMessages = make([]whatsapp.Message, 0)

	for _, message := range p.messages {

		var attachment whatsapp.Attachment
		if len(message.Attachment.FileName) > 0 {
			attachment = attachURLFile(attachmentURL, attachmentFiles, message.Attachment)
		}
		newMessages = append(newMessages, whatsapp.Message{
			Attachment: attachment,
			Date:       message.Date,
			Author:     message.Author,
			IsSender:   message.IsSender,
			IsReceiver: message.IsReceiver,
			IsInfo:     message.IsInfo,
			Message:    message.Message,
		})
	}

	return &export{
		messages: newMessages,
	}
}

func parseHTMLEntities(scapedHTML string) string {
	return html.UnescapeString(scapedHTML)
}

func renderTemplate(bookTemplate string, data BookData, funcs template.FuncMap, buffer io.Writer) error {
	tmpl, err := template.New("book-rendered").Funcs(funcs).Parse(bookTemplate)
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
func (e *export) Export() []whatsapp.Message {
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

	if err := renderTemplate(bookTemplate, bookData, templateFuncs, buffer); err != nil {
		return book, err
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
	styles := bookStyle()

	return bookTemplate(styles)
}
