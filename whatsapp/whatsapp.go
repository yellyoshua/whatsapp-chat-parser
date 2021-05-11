package whatsapp

// RegexAttachment format input `$date - $contact: IMG-20200319-WA0011.jpg (file attached)`
var RegexAttachment string = `(: )+[\S\s]+(\.\w{2,4}\s)+\(+(file attached)+\)`

// RegexAttachmentText format input `IMG-20200319-WA0011.jpg (file attached)`
var RegexAttachmentText string = `\(file attached\)`

var RegexpMessage string = `^(?:U+200E|U+200F)*\[?(\d{1,4}[-/.] ?\d{1,4}[-/.] ?\d{1,4})[,.]? \D*?(\d{1,2}[.:]\d{1,2}(?:[.:]\d{1,2})?)(?: ([ap]\.? ?m\.?))?\]?(?: -|:)? (.+?): (.+)`
var RegexpMessageSystem string = `^(?:U+200E|U+200F)*\[?(\d{1,4}[-/.] ?\d{1,4}[-/.] ?\d{1,4})[,.]? \D*?(\d{1,2}[.:]\d{1,2}(?:[.:]\d{1,2})?)(?: ([ap]\.? ?m\.?))?\]?(?: -|:)? (.+)`
var RegexMessageAttachment string = `<.+:(.+)>`

type DateFormat struct {
	Hours  string `json:"hours"`
	Mins   string `json:"mins"`
	Format string `json:"format"`
	Day    int    `json:"day"`
	Month  int    `json:"month"`
	Year   int    `json:"year"`
	UTC    string `json:"utc"`
}

// Attachment _
type Attachment struct {
	Exist     bool   `json:"exist"`
	FileName  string `json:"fileName,omitempty"`
	Extension string `json:"extension,omitempty"`
}

// RawMessage _
type RawMessage struct {
	IsSystem bool   `json:"isSystem"`
	Message  string `json:"msg"`
}

// Message _
type Message struct {
	Date       DateFormat `json:"date"`
	Author     string     `json:"author"`
	IsSender   bool       `json:"isSender"`
	IsInfo     bool       `json:"isInfo"`
	IsReceiver bool       `json:"isReceiver"`
	Message    string     `json:"message"`
	Attachment Attachment `json:"attachment"`
}

type MessagesBroker struct {
	splitMessages []string
}

func New(user_id string, chat string) *MessagesBroker {
	splitedMessages := splitChatMessages(chat)
	return &MessagesBroker{
		splitMessages: splitedMessages,
	}
}

func (mc *MessagesBroker) parsePlainMessages() []RawMessage {
	rawMessages := parsePlainMessages(mc.splitMessages)
	return rawMessages
}

func (mc *MessagesBroker) Messages() []Message {
	rawMessage := mc.parsePlainMessages()
	return rawMessagesToMessages(rawMessage)
}
