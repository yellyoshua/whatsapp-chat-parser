FROM golang:1.16 as wpparserbuild

WORKDIR /app

COPY . .

RUN go build .

FROM alpine:latest

WORKDIR /app

COPY --from=wpparserbuild /app/whatsapp-chat-parser .

ENV PATH /app/:$PATH

CMD [ "whatsapp-chat-parser" ]