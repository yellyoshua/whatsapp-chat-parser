# 💬 Whatsapp chat parser 📋

#### `Powered by yellyoshua `

<a href="https://www.buymeacoffee.com/yellyoshua" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me A Coffee" height="41" width="174" ></a>

> `Actualizado 10 de mayo del 2021`

# Introducci&oacute;n

Este paquete le permite convertir sus chats de whatsapp a formato JSON, PDF ó HTML.

# Pr&oacute;ximas Features
- Opci&oacute;n para editar el fondo del chat
- Opcion para editar los colores del chat
- Opcion para agregar una portada
- SOPORTE PARA CHAT GRUPAL

# Features
- Imagen QR como link a archivos de audio y video subidos a la nube
- Archivos audio y video reemplazados por texto `"Archivo de audio" ó "Archivo de video"`

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

	"github.com/yellyoshua/whatsapp-chat-parser/whatsapp"
)

func readFile(filename string) []byte {
	data, err := ioutil.ReadFile("test.txt")
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func main() {
	var whatsappChat = readFile("file.txt")
	
	// Este parametro es requerido como una 'clave' para la
	// creacion de la carpeta que contendra los archivos
	// del chat. Puede traducirse a un 'userID'
	uuid := "1e3e4e5e6e7e8e9e10e"
	
	chatBuilder := whatsapp.New()
	messages, err := chatBuilder.Parser(uuid, whatsappChat)
	if err != nil {
		logger.Fatalf("error al parsear el chat -> %s", err.Error())
	}

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
AWS_REGION=
S3_BUCKET_NAME="whatsapp-chat-parser"
```

```
$ docker run --name whatsapp -p 4000:4000 --rm --env-file ./.env.production yellyoshua/whatsapp-chat-parser
```
&Oacute; agregar las variables de entorno por linea de comandos

```
$ docker run --name whatsapp -p 4000:4000 --rm -e S3_BUCKET_NAME="whatsapp-chat-parser" -e PORT=4000 -e AWS_REGION="" -e AWS_ACCESS_KEY="" -e AWS_SECRET_KEY="" yellyoshua/whatsapp-chat-parser
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
Name | Type | Description
| -- | -- | -- |
PORT | Default: 4000| Puerto de escucha de la API REST |
S3_BUCKET_NAME | `String`| Nombre del bucket de S3 |
AWS_REGION | `String`| Regi&oacute;n del bucket de S3 |
AWS_ACCESS_KEY | `String`| Clave de acceso S3 del bucket |
AWS_SECRET_KEY | `String`| Clave secreta de acceso S3 del bucket |
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

💬 Estructura del `whatsapp.Message`

`github.com/yellyoshua/whatsapp-chat-parser/whatsapp`

```golang
type DateFormat struct {
	Hours  string `json:"hours"`
	Mins   string `json:"mins"`
	Format string `json:"format"`
	Day    int    `json:"day"`
	Month  int    `json:"month"`
	Year   int    `json:"year"`
	UTC    string `json:"utc"`
}

type Attachment struct {
	Exist     bool   `json:"exist"`
	FileName  string `json:"fileName,omitempty"`
	Extension string `json:"extension,omitempty"`
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
```

📚 Interfaz de `paper.Book`

`github.com/yellyoshua/whatsapp-chat-parser/paper`

```golang
// Book __
type Book interface {
	Export() []whatsapp.Message
	ExportJSON() (MessagesJSON, error)
	ExportHTML(paper Type) (string, error)
	ExportHTMLFile(paper Type, filePathName string) error
}
```

📋 Estructura de `whatsapp.ChatBuilder`

`github.com/yellyoshua/whatsapp-chat-parser/whatsapp`

```golang

func New() *ChatBuilder

func (r *ChatBuilder) Parser(user_id string, chat []byte) ([]Message, error)

```

## Tecnolog&iacute;as usadas ``GENERAL``

- HTML templates para la renderizaci&oacute;n del Whatsapp Book
- Lenguaje: [Golang](https://golang.org/)
- Para la creaci&oacute;n de images QR: [GO-QRCODE](https://github.com/skip2/go-qrcode)

## Tecnolog&iacute;as usadas ``API-REST``

- API-REST: [GIN-GONIC](github.com/gin-gonic/gin/)
- Almacenamiento: [S3](https://aws.amazon.com/s3/)

#### `Powered by yellyoshua `