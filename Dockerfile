FROM golang:1.16 as wpparserbuild

WORKDIR /app

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o whatsapp-chat-parser cmd/cmd.go

FROM alpine:latest

WORKDIR /app

COPY --from=wpparserbuild /app/whatsapp-chat-parser .

ENV PATH $PATH:/app/

ENV PORT 4000
ENV GIN_MODE release

CMD [ "whatsapp-chat-parser" ]