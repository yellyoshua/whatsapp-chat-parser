package paper

func bookStyle() string {
	stylesBook := `
  <style>
    :root {
      --sender-bg-color: white;
      --receiver-bg-color: rgb(39, 255, 118);
			--default-bg-color: #ff4465;
    }

    * {
      box-sizing: border-box;
      padding: 0;
      margin: 0;
      font-family: 'Noto Sans';
      font-weight: normal;
      font-style: normal;
			background: var(--default-bg-color);
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

    .banner-file-type {
      text-align: center !important;
      color: white !important;
      font-size: 14px !important;
      background: black !important;
      padding: 5px !important;
    }

    .sender-mini-box {
      left: -3px;
      background: linear-gradient(45deg, var(--sender-bg-color) 50%, transparent 50%) !important;
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

		.message-bubble.sender-message-bubble * {
			background: var(--sender-bg-color);
		}

		.message-bubble.receiver-message-bubble * {
			background: var(--receiver-bg-color);
		}

    .message-bubble {
      position: relative;
      width: 75%;
      padding: 10px;
      border-radius: 5px;
    }

    .message-author, .message-author p {
      font-size: 16px;
      margin-bottom: 10px;
      font-weight: bold;
    }

    .message-message {
      font-size: 16px;
      text-align: justify;
			word-break: break-word;
    }

    img.emoji {
      width: 19px;
      height: 19px;
    }

    .message-date, .message-date p {
      font-size: 12px;
      float: right;
      font-weight: bold;
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

	return stylesBook
}

// {{ if .Attachment.Extension = ".opus" }}
// <p class="banner-file-type" >Archivo de audio</p>
// {{ end }}
// {{ if .Attachment.Extension = ".mp4" }}
// <p class="banner-file-type" >Archivo de video</p>
// {{ end }}

func cardImage() string {
	cardImagePreview := `{{ if .Attachment.FileName}}
{{ if .Attachment.Exist}}
  <img src="{{.Attachment.FileName}}" width="100%" height="auto" alt="{{.Attachment.FileName}}" />
  {{ if .Attachment.Extension}}
    {{ if isAttachmentAudio .Attachment.Extension }}
    <p class="banner-file-type" >Archivo de audio</p>
    {{ else if isAttachmentVideo .Attachment.Extension }}
    <p class="banner-file-type" >Archivo de video</p>
    {{ else if isAttachmentImage .Attachment.Extension }}
    <p class="banner-file-type" >Archivo de imagen</p>
    {{ else }}
    <p class="banner-file-type">Archivo de documento</p>
    {{ end }}
  {{ else }}
  <p class="banner-file-type">Archivo Desconocido</p>
  {{ end }}
{{ else }}
<p class="banner-file-type" >El archivo no existe o eliminado del chat</p>
{{ end }}
{{ end }}`
	return cardImagePreview
}

func cardChatSender() string {
	cardSender := `{{ if .IsSender -}}
<div class="pg-break message-container sender-message-container">
	<div class="message-bubble sender-message-bubble">
		<div class="sender-mini-box mini-box"></div>
		<div class="message-author"><p>{{.Author}}</p></div>
    ` + cardImage() + `
		<div class="message-message">{{.Message}}</div>
		<div class="message-date"><p>{{.Date}}</p></div>
	</div>
</div>
{{ end }}`
	return cardSender
}

func cardChatReceiver() string {
	cardReceiver := `{{ if .IsReceiver -}}
<div class="pg-break message-container receiver-message-container">
  <div class="message-bubble receiver-message-bubble">
    <div class="receiver-mini-box mini-box"></div>
    <div class="message-author"><p>{{.Author}}</p></div>
    ` + cardImage() + `
    <div class="message-message">{{.Message}}</div>
    <div class="message-date"><p>{{.Date}}</p></div>
  </div>
</div>
{{ end }}`
	return cardReceiver
}

func bookWrapperChats() string {
	mappingMessages := `
  {{- range .Messages }}
  ` + cardChatSender() + `
  ` + cardChatReceiver() + `
	{{ end }}
  `
	return mappingMessages
}

func bookTemplate(bookStyles string) string {
	baseHTMLTemplate := `
  <!DOCTYPE html>
  <html lang="en">
  <head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
		<link rel="preconnect" href="https://fonts.gstatic.com">
    <link href="https://fonts.googleapis.com/css2?family=Noto+Sans:wght@400;700&display=swap" rel="stylesheet">
    <title>Document</title>
  </head>
  <body>
  ` + bookWrapperChats() + `
  ` + bookStyles + `
  </body>
  </html>
  `
	return baseHTMLTemplate
}
