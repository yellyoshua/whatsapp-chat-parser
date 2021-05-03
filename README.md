# 💬 Whatsapp chat parser 📋

> `Actualizado 15 de abril del 2021`

# Introducci&oacute;n

Este paquete le permite convertir sus chats de whatsapp a formato JSON, PDF ó HTML.

# Pr&oacute;ximas Features
- Imagen QR como link a archivos de audio y video subidos a la nube
- Archivos audio y video reemplazados por texto `Archivo de audio ó Archivo de video`

# Instalar 💻

```
$ go get github.com/yellyoshua/whatsapp-chat-parser
```

Para ejecutar:

## Como libreria 📜

```golang
package main

import (
	"io/ioutil"
	"log"
	"path"

	"github.com/yellyoshua/whatsapp-chat-parser/paper"
	"github.com/yellyoshua/whatsapp-chat-parser/whatsapp"
)

func readFile(filename string) []byte {
	data, err := ioutil.ReadFile("test.txt")
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func initParsers() (whatsapp.Parser, paper.Writer) {
	var chat whatsapp.Parser = whatsapp.New()
	var writer = paper.New()
	return chat, writer
}

func main() {
	var whatsappChat = readFile("file.txt")
	var plainMessages string
	chat, writer := initParsers()

	err := chat.ParserMessages(whatsappChat, &plainMessages)
	if err != nil {
		log.Fatal(err)
	}

	book := writer.UnmarshalJSONMessages(
		plainMessages,
		path.Join("http://localhost:4000/public"),
	)

	messages := book.Export()

	log.Printf("Se procesaron %v messages", len(messages))
	log.Printf("%s - %s - %s",
		messages[0].Date,
		messages[0].Author,
		messages[0].Message,
	)
}
```

## Correr API REST en docker 🐳

Imagen subida a [Docker Registry](https://hub.docker.com/r/yellyoshua/whatsapp-chat-parser)

```
$ docker pull yellyoshua/whatsapp-chat-parser
```
Archivo `.env.production` para cargar las variables de entorno a la imagen de Docker

```.env
PORT=4000
AWS_ACCESS_KEY=
AWS_SECRET_KEY=
AWS_REGION="earth"
S3_BUCKET_NAME="whatsapp-chat-parser"
```

```
$ docker run --name whatsapp -p 4000:4000 --rm --env-file ./.env.production yellyoshua/whatsapp-chat-parser
```
&Oacute; agregar las variables de entorno por linea de comandos

```
$ docker run --name whatsapp -p 4000:4000 --rm -e S3_BUCKET_NAME="whatsapp-chat-parser" -e PORT=4000 -e AWS_REGION="earth" -e AWS_ACCESS_KEY="" -e AWS_SECRET_KEY="" yellyoshua/whatsapp-chat-parser
```


## Inicializar como API REST 😎

### Instalaci&oacute;n

```
$ git clone github.com/yellyoshua/whatsapp-chat-parser
$ cd whatsapp-chat-parser
$ make install
```

Para inicializar como API debe definir las siguientes variables de entorno.

<!-- prettier-ignore-start -->
| Name | Type | Description |
| :--- | :--- | :--- |
| PORT | Default: 4000| Puerto de escucha de la API REST |
| S3_BUCKET_NAME | `String`| Nombre del bucket de S3 |
| AWS_REGION | `String`| Regi&oacute;n del bucket de S3 |
| AWS_ACCESS_KEY | `String`| Clave de acceso S3 del bucket |
| AWS_SECRET_KEY | `String`| Clave secreta de acceso S3 del bucket |
<!-- prettier-ignore-end -->

____
____


Desp&uacute;es de haber definido las variables de entorno corremos los test's, construimos la api y la inicializamos:

```
$ make test
$ make build
$ ./whatsapp-chat-parser
```

# API

💬 Estructura del `paper.Message`

`github.com/yellyoshua/whatsapp-chat-parser/paper`

```golang
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
```

📚 Interfaz de `paper.Book`

`github.com/yellyoshua/whatsapp-chat-parser/paper`

```golang
// Book __
type Book interface {
	Export() []Message
	ExportJSON() (MessagesJSON, error)
	ExportHTML(paper Type) (string, error)
	ExportHTMLFile(paper Type, filePathName string) error
}
```

📋 Interfaz de `whatsapp.Parser`

`github.com/yellyoshua/whatsapp-chat-parser/whatsapp`

```golang
// Parser _
type Parser interface {
	ParserMessages(data []byte, outputMessages *string) error
}
```

## Tecnolog&iacute;as usadas

- Lenguaje: [Golang](https://golang.org/)
- Almacenamiento API REST: [S3](https://aws.amazon.com/s3/)
- Comunicaci&oacute;n API REST: [GIN-GONIC](github.com/gin-gonic/gin)


#### `Powered by yellyoshua `