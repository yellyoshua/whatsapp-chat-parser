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
	"strings"
	"time"

	"github.com/urakozz/go-emoji"
	"github.com/yellyoshua/whatsapp-chat-parser/logger"
	"github.com/yellyoshua/whatsapp-chat-parser/utils"
)

var translations = map[string]interface{}{
	"es": map[string]interface{}{
		"months": map[string]string{
			"01": "Enero",
			"02": "Febrero",
			"03": "Marzo",
			"04": "Abril",
			"05": "Mayo",
			"06": "Junio",
			"07": "Julio",
			"08": "Agosto",
			"09": "Septiembre",
			"10": "Octubre",
			"11": "Noviembre",
			"12": "Diciembre",
		},
		"date_template": "%s %s, %s",
	},
}

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
	UnmarshalJSONMessages(json_messages string, attachmentFiles map[string]string, attachmentURL string) Book
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
			Exist:     true,
			Extension: filepath.Ext(attachment.FileName),
			FileName:  path.Join(attachmentURL, attachmentFiles[attachment.FileName]),
		}
	}

	return Attachment{
		Exist:     false,
		Extension: filepath.Ext(attachment.FileName),
		FileName:  attachment.FileName,
	}
}

// parserDate recieve the date with this format `06_01_2020=23:25`
func parserDate(date string) (string, string, string, string) {
	vals := strings.Split(date, "=")
	dates := strings.Split(vals[0], "_")

	month := dates[0]
	day := dates[1]
	year := dates[2]
	hours := vals[1]

	return month, day, year, hours
}

func getTranslateDate(lang string, month string, day string, year string) string {
	t, _ := translations[lang].(map[string]interface{})

	months, _ := t["months"].(map[string]string)
	date_template, _ := t["date_template"].(string)

	return fmt.Sprintf(date_template, months[month], day, year)
}

func (p *writertruct) UnmarshalJSONMessages(json_messages string, attachmentFiles map[string]string, attachmentURL string) Book {
	var language = "es"
	var temporalMessages []Message
	var messages []Message

	json.Unmarshal([]byte(json_messages), &temporalMessages)

	var sender string
	var receiver string
	var lastDate string

	// TODO: color to the badge of the name of the Author
	// sample: https://www.beautypunk.com/wp-content/uploads/2016/11/whatsapp-zapptales-buecher.jpg

	emojiConvert := emoji.NewEmojiParser()

	for _, m := range temporalMessages {
		month, day, year, hours := parserDate(m.Date)
		currentDate := fmt.Sprintf("%s_%s_%s", month, day, year)

		if len(lastDate) == 0 || !utils.IsEqualString(currentDate, lastDate) {
			badgeMessagesDate := Message{
				Date:       hours,
				Author:     m.Author,
				Message:    getTranslateDate(language, month, day, year),
				Attachment: Attachment{},
				IsSender:   false,
				IsReceiver: false,
				IsInfo:     true,
			}

			messages = append(messages, badgeMessagesDate)
		}

		lastDate = currentDate
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
				Date:       hours,
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
				Date:       hours,
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
			Date:       hours,
			Author:     "Info",
			Message:    messageValue,
			Attachment: Attachment{},
			IsSender:   false,
			IsReceiver: false,
			IsInfo:     true,
		}
		messages = append(messages, currentMessage)
	}

	return &export{
		messages: messages,
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
func (e *export) Export() []Message {
	return e.messages
}

// ExportHTML __
func (e *export) ExportHTML(paper Type) (string, error) {
	paperProps := *paper(e.messages)
	return paintPaper(paperProps)
}

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
