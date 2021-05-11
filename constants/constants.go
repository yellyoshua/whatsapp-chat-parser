package constants

const (
	// DefaultPort default port of httpServer server
	DefaultPort      string = "4000"
	S3BucketName     string = "whatsapp-chat-parser"
	S3BucketEndpoint string = "https://BUCKET_NAME.s3.BUCKET_REGION.amazonaws.com"
)

const (
	MESSAGE_INDEX_DATE        = 1
	MESSAGE_INDEX_TIME        = 2
	MESSAGE_INDEX_TIME_FORMAT = 3
	MESSAGE_INDEX_AUTHOR      = 4
	MESSAGE_INDEX_MESSAGE     = 5
)

const (
	ContentTypeJSON string = "application/json; charset=utf-8"
)

const (
	KEY_MIDDLEWARE_ATTACHMENT_URL = "attachment_url"
)

const (
	FormatHTML string = "html"
	FormatFile string = "file"
	FormatJSON string = "json"
)
