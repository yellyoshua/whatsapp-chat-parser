FROM golang:1.16 as wpparserbuild

WORKDIR /app

COPY . .

# RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o whatsapp-chat-parser cmd/cmd.go
RUN echo "echo hola" > whatsapp-chat-parser

FROM hayd/alpine-deno:1.9.2

WORKDIR /app

COPY --from=wpparserbuild /app/whatsapp-chat-parser .
COPY --from=wpparserbuild /app/bin ./bin

# RUN yum install unzip -y && curl -fsSL https://deno.land/x/install/install.sh | sh
RUN apk update && apk upgrade && deno install -n chat-parser ./bin/parser-chat/index.ts

ENV PATH="$PATH:/app/"
ENV PORT 4000
ENV GIN_MODE release

CMD [ "whatsapp-chat-parser" ]